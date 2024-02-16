package core

import (
	"log"

	"github.com/Rosa-Devs/Rosa-Desktop/models"
)

func (c *Core) TrusNewProfile(p models.ProfileStorePublic) {
	if !c.Started {
		log.Println("Db manger not started yet")
		return
	}

	pool, err := c.Service_DB.GetPool("trust")
	if err != nil {
		log.Println(err)
		return
	}

	data, err := p.Serialize()
	if err != nil {
		log.Println("Fail to serialize public profile")
		return
	}

	err = pool.Record(data)
	if err != nil {
		log.Println(err)
		return
	}
}
