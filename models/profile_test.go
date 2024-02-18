package models

import (
	"encoding/json"
	"testing"
)

func TestProfileSignAndValidate(t *testing.T) {

	profile, err := CreateProfile("Mihalic2040", "")
	if err != nil {
		t.Error("Fail to creat profile")
	}
	profilePublic := profile.GetPublic()

	message := Message{
		DataType: 1,
		Sender:   profile.GetPublic(),
		SenderId: profile.Id,
		Data:     "Hello, World!",
		Time:     "2022-02-15T12:34:56Z",
	}

	// Test signing and validation
	err = profile.Sign(&message)
	if err != nil {
		t.Errorf("Error signing message: %v", err)
	}

	isValid := profilePublic.ValidateMsg(message)
	if !isValid {
		t.Errorf("Message signature is not valid.")
	}
}

func TestProfileSignAndValidateSerialized(t *testing.T) {
	// Create a profile
	profile, err := CreateProfile("Mihalic2040", "")
	if err != nil {
		t.Error("Fail to create profile")
	}

	// Get the public profile
	profilePublic := profile.GetPublic()

	// Create a message
	message := Message{
		DataType: 1,
		Sender:   profile.GetPublic(),
		SenderId: profile.Id,
		Data:     "Hello, World!",
		Time:     "2022-02-15T12:34:56Z",
	}

	// Test signing
	err = profile.Sign(&message)
	if err != nil {
		t.Errorf("Error signing message: %v", err)
	}

	// Serialize the message
	messageJSON, err := json.Marshal(message)
	if err != nil {
		t.Errorf("Error serializing message: %v", err)
	}

	// Deserialize the message
	var deserializedMessage Message
	err = json.Unmarshal(messageJSON, &deserializedMessage)
	if err != nil {
		t.Errorf("Error deserializing message: %v", err)
	}

	// Test validation
	isValid := profilePublic.ValidateMsg(deserializedMessage)
	if !isValid {
		t.Errorf("Message signature is not valid.")
	}
}
