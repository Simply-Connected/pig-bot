package pig_bot

import "fmt"

func getPositiveGrowPhrase(pig *Pig, diff int32) string {
	return fmt.Sprintf("Ваш *%s* съел морковки и поправился на *%d* килограмм.\n\nВес вашего свина: *%d*", pig.Name, diff,
		pig.Weight)
}

func getNegativeGrowPhrase(pig *Pig, diff int32) string {
	return fmt.Sprintf("Вашему *%s* попался сгнивший корм и он потерял *%d* килограмм.\n\nВес вашего свина: *%d*", pig.Name, -diff,
		pig.Weight)
}

func getZeroGrowPhrase(pig *Pig, diff int32) string {
	return fmt.Sprintf("Ваш *%s* не нашел своего корма и лег спать.\n\nВес вашего свина: *%d*", pig.Name,
		int32(pig.Weight))
}

func getGrowPhrase(pig *Pig, diff int32) string {
	if diff > 0 {
		return getPositiveGrowPhrase(pig, diff)
	} else if diff == 0 {
		return getZeroGrowPhrase(pig, diff)
	} else {
		return getNegativeGrowPhrase(pig, diff)
	}
}
