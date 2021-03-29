package main

import (
	"math/rand"
	"os"
)

const MaxInt32 = int(^uint(0) >> 1)

func GetEnvVariable(key string, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func SliceContains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func RandomString(length int) string {
	var output string
	for i := 0; i < length; i++ {
		z := rand.Intn(3)
		switch z {
		case 0: // Numbers
			output += string(rune(rand.Intn(57-48) + 48))
		case 1: // Uppercase letters
			output += string(rune(rand.Intn(90-65) + 65))
		case 2: // Lowercase letters
			output += string(rune(rand.Intn(122-97) + 97))
		}
	}
	return output
}
