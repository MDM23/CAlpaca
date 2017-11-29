package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log"
)

type User struct {
	ID  string
	Key string
}

type Chain struct {
	Payload string
	UserIDS []string
}

func main() {
	users := []User{
		User{"A", "SomeKeyOfA"},
		User{"B", "SomeKeyOfB"},
		User{"C", "SomeKeyOfC"},
		User{"D", "SomeKeyOfD"},
		User{"E", "SomeKeyOfE"},
	}

	// This specifies how many distinct users are needed to decrypt the real key
	consensusParam := 3

	chain := Chain{"This will hold the RSA private key data", []string{}}
	chains, err := makeChains(chain, users, consensusParam)

	if err != nil {
		log.Fatalln(err)
	}

	for _, c := range chains {
		fmt.Printf("Keychain %v provides key %s\n", c.UserIDS, c.Payload)
	}
}

func makeChains(chain Chain, users []User, chainLen int) ([]Chain, error) {
	var chains []Chain

	for i := 0; i < len(users); i++ {
		var uids []string
		uids = append(uids, users[i].ID)
		uids = append(uids, chain.UserIDS...)

		key := sha256.Sum256([]byte(users[i].Key))
		encryptedText, err := encrypt(key[:], chain.Payload)
		if err != nil {
			return nil, err
		}

		chain := Chain{encryptedText, uids}

		if chainLen == 1 {
			chains = append(chains, chain)
			continue
		}

		var others []User
		others = append(others, users[:i]...)
		others = append(others, users[i+1:]...)

		newChains, err := makeChains(chain, others, chainLen-1)
		if err != nil {
			return nil, err
		}

		chains = append(chains, newChains...)
	}

	return chains, nil
}

func encrypt(key []byte, message string) (encmess string, err error) {
	plainText := []byte(message)

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	encmess = base64.URLEncoding.EncodeToString(cipherText)
	return
}
