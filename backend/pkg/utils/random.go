package utils

import (
	crypto "crypto/rand"
	"encoding/base64"
	"math/rand"
	"strings"
	"time"
)

// rand.Read is deprecated: For almost all use cases, [crypto/rand.Read] is more appropriate.deprecateddefault
// func rand.Read(p []byte) (n int, err error)
// Read generates len(p) random bytes from the default [Source] and writes them into p. It always returns len(p) and a nil error. Read, unlike the [Rand.Read] method, is safe for concurrent use.

// Deprecated: For almost all use cases, crypto/rand.Read is more appropriate.

// Added in go1.6

func GenerateRandomString(length int) (string, error) {
	b := make([]byte, length)
	_, err := crypto.Read(b)

	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

var (
	// Simple list of common words for generating memorable passphrases
	wordList = []string{
		"apple", "banana", "cherry", "date", "elder", "fig", "grape", "honey",
		"indigo", "juniper", "kiwi", "lemon", "mango", "nectar", "orange", "peach",
		"quince", "raspberry", "strawberry", "tangerine", "umbrella", "vanilla", "watermelon",
		"xylophone", "yellow", "zucchini", "almond", "butter", "coffee", "donut",
		"egg", "flour", "garlic", "honey", "ice", "jam", "kale", "lentil",
		"milk", "nut", "oats", "pepper", "quinoa", "rice", "sugar", "tea",
		"vinegar", "water", "yogurt", "zest", "acorn", "butterfly", "cat", "dog",
		"elephant", "fox", "giraffe", "horse", "iguana", "jaguar", "kangaroo", "lion",
		"monkey", "newt", "octopus", "penguin", "quail", "rabbit", "snake", "tiger",
		"unicorn", "vulture", "whale", "xerox", "yak", "zebra", "airplane", "book",
		"canoe", "desk", "elevator", "flute", "guitar", "hat", "igloo", "jacket",
		"kite", "lamp", "mirror", "notebook", "oven", "pencil", "quilt", "radio",
		"scissors", "telephone", "umbrella", "violin", "watch", "xylophone", "yacht", "zipper",
	}
)

// Init ensures random number generation is seeded
func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// GenerateRandomWords generates a passphrase with the given number of random words
func GenerateRandomWords(count int) (string, error) {
	if count <= 0 {
		return "", nil
	}

	// Initialize a slice to hold selected words
	words := make([]string, count)

	// Select random words
	for i := 0; i < count; i++ {
		words[i] = wordList[rand.Intn(len(wordList))]
	}

	// Join words with hyphens to form a passphrase
	return strings.Join(words, "-"), nil
}
