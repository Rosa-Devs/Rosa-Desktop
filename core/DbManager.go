package core

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/Rosa-Devs/Database/src/manifest"
	db "github.com/Rosa-Devs/Database/src/store"
	"github.com/Rosa-Devs/Rosa-Desktop/models"
	"github.com/Rosa-Devs/Rosa-Desktop/network"
	"github.com/Rosa-Devs/Rosa-Desktop/store"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const Randevuz = "RosaApp"

type Core struct {
	Icon    []byte
	Started bool
	DbPath  string

	//TEST
	Store   store.Store
	profile models.Profile

	ctx context.Context

	//Host
	host network.Host

	dbs map[manifest.Manifest]*db.Database

	Driver     *db.DB
	Service_DB db.Database

	//Event server
	stopCh     chan struct{}
	waitGrp    sync.WaitGroup
	cancelFunc context.CancelFunc
	wailsctx   context.Context
}

func (d *Core) GetProfile() string {
	if d.Started == false {
		return ""
	}
	return d.profile.Name
}

func (d *Core) OnWailsInit(ctx context.Context) {
	d.wailsctx = ctx
}

func (d *Core) StartManager() {
	d.stopCh = make(chan struct{})

	var err error
	d.profile, err = models.LoadFromFile(d.Store.Profile)
	if err != nil {
		fmt.Println("Error loading profile:", err)
		d.profile = models.Profile{
			Id: "UAUNT",
		}
		return
	}

	if d.Started == true {
		log.Println("Dbs Manager already started..")
		return
	}
	d.Started = true

	d.DbPath = d.Store.Database
	d.dbs = make(map[manifest.Manifest]*db.Database)
	d.ctx = context.Background()

	// Create new Host instance with properties
	d.host = network.Host{
		MDnsServie: true,
		DhtService: true,
	}

	if d.host.InitHost(d.ctx) != nil {
		log.Println("Failt to init HOST module. Crytical error")
		return
	}

	d.Driver = &db.DB{
		H:  d.host.H,
		Pb: d.host.Ps,
	}
	d.Driver.Start(d.DbPath)

	m_db := manifest.Manifest{
		Name:   "Service",
		PubSub: manifest.GenerateNoise(15),
		Chiper: manifest.GenerateNoise(32),
	}

	d.Driver.CreateDb(m_db)
	d.Service_DB = d.Driver.GetDb(m_db)

	err = d.Service_DB.CreatePool("manifests")
	if err != nil {
		//log.Println("Not recreating pool:", err)
	}
	err = d.Service_DB.CreatePool("trust")
	if err != nil {
		//log.Println("Not recreating pool:", err)
	}

	//READ MANIFET DB AND CREATE DBS
	pool, err := d.Service_DB.GetPool("manifests")
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
		db.StartWorker(35)
		d.dbs[*m] = &db
	}

	log.Println("All database are create and ready to use")
}

func (d *Core) CreateManifest(name string, opts string) string {

	m_json, err := manifest.GenereateManifest(name, false, opts).Serialize()
	if err != nil {
		log.Println("Fail to create manifest")
	}

	return string(m_json)

}

func (d *Core) AddManifets(manifestJson string) error {
	if d.Started == false {
		log.Println("Db manager is not started")
		return fmt.Errorf("Db manager is not started")
	}
	err := d.Service_DB.CreatePool("manifests")
	if err != nil {
		log.Println("Not recreating pool:", err)
	}
	pool, err := d.Service_DB.GetPool("manifests")
	if err != nil {
		log.Println("Fail to get pool", err)
		return err
	}

	m := new(manifest.Manifest)
	m.Deserialize([]byte(manifestJson))

	m_s := new(models.MStore)
	m_s.Data, err = m.Serialize()
	m_s.Type = models.MStore_TYPE_Manifet
	if err != nil {
		log.Panicln("Fail to serialize manifest", err)
		return err
	}
	jsonData, err := json.Marshal(m_s)
	if err != nil {
		log.Println("Fail to marshal manifest", err)
	}

	//Create new db
	//Try to create db
	err = d.Driver.CreateDb(*m)
	if err != nil {
		log.Println("Not recreating db db, err:", err)
		return fmt.Errorf("Not reacreatinng db, err:", err)
	}
	db := d.Driver.GetDb(*m)
	db.StartWorker(60)
	d.dbs[*m] = &db

	err = pool.Record(jsonData)
	if err != nil {
		log.Println("Fail to update pool", err)
		return err
	}

	return nil

}

func (d *Core) ManifestList() []manifest.Manifest {
	if d.Started == false {
		log.Println("Db manager is not started")
		return append([]manifest.Manifest{}, manifest.Manifest{Name: "Db Manager not started", PubSub: "0"})
	}
	err := d.Service_DB.CreatePool("manifests")
	if err != nil {
		//log.Println("Not recreating pool:", err)
	}
	//READ MANIFET DB AND CREATE DBS
	pool, err := d.Service_DB.GetPool("manifests")
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

func (d *Core) DeleteManifest(m manifest.Manifest) error {
	if d.Started == false {
		log.Println("Db manager is not started")
		return fmt.Errorf("Db manager is not started")
	}
	err := d.Service_DB.CreatePool("manifests")
	if err != nil {
		log.Println("Not recreating pool:", err)
	}
	//READ MANIFET DB AND CREATE DBS
	pool, err := d.Service_DB.GetPool("manifests")
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

func (c *Core) CreateNewAccount(name string, avatar string) error {

	profile, err := models.CreateProfile(name, avatar)
	if err != nil {
		log.Println(err)
		return err
	}
	err = models.WriteToFile(c.Store.Profile, profile)
	if err != nil {
		log.Println(err)
		return err
	}

	//For frontend
	c.StartManager()
	return nil
}

func (c *Core) Autorized() bool {
	if !c.Started {
		return false
	}
	if c.profile.Id == "UAUNT" {
		return false
	}

	return true
}

func (c *Core) ExportManifest(m manifest.Manifest) {
	path, err := runtime.SaveFileDialog(c.wailsctx, runtime.SaveDialogOptions{
		Title:                "Save chat file...",
		CanCreateDirectories: true,
		DefaultFilename:      m.Name + ".json",
	})
	if err != nil {
		log.Println("Fail to chosee export path")
		return
	}

	data, err := m.Serialize()
	if err != nil {
		log.Println("Fail to serialize")
		return
	}

	err = os.WriteFile(path, data, os.FileMode(0775))
	if err != nil {
		log.Println("Fail to write a file!!")
		return
	}

	log.Println(path)
}
