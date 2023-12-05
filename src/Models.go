package src

//POOLS

const MsgPool = "message_pool"
const UserPool = "user_pool"
const MessageType = 1
const UserType = 2

type Message struct {
	DataType int    `json:"datatype"`
	Sender   string `json:"sender"`
	Data     string `json:"data"`
	Time     string `json:"time"`
}

type User struct {
	NickName string
}

const MStore_TYPE_Manifet = 1

type MStore struct {
	Type int    `json:"type"`
	Data []byte `json:"data"`
}
