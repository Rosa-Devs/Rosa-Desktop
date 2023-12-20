package src

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Rosa-Devs/Database/src/manifest"
	db "github.com/Rosa-Devs/Database/src/store"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

type DbManager struct {
	Started bool
	DbPath  string

	//TEST
	Name string

	ctx context.Context
	h   host.Host
	dht *dht.IpfsDHT
	ps  *pubsub.PubSub
	dbs map[manifest.Manifest]*db.Database

	Driver      *db.DB
	Manifest_DB db.Database

	//Event server
	stopCh     chan struct{}
	waitGrp    sync.WaitGroup
	cancelFunc context.CancelFunc
	wailsctx   context.Context
}

func (d *DbManager) GetProfile() string {
	if d.Started == false {
		return ""
	}
	return d.Name
}

func (d *DbManager) OnWailsInit(ctx context.Context) {
	d.wailsctx = ctx
}

func (d *DbManager) StartManager(dbPath string, N string) {
	d.stopCh = make(chan struct{})
	d.Name = N

	if d.Started == true {
		log.Println("Dbs Manager already started..")
		return
	}
	d.Started = true

	d.DbPath = dbPath
	d.dbs = make(map[manifest.Manifest]*db.Database)
	d.ctx = context.Background()

	var err error
	d.h, err = libp2p.New()
	if err != nil {
		panic(err)
	}

	d.ps, err = pubsub.NewGossipSub(d.ctx, d.h)
	if err != nil {
		panic(err)
	}

	if err := setupDiscovery(d.h); err != nil {
		panic(err)
	}

	d.dht = init_DHT(d.ctx, d.h)
	go bootstrap(d.ctx, d.dht)
	go boot(d.ctx, "rosa", d.h)

	d.Driver = &db.DB{
		H:  d.h,
		Pb: d.ps,
	}
	d.Driver.Start(d.DbPath)

	m_db := manifest.Manifest{
		Name:   "Manifests",
		PubSub: manifest.GenerateNoise(15),
	}

	d.Driver.CreateDb(m_db)
	d.Manifest_DB = d.Driver.GetDb(m_db)

	err = d.Manifest_DB.CreatePool("manifests")
	if err != nil {
		log.Println("Not recreating pool:", err)
	}
	//READ MANIFET DB AND CREATE DBS
	pool, err := d.Manifest_DB.GetPool("manifests")
	if err != nil {
		log.Println("Failed to get pool")
	}

	filter := map[string]interface{}{
		"type": 1, // All manifests
	}

	data, err := pool.Filter(filter)
	if err != nil {
		fmt.Println("Data:", data)
		fmt.Println("Error filtering data:", err)
	}

	for _, record := range data {
		//log.Println(record)
		manifestData, ok := record["data"].(string)
		if !ok {
			fmt.Println("Data field not found in map")
			continue
		}

		decodedData, err := base64.StdEncoding.DecodeString(manifestData)
		if err != nil {
			log.Println("Error decoding base64 data:", err)
			continue
		}

		m := new(manifest.Manifest)
		err = m.Deserialize(decodedData)
		if err != nil {
			log.Println("Error deserializing manifest, err:", err)
			continue
		}

		//Try to create db
		err = d.Driver.CreateDb(*m)
		if err != nil {
			log.Println("Not recreating db db, err:", err)
		}
		//Get db by manifest
		db := d.Driver.GetDb(*m)
		db.StartWorker(15)
		d.dbs[*m] = &db
	}

	log.Println("All database are create and ready to use")

}

func (d *DbManager) AddManifets(manifestJson string) error {
	if d.Started == false {
		log.Println("Db manager is not started")
		return fmt.Errorf("Db manager is not started")
	}
	err := d.Manifest_DB.CreatePool("manifests")
	if err != nil {
		log.Println("Not recreating pool:", err)
	}
	pool, err := d.Manifest_DB.GetPool("manifests")
	if err != nil {
		log.Println("Fail to get pool", err)
		return err
	}

	m := new(manifest.Manifest)
	m.Deserialize([]byte(manifestJson))

	m_s := new(MStore)
	m_s.Data, err = m.Serialize()
	m_s.Type = MStore_TYPE_Manifet
	if err != nil {
		log.Panicln("Fail to serialize manifest", err)
		return err
	}
	jsonData, err := json.Marshal(m_s)
	if err != nil {
		log.Println("Fail to marshal manifest", err)
	}

	//Create new db
	db := d.Driver.GetDb(*m)
	db.StartWorker(15)
	d.dbs[*m] = &db

	err = pool.Record(jsonData)
	if err != nil {
		log.Println("Fail to update pool", err)
		return err
	}

	return nil
}

func (d *DbManager) ManifestList() []manifest.Manifest {
	if d.Started == false {
		log.Println("Db manager is not started")
		return append([]manifest.Manifest{}, manifest.Manifest{Name: "Db Manager not started", PubSub: "0"})
	}
	err := d.Manifest_DB.CreatePool("manifests")
	if err != nil {
		//log.Println("Not recreating pool:", err)
	}
	//READ MANIFET DB AND CREATE DBS
	pool, err := d.Manifest_DB.GetPool("manifests")
	if err != nil {
		log.Println("Failed to get pool")
	}

	filter := map[string]interface{}{
		"type": 1, // All manifests
	}

	data, err := pool.Filter(filter)
	if err != nil {
		fmt.Println("Data:", data)
		fmt.Println("Error filtering data:", err)
	}

	var manifetss []manifest.Manifest
	for _, record := range data {
		//log.Println(record)
		manifestData, ok := record["data"].(string)
		if !ok {
			fmt.Println("Data field not found in map")
			continue
		}

		decodedData, err := base64.StdEncoding.DecodeString(manifestData)
		if err != nil {
			log.Println("Error decoding base64 data:", err)
			continue
		}

		m := new(manifest.Manifest)
		err = m.Deserialize(decodedData)
		if err != nil {
			log.Println("Error deserializing manifest, err:", err)
			continue
		}

		manifetss = append(manifetss, *m)
	}
	return manifetss
}

func (d *DbManager) DeleteManifest(m manifest.Manifest) error {
	if d.Started == false {
		log.Println("Db manager is not started")
		return fmt.Errorf("Db manager is not started")
	}
	err := d.Manifest_DB.CreatePool("manifests")
	if err != nil {
		log.Println("Not recreating pool:", err)
	}
	//READ MANIFET DB AND CREATE DBS
	pool, err := d.Manifest_DB.GetPool("manifests")
	if err != nil {
		log.Println("Failed to get pool")
		return err
	}

	m_d, err := m.Serialize()
	if err != nil {
		log.Println("Failed to serialize manifest:", err)
		return err
	}
	encodedData := base64.StdEncoding.EncodeToString(m_d)

	filter := map[string]interface{}{
		"type": 1, // All manifests
		"data": encodedData,
	}

	data, err := pool.Filter(filter)
	if err != nil {
		fmt.Println("Data:", data)
		fmt.Println("Error filtering data:", err)
	}

	log.Println(data)

	for _, record := range data {
		err := pool.Delete(record["_id"].(string))
		if err != nil {
			log.Println("Error deleting record:", err)
			return err
		}
	}

	return nil
}

// DiscoveryInterval is how often we re-publish our mDNS records.
const DiscoveryInterval = time.Hour

// DiscoveryServiceTag is used in our mDNS advertisements to discover other chat peers.
const DiscoveryServiceTag = "pubsub-chat-example"

// discoveryNotifee gets notified when we find a new peer via mDNS discovery
type discoveryNotifee struct {
	h host.Host
}

// HandlePeerFound connects to peers discovered via mDNS. Once they're connected,
// the PubSub system will automatically start interacting with them if they also
// support PubSub.
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	fmt.Printf("discovered new peer %s\n", pi.ID)
	err := n.h.Connect(context.Background(), pi)
	if err != nil {
		fmt.Printf("error connecting to peer %s: %s\n", pi.ID, err)
	}
}

// setupDiscovery creates an mDNS discovery service and attaches it to the libp2p Host.
// This lets us automatically discover peers on the same LAN and connect to them.
func setupDiscovery(h host.Host) error {
	// setup mDNS discovery to find local peers
	s := mdns.NewMdnsService(h, DiscoveryServiceTag, &discoveryNotifee{h: h})
	return s.Start()
}
