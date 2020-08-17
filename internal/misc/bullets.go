package misc

import "math/rand"

var bullets = []string{":green_heart:", ":heart:", ":blue_heart:", ":yellow_heart:", ":purple_heart:", ":orange_heart:", ":black_heart:"}

// RandomBullet - returns a random exclamation
func RandomBullet() string {
	return bullets[rand.Intn(len(bullets))]
}
