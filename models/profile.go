package models

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	"github.com/Rosa-Devs/Database/src/manifest"
)

const Profile_type int = 2

type Profile struct {
	Name    string `json:"name"`
	Id      string `json:"id"`
	PrivKey string `json:"priv"`
	PubKey  string `json:"pub"`
	Avatar  string `json:"avatar"`
}

type ProfileStorePublic struct {
	Type   string `json:"type"`
	Name   string `json:"name"`
	Id     string `json:"id"`
	PubKey string `json:"pub"`
	Avatar string `json:"avatar"`
}

func (p *Profile) Serialize() ([]byte, error) {
	// Using JSON encoding for serialization
	data, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (p *Profile) Deserialize(s string) error {
	// Using JSON decoding for deserialization
	err := json.Unmarshal([]byte(s), p)
	if err != nil {
		return err
	}
	return nil
}

func (p *Profile) GetPublic() ProfileStorePublic {
	return ProfileStorePublic{
		Name:   p.Name,
		Id:     p.Id,
		Avatar: p.Avatar,
		PubKey: p.PubKey,
	}
}

func (p *ProfileStorePublic) Serialize() ([]byte, error) {
	// Using JSON encoding for serialization
	data, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (p *ProfileStorePublic) Deserialize(s string) error {
	// Using JSON decoding for deserialization
	err := json.Unmarshal([]byte(s), p)
	if err != nil {
		return err
	}
	return nil
}

func CreateProfile(name string, Avatar string) (Profile, error) {
	// Generate RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Println("Failed to create account")
		return Profile{}, fmt.Errorf("failed to generate private key: %v", err)
	}
	publicKey := &privateKey.PublicKey

	// Convert private key to PEM format
	privkey_bytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkey_bytes,
		},
	)
	privKeyBase64 := base64.StdEncoding.EncodeToString(privkey_pem)

	// Convert public key to PEM format
	pubkey_bytes := x509.MarshalPKCS1PublicKey(publicKey)
	pubkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: pubkey_bytes,
		},
	)
	pubKeyBase64 := base64.StdEncoding.EncodeToString(pubkey_pem)

	if name == "" {
		log.Println("Name is required")
		return Profile{}, fmt.Errorf("name is requered")
	}

	return Profile{
		Name:    name,
		Id:      manifest.GenerateNoise(64),
		PrivKey: privKeyBase64,
		PubKey:  pubKeyBase64,
		Avatar:  Avatar,
	}, nil
}

func (p *ProfileStorePublic) ValidateMsg(m Message) bool {
	pubKeyBytes, err := base64.StdEncoding.DecodeString(p.PubKey)
	if err != nil {
		fmt.Println("Error decoding public key:", err)
		return false
	}

	block, _ := pem.Decode(pubKeyBytes)
	if block == nil {
		return false
	}

	// Parse the public key
	pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		fmt.Println("Error parsing public key:", err)
		return false
	}

	// Decode base64-encoded signature
	signature, err := base64.StdEncoding.DecodeString(m.Signature)
	if err != nil {
		fmt.Println("Error decoding signature:", err)
		return false
	}

	hashed := sha256.Sum256([]byte(m.Data + m.Time + m.SenderId))
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		//fmt.Println("Signature verification failed:", err)
		return false
	}

	return true
}

func (p *Profile) Sign(m *Message) error {
	// Decode base64-encoded private key
	privKeyBytes, err := base64.StdEncoding.DecodeString(p.PrivKey)
	if err != nil {
		return fmt.Errorf("error decoding private key: %v", err)
	}

	block, _ := pem.Decode(privKeyBytes)
	if block == nil {
		return fmt.Errorf("fail to get block of key")
	}
	// Parse the private key
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("error parsing private key: %v", err)
	}

	// Sign the message
	hashed := sha256.Sum256([]byte(m.Data + m.Time + m.SenderId))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashed[:])
	if err != nil {
		return fmt.Errorf("error signing message: %v", err)
	}

	// Encode the signature to base64
	m.Signature = base64.StdEncoding.EncodeToString(signature)

	return nil
}

func LoadFromFile(file string) (Profile, error) {
	// Read the content of the file
	content, err := os.ReadFile(file)
	if err != nil {
		return Profile{}, fmt.Errorf("error reading file: %v", err)
	}

	// Unmarshal JSON content into a Profile struct
	var profile Profile
	err = json.Unmarshal(content, &profile)
	if err != nil {
		return Profile{}, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return profile, nil
}

func WriteToFile(file string, profile Profile) error {
	// Marshal the Profile struct to JSON
	jsonData, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling profile to JSON: %v", err)
	}

	// Write the JSON data to the file
	err = os.WriteFile(file, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}
