package misc

import "math/rand"

var exclamations = []string{"What an absolute unit!", "Bravo! :clap:", "What a trooper!", ":tada:", "Look at you go!", "So classy! :face_with_monocle:"}

// RandomExclamation - returns a random exclamation
func RandomExclamation() string {
	return exclamations[rand.Intn(len(exclamations))]
}
