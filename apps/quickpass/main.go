package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+"

func generatePassword(length int) string {
	b := make([]byte, length)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[num.Int64()]
	}
	return string(b)
}

func main() {
	fmt.Println("Generated 20-char secure passwords:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("  %d) %s\n", i, generatePassword(20))
	}
}
