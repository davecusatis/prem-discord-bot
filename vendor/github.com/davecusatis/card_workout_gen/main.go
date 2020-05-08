package main

import (
	"math/rand"
	"time"
)

const workoutLength = 20

func main() {
	rand.Seed(time.Now().Unix())

	// cardSet := cardset.InitSet()
	// var randomCard string
	// var randIndex int
	// var output strings.Builder
	// for i := 0; i < workoutLength; i++ {
	// 	randIndex = rand.Intn(len(cardSet.Cards))
	// 	randomCard = cardSet.Cards[randIndex]
	// 	output.WriteString(fmt.Sprintf("%d. %s \n", i+1, cardSet.getCardTranslation(randomCard)))
	// }
	// fmt.Println(output.String())
}
