package network

import (
	"context"
	"log"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
)

type Host struct {
	ctx context.Context

	MDnsServie bool
	DhtService bool

	// Libp2p host for routing
	H   host.Host
	Dht *dht.IpfsDHT
	Ps  *pubsub.PubSub
}

func (h *Host) InitHost(ctx context.Context) error {
	//Init a context
	h.ctx = ctx

	//Defiend opts for libp2p stack
	opts := libp2p.ChainOptions(
		libp2p.EnableNATService(),
		libp2p.EnableRelayService(),
		libp2p.EnableRelay(),
		libp2p.EnableHolePunching(),
	)

	//Create new libp2p instance

	var err error
	h.H, err = libp2p.New(opts)
	if err != nil {
		return err
	}

	//Cerate new Gossip Router
	h.Ps, err = pubsub.NewGossipSub(h.ctx, h.H)
	if err != nil {
		return err
	}

	//Run host if needed
	if h.MDnsServie {
		if h.StartmDSNService() != nil {
			log.Println("Failed to launch mDns service:", err)
		}
	}
	if h.DhtService {
		if h.init_DHT() != nil {
			log.Println("Failede to launch Dht service:", err)
			return nil
		}
		go h.connectToDefaultNodes()
	}

	return nil
}
