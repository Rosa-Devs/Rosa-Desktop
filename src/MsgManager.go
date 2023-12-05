package src

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Rosa-Devs/Database/src/manifest"
)

func (Mgr *DbManager) NewMessage(m manifest.Manifest, msg string) error {

	msg_stuct := new(Message)

	currentTime := time.Now().UTC()

	// Format the time in RFC3339 format
	timestamp := currentTime.Format(time.RFC3339)

	msg_stuct.Data = msg
	msg_stuct.Sender = Mgr.Name
	msg_stuct.Time = timestamp + "1"
	msg_stuct.DataType = MessageType

	//GET DB FROM DATABASE MGR

	db, ok := Mgr.dbs[m]
	if !ok {
		log.Println("Failed to get database from dbs manager.")
		return fmt.Errorf("Failed to get database from dbs manager")
	}

	err := db.CreatePool(MsgPool)
	if err != nil {
		log.Println("Not recreating pool:", err)
	}

	pool, err := db.GetPool(MsgPool)
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

func (Mgr *DbManager) GetMessages(m manifest.Manifest) ([]Message, error) {

	db, ok := Mgr.dbs[m]
	log.Println(Mgr.dbs)
	if !ok {
		log.Println("Failed to get database from dbs manager.")
		return nil, fmt.Errorf("Failed to get database from dbs manager 1")
	}

	err := db.CreatePool(MsgPool)
	if err != nil {
		log.Println("Not recreating pool:", err)
	}

	pool, err := db.GetPool(MsgPool)
	if err != nil {
		log.Println("Failed to get pool:", err)
		return nil, err
	}

	filter := map[string]interface{}{
		"datatype": MessageType,
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

	log.Println(msg_data)

	return msg_data, nil
}

func convertToMessages(data []map[string]interface{}) []Message {
	messages := make([]Message, len(data))

	for i, item := range data {
		// Assuming your map contains fields like "ID" and "Text"
		// Adjust these according to your actual map structure
		d_type, _ := item["datatype"].(int)
		sender, _ := item["sender"].(string)
		data, _ := item["data"].(string)
		time, _ := item["time"].(string)

		// Create a new Message and append it to the result slice
		messages[i] = Message{
			DataType: d_type,
			Sender:   sender,
			Data:     data,
			Time:     time,
		}
	}

	return messages
}
