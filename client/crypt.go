package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
)

func makeMultiple16Bit(plaintext []byte) []byte {
	str := string(plaintext)

	if len(str) < aes.BlockSize {
		left := aes.BlockSize - len(str)

		for i := 0; i < left; i++ {
			str = str + "|"
		}
		makeMultiple16Bit([]byte(str))
	}

	if len(str)%aes.BlockSize != 0 {
		left := len(str) % aes.BlockSize

		for i := 0; i < (aes.BlockSize - left); i++ {
			str = str + "|"
		}
		makeMultiple16Bit([]byte(str))
	}

	fmt.Println(str)

	return []byte(str)
}

// EncryptCBC ...
func EncryptCBC(key, plaintext []byte) (ciphertext []byte, err error) {
	plaintext = makeMultiple16Bit(plaintext)

	log.Println(len(plaintext))

	if len(plaintext)%aes.BlockSize != 0 {
		panic("plaintext is not a multiple of the block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext = make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	//iv, _ := hex.DecodeString("acfa7a047800b2f221f2c4f7d626eafb")
	//copy(ciphertext[:aes.BlockSize], iv)

	// fmt.Printf("CBC Key: %s\n", hex.EncodeToString(key))
	// fmt.Printf("CBC IV: %s\n", hex.EncodeToString(iv))

	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return
}

// DecryptCBC ...
func DecryptCBC(key, ciphertext []byte) (plaintext []byte, err error) {
	var block cipher.Block

	if block, err = aes.NewCipher(key); err != nil {
		return
	}

	if len(ciphertext) < aes.BlockSize {
		// fmt.Println("ciphertext too short")
		return
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(ciphertext, ciphertext)

	plaintext = ciphertext

	return
}
