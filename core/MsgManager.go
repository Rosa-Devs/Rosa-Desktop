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
	msg_stuct.Sender = Mgr.profile.Name
	msg_stuct.Time = timestamp + "1"
	msg_stuct.DataType = models.MessageType

	//GET DB FROM DATABASE MGR
	db, ok := Mgr.dbs[m]
	if !ok {
		log.Println("Failed to get database from dbs manager.")
		return fmt.Errorf("Failed to get database from dbs manager")
	}

	err := db.CreatePool(models.MsgPool)
	if err != nil {
		log.Println("Not recreating pool:", err)
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

	msg_data := convertToMessages(data)

	// sort.Slice(msg_data, func(i, j int) bool {
	// 	timei, _ := time.Parse(time.RFC3339, msg_data[i].Time) // Assuming Data field contains timestamp in RFC3339 format
	// 	timej, _ := time.Parse(time.RFC3339, msg_data[j].Time)

	// 	// Compare timestamps to sort from latest to newest
	// 	return timei.After(timej)
	// })

	//log.Println(msg_data)

	return msg_data, nil
}

func convertToMessages(data []map[string]interface{}) []models.Message {
	messages := make([]models.Message, len(data))

	for i, item := range data {
		// Assuming your map contains fields like "ID" and "Text"
		// Adjust these according to your actual map structure
		d_type, _ := item["datatype"].(int)
		sender, _ := item["sender"].(string)
		data, _ := item["data"].(string)
		time, _ := item["time"].(string)

		// Create a new Message and append it to the result slice
		messages[i] = models.Message{
			DataType: d_type,
			Sender:   sender,
			Data:     data,
			Time:     time,
		}
	}

	return messages
}
