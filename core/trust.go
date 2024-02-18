package core

import (
	"log"

	"github.com/Rosa-Devs/Rosa-Desktop/models"
)

func (c *Core) TrustNewProfile(p models.Message) {
	if !c.Started {
		log.Println("Db manger not started yet")
		return
	}
	pool, err := c.Service_DB.GetPool("trust")
	if err != nil {
		log.Println(err)
		return
	}

	filter := map[string]interface{}{
		"id": p.Sender.Id,
	}

	profiles, err := pool.Filter(filter)
	if err != nil {
		log.Println("Fail to get pool:", err)
		return
	}

	if len(profiles) > 0 {
		log.Println("This account exist!")
		return
	}

	data, err := p.Sender.Serialize()
	if err != nil {
		log.Println("Fail to serialize public profile", err)
		return
	}

	err = pool.Record(data)
	if err != nil {
		log.Println(err)
		return
	}

}

func (c *Core) FindUserById(id string) models.ProfileStorePublic {

	if !c.Started {
		log.Println("Db manager is not stated")
		return models.ProfileStorePublic{}
	}

	pool, err := c.Service_DB.GetPool("trust")
	if err != nil {
		log.Println("Fail to get trust pool:", err)
		return models.ProfileStorePublic{}
	}

	filter := map[string]interface{}{
		"id": id, // Random integer between 0 and 100
	}

	data, err := pool.Filter(filter)
	if err != nil {
		log.Println("Fail to get pool:", err)
		return models.ProfileStorePublic{}
	}

	profile := new(models.ProfileStorePublic)
	if len(data) > 0 {
		profile.ProfileFromMap(data[0])
	}

	return *profile
}
