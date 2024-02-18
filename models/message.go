package models

//POOLS

const MsgPool = "message_pool"
const UserPool = "user_pool"
const MessageType = 1
const UserType = 2

type Message struct {
	DataType  int                `json:"datatype"`
	Sender    ProfileStorePublic `json:"sender"`
	SenderId  string             `json:"senderid"`
	Data      string             `json:"data"`
	Time      string             `json:"time"`
	Signature string             `json:"sign"`
	Valid     bool               `json:"valid"` // This is only for chat in db it always false
}

type User struct {
	NickName string
}

const MStore_TYPE_Manifet = 1

type MStore struct {
	Type int    `json:"type"`
	Data []byte `json:"data"`
}
