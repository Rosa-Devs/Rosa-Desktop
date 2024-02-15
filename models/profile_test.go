package models

import (
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
		Sender:   profile.Name,
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
