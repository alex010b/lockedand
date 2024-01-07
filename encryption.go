package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
)

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func comparePasswordWithHash(password string, storedHash []byte) error {
	return bcrypt.CompareHashAndPassword(storedHash, []byte(password))
}

func encrypt(key string, plaintext string) []byte {
	salt := []byte("7*asdih*bs89db(743)")
	keyLen := 32 // AES-256
	iterations := 4096

	newKey := pbkdf2.Key([]byte(key), salt, iterations, keyLen, sha256.New)
	fmt.Printf("Generated key: %x\n", newKey)

	block, err := aes.NewCipher(newKey)
	if err != nil {
		fmt.Println("couldnt encrypt", err)
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Println("couldnt encrypt", err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	copy(ciphertext[:aes.BlockSize], iv)

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return ciphertext
}

func decrypt(key []byte, ciphertext []byte) (string, error) {

	password := key
	salt := []byte("7*asdih*bs89db(743)")
	keyLen := 32 // AES-256
	iterations := 4096

	newKey := pbkdf2.Key([]byte(password), salt, iterations, keyLen, sha256.New)

	block, err := aes.NewCipher(newKey)
	if err != nil {
		fmt.Println(err)
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)

	return string(plaintext), nil
}

func decryptFileWithKey(key []byte, filePath string) (string, error) {

	ciphertext, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	plaintext, err := decrypt(key, ciphertext)
	if err != nil {
		return "", err
	}

	return plaintext, nil
}
