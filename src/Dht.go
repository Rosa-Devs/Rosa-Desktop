package src

import (
	"context"
	"log"
	"sync"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
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

func boot(ctx context.Context, randevu string, host host.Host) {

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

}
