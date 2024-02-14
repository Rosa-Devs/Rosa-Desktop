package core

import (
	"context"
	"fmt"
	"log"
	"sync"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
)

func bootstrap(ctx context.Context, dht *dht.IpfsDHT) {
	//Bootstrap
	log.Println("[DHT] Running bootstrap thread")
	if err := dht.Bootstrap(ctx); err != nil {
		log.Println("[DHT] Fail to bootstrap DHT: ", err)
	}

}

func init_DHT(ctx context.Context, host host.Host) *dht.IpfsDHT {
	// init DHT

	kademliaDHT, err := dht.New(ctx, host)
	if err != nil {
		log.Println("[DHT] Fail to init DHT: ", err)
	}
	log.Println("[DHT] Init sucsesfull")
	return kademliaDHT

}

func boot(ctx context.Context, host host.Host, d *dht.IpfsDHT) {

	var wg sync.WaitGroup
	for _, peerAddr := range dht.DefaultBootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := host.Connect(ctx, *peerinfo); err != nil {
				log.Println("[DHT:Bootstrap] Fail to connect:", err)
			} else {
				log.Println("[DHT:Bootstrap] Successfully connected to node: ", *peerinfo)
			}
		}()
	}
	wg.Wait()

	routingDiscovery := drouting.NewRoutingDiscovery(d)
	dutil.Advertise(ctx, routingDiscovery, Randevuz)

	// Look for others who have announced and attempt to connect to them
	anyConnected := false
	for !anyConnected {
		fmt.Println("Searching for peers...")
		peerChan, err := routingDiscovery.FindPeers(ctx, Randevuz)
		if err != nil {
			panic(err)
		}
		for p := range peerChan {
			if p.ID == host.ID() {
				continue // No self connection
			}
			go func() {
				err := host.Connect(ctx, p)
				if err != nil {
					fmt.Println("[DHT] Failed connecting to ", p.ID)
				} else {
					fmt.Println("[DHT] Connected to:", p.ID)

				}
			}()
			anyConnected = true
		}
	}
	fmt.Println("[DHT] Peer discovery complete")

}
