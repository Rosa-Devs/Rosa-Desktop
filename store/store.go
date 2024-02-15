package store

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

// This module giving defaul store paths for app
type Store struct {
	Global string

	Profile  string
	Database string
}

func NewStore(custom string) (*Store, error) {
	var homePath string

	switch oss := runtime.GOOS; oss {
	case "darwin", "linux":
		homePath = os.Getenv("HOME")
	case "windows":
		homePath = os.Getenv("USERPROFILE")
	default:
		return nil, fmt.Errorf("unsupported operating system: %s", oss)
	}
	if homePath == "" {
		return nil, fmt.Errorf("unable to determine home folder for operating system!")
	}

	homePath = homePath + "/rosa"

	if custom != "" {
		homePath = custom + "/rosa"
	}

	if _, err := os.Stat(homePath); os.IsNotExist(err) {
		// Directory does not exist, create it
		err := os.Mkdir(homePath, os.FileMode(0775))
		if err != nil {
			log.Fatal("Error creating directory:", err)
		}
		log.Println("DataStoreFolder created with permissions 775")
	} else {
		log.Println("DataStoreFolder already exists")
	}

	return &Store{
		Global:   homePath,
		Profile:  homePath + "/profile.json",
		Database: homePath + "/database",
	}, nil
}
