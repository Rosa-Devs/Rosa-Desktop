package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Rosa-Devs/Database/src/manifest"
	db "github.com/Rosa-Devs/Database/src/store"
	"github.com/Rosa-Devs/Rosa-Desktop/models"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (Mgr *Core) DatabaseUpdateEventServer(ctx context.Context, m manifest.Manifest) {
	Mgr.waitGrp.Add(1)

	go func(ctx context.Context, m manifest.Manifest, M *Core) {
		defer M.waitGrp.Done()

		// Get database
		database, ok := M.dbs[m]
		if !ok {
			log.Println("Fail to get db!")
			return
		}

		// Subscribe to event channel
		eventListener := make(chan db.Event)
		database.EventBus.Subscribe(db.DbUpdateEvent, eventListener)

		for {
			select {
			case <-ctx.Done():
				log.Println("DatabaseUpdateEventServer exiting.")
				return
			case <-eventListener:
				log.Println("Update event")
				runtime.EventsEmit(M.wailsctx, "update")
				time.Sleep(time.Second)
				//log.Println("Event")
			}
		}
	}(ctx, m, Mgr)
}

func (Mgr *Core) ChangeListeningDb(m manifest.Manifest) {
	if Mgr.Started == false {
		log.Println("Cannot change event server because DBMgr not running")
		return
	}

	// Create a cancelable context
	ctx, cancel := context.WithCancel(context.Background())

	// Replace the existing cancel function with the new one
	// to ensure that the old goroutine will exit.
	if Mgr.cancelFunc != nil {
		Mgr.cancelFunc()
	}

	// Set the new cancel function
	Mgr.cancelFunc = cancel

	// Create new DatabaseUpdateEventServer with new data!
	log.Println("Recreating events server")
	Mgr.DatabaseUpdateEventServer(ctx, m)
}

func (Mgr *Core) Nodes() int {
	if Mgr.Started == false {
		return 0
	}
	return len(Mgr.host.Dht.Host().Network().Peers())
}

func (Mgr *Core) NewMessage(m manifest.Manifest, msg string) error {

	msg_stuct := new(models.Message)

	currentTime := time.Now().UTC()

	// Format the time in the desired format
	timestamp := currentTime.Format("2006-01-02T15:04:05.000")

	msg_stuct.Data = msg
	msg_stuct.Sender = Mgr.profile.GetPublic()
	msg_stuct.Sender.Avatar = ""
	msg_stuct.SenderId = Mgr.profile.Id
	msg_stuct.Time = timestamp + "1"
	msg_stuct.DataType = models.MessageType
	msg_stuct.Valid = false
	Mgr.profile.Sign(msg_stuct)
	p_pub := Mgr.profile.GetPublic()
	b := p_pub.ValidateMsg(*msg_stuct)
	log.Println("Self signed msg. Sigh verified:", b)

	log.Println(msg_stuct.Signature)

	//GET DB FROM DATABASE MGR
	db, ok := Mgr.dbs[m]
	if !ok {
		log.Println("Failed to get database from dbs manager.")
		return fmt.Errorf("Failed to get database from dbs manager")
	}

	err := db.CreatePool(models.MsgPool)
	if err != nil {
		//log.Println("Not recreating pool:", err)
	}

	pool, err := db.GetPool(models.MsgPool)
	if err != nil {
		log.Println("Failed to get pool:", err)
		return err
	}

	//SERIALIZE MSG

	msgBytes, err := json.Marshal(msg_stuct)
	if err != nil {
		log.Println("Failed to serialize message:", err)
		return err
	}

	pool.Record(msgBytes)

	return nil
}

func (Mgr *Core) GetMessages(m manifest.Manifest) ([]models.Message, error) {

	startTime := time.Now()

	db, ok := Mgr.dbs[m]
	//log.Println(Mgr.dbs)
	if !ok {
		log.Println("Failed to get database from dbs manager.")
		return nil, fmt.Errorf("Failed to get database from dbs manager 1")
	}

	err := db.CreatePool(models.MsgPool)
	if err != nil {
		//log.Println("Not recreating pool:", err)
	}

	pool, err := db.GetPool(models.MsgPool)
	if err != nil {
		log.Println("Failed to get pool:", err)
		return nil, err
	}

	filter := map[string]interface{}{
		"datatype": models.MessageType,
	}

	data, err := pool.Filter(filter)
	if err != nil {
		fmt.Println("Data:", data)
		fmt.Println("Error filtering data:", err)
	}

	msg_data := convertToMessages(data, db)
	Mgr.Validator(&msg_data)

	// log.Println(msg_data[0])
	// sort.Slice(msg_data, func(i, j int) bool {
	// 	timei, _ := time.Parse(time.RFC3339, msg_data[i].Time) // Assuming Data field contains timestamp in RFC3339 format
	// 	timej, _ := time.Parse(time.RFC3339, msg_data[j].Time)

	// 	// Compare timestamps to sort from latest to newest
	// 	return timei.After(timej)
	// })

	stopTime := time.Now()
	//log.Println(msg_data)
	elapsedTime := stopTime.Sub(startTime)
	log.Printf("GetMessages execution time: %v\n", elapsedTime)

	return msg_data, nil
}

func convertToMessages(data []map[string]interface{}, db *db.Database) []models.Message {
	messages := make([]models.Message, len(data))

	for i, item := range data {
		// Assuming your map contains fields like "ID" and "Text"
		// Adjust these according to your actual map structure
		d_type, _ := item["datatype"].(int)
		sender_map, _ := item["sender"].(map[string]interface{})
		sender := new(models.ProfileStorePublic)
		sender.ProfileFromMap(sender_map)

		filter := map[string]interface{}{
			"id": sender.Id, // Random integer between 0 and 100
		}

		pool, err := db.GetPool(models.UserPool)
		if err != nil {
			log.Println(err)
			return nil
		}

		profiles, err := pool.Filter(filter)
		if err != nil {
			log.Println("Fail to get pool:", err)
			return nil
		}

		if len(profiles) > 0 {
			p := new(models.ProfileStorePublic)
			p.ProfileFromMap(profiles[0])
			sender.Avatar = p.Avatar
		} else {
			sender.Avatar = ""
		}

		data, _ := item["data"].(string)
		time, _ := item["time"].(string)
		sign, _ := item["sign"].(string)
		s_id, _ := item["senderid"].(string)

		// Create a new Message and append it to the result slice
		messages[i] = models.Message{
			DataType:  d_type,
			Sender:    *sender,
			SenderId:  s_id,
			Data:      data,
			Time:      time,
			Signature: sign,
		}
	}

	return messages
}

func (c *Core) Validator(m *[]models.Message) {
	for i := range *m {
		user := c.FindUserById((*m)[i].Sender.Id)
		if user.Id != (*m)[i].Sender.Id {
			//log.Println("This account not is exsit:", user.Id)
		}

		msg := (*m)[i]

		if validated, ok := c.MessageValidateСache[msg.Data+msg.Time]; ok {
			//log.Println("Msg:", msg.Data, "Already Validated for user:", user.Name)
			(*m)[i].Valid = validated
			continue
		} else {
			if user.ValidateMsg(msg) {
				//log.Println("Msg:", (*m)[i].Data, "Validated:", "true", "With user:", user.Name)
				(*m)[i].Valid = true
				//.MessageValidateСache[msg.Data+msg.Time] = true
			} else {
				(*m)[i].Valid = false
				c.MessageValidateСache[msg.Data+msg.Time] = false
				//log.Println("Msg:", (*m)[i].Data, "Validated:", "false", "With user:", user.Name)
			}
		}
	}
}

// func (c *Core) Validator(m *[]models.Message) {
// 	for i := range *m {
// 		user := c.FindUserById((*m)[i].Sender.Id)
// 		if user.Id != (*m)[i].Sender.Id {
// 			log.Println("This account not is exsit:", user.Id)
// 		}

// 		msg := (*m)[i]

// 		if user.ValidateMsg(msg) {
// 			log.Println("Msg:", (*m)[i].Data, "Validated:", "true", "With user:", user.Name)
// 			(*m)[i].Valid = true
// 		} else {
// 			(*m)[i].Valid = false
// 			log.Println("Msg:", (*m)[i].Data, "Validated:", "false", "With user:", user.Name)
// 		}
// 	}
// }
