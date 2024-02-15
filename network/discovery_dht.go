package network

import (
	"fmt"
	"log"
	"sync"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/peer"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
)

func (h *Host) init_DHT() error {
	// init DHT

	kademliaDHT, err := dht.New(h.ctx, h.H)
	if err != nil {
		log.Println("[DHT] Fail to init DHT: ", err)
		return err
	}

	log.Println("[DHT] Running bootstrap thread")
	if err := kademliaDHT.Bootstrap(h.ctx); err != nil {
		log.Println("[DHT] Fail to bootstrap DHT: ", err)
		return err
	}
	log.Println("[DHT] Init sucsesfull")

	h.Dht = kademliaDHT
	return nil

}

func (h *Host) connectToDefaultNodes() {

	var wg sync.WaitGroup
	for _, peerAddr := range dht.DefaultBootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := h.H.Connect(h.ctx, *peerinfo); err != nil {
				log.Println("[DHT:Bootstrap] Fail to connect:", err)
			} else {
				log.Println("[DHT:Bootstrap] Successfully connected to node: ", *peerinfo)
			}
		}()
	}
	wg.Wait()

	routingDiscovery := drouting.NewRoutingDiscovery(h.Dht)
	dutil.Advertise(h.ctx, routingDiscovery, Randevuz)

	// Look for others who have announced and attempt to connect to them
	anyConnected := false
	for !anyConnected {
		fmt.Println("Searching for peers...")
		peerChan, err := routingDiscovery.FindPeers(h.ctx, Randevuz)
		if err != nil {
			panic(err)
		}
		for p := range peerChan {
			if p.ID == h.H.ID() {
				continue // No self connection
			}
			go func() {
				err := h.H.Connect(h.ctx, p)
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
