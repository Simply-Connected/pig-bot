package pig_bot

import "math/rand"

func Duel(first *Pig, second *Pig) *Pig {
	totalWeight := float64(first.Weight + second.Weight)
	firstChance := float64(first.Weight) / totalWeight
	roll := rand.Float64()
	if roll <= firstChance {
		return first
	} else {
		return second
	}
}
