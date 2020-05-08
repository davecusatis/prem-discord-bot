package cardset

import (
	"fmt"
	"strconv"
	"strings"
)

var workoutMap map[string][]string

const colonDelim = ":"

// CardSet is the type we use to represent cards
type CardSet struct {
	SCount int
	HCount int
	DCount int
	CCount int
	Cards  []string
}

func (c *CardSet) removeCard(card string) []string {
	for i, ca := range c.Cards {
		if card == ca {
			return append(c.Cards[:i], c.Cards[i+1:]...)
		}
	}
	return c.Cards
}

// GetCardTranslation translates the card to some movement and reps or rest
func (c *CardSet) GetCardTranslation(card string) string {
	cardStrs := strings.Split(card, colonDelim)
	if len(cardStrs) != 2 {
		return ""
	}

	value := cardStrs[0]
	suit := cardStrs[1]

	// if ace rest 2 min
	if value == "A" {
		return "Rest 2 min"
	}

	// face cards = 20 reps
	var translation strings.Builder
	if value == "K" || value == "J" || value == "Q" {
		translation.WriteString("20 reps ")
	} else {
		// all other values are the value + 10
		i, err := strconv.Atoi(value)
		if err != nil {
			return ""
		}
		translation.WriteString(fmt.Sprintf("%d reps ", i+10))
	}

	// determine exercise and increment counter
	workouts := workoutMap[suit]
	switch suit {
	case "S":
		translation.WriteString(workouts[c.SCount%4])
		c.SCount = c.SCount + 1
	case "C":
		translation.WriteString(workouts[c.CCount%4])
		c.CCount = c.CCount + 1
	case "H":
		translation.WriteString(workouts[c.HCount%4])
		c.HCount = c.HCount + 1
	case "D":
		translation.WriteString(workouts[c.DCount%4])
		c.DCount = c.DCount + 1
	}

	// remove card so we don't have repeats
	c.Cards = c.removeCard(card)
	return translation.String()
}

// InitSet inits the card set
func InitSet() *CardSet {
	workoutMap = map[string][]string{
		"S": []string{
			"push up",
			"pike push up",
			"banded press",
			"lateral raises",
		},
		"C": []string{
			"banded pull down",
			"horizontal banded row",
			"banded rear delt flyes",
			"banded upright row",
		},
		"H": []string{
			"walking lunges",
			"bulgarian split squat",
			"single leg hip thrust",
			"banded romanian dead lift",
		},
		"D": []string{
			"banded curl",
			"banded skull crushers",
			"reverse crunches",
			"reverse crunches",
		},
	}
	return &CardSet{
		SCount: 0,
		HCount: 0,
		DCount: 0,
		CCount: 0,
		Cards: []string{
			"A:S",
			"2:S",
			"3:S",
			"4:S",
			"5:S",
			"6:S",
			"7:S",
			"8:S",
			"9:S",
			"10:S",
			"J:S",
			"Q:S",
			"K:S",
			"A:H",
			"2:H",
			"3:H",
			"4:H",
			"5:H",
			"6:H",
			"7:H",
			"8:H",
			"9:H",
			"10:H",
			"J:H",
			"Q:H",
			"K:H",
			"A:D",
			"2:D",
			"3:D",
			"4:D",
			"5:D",
			"6:D",
			"7:D",
			"8:D",
			"9:D",
			"10:D",
			"J:D",
			"Q:D",
			"K:D",
			"A:C",
			"2:C",
			"3:C",
			"4:C",
			"5:C",
			"6:C",
			"7:C",
			"8:C",
			"9:C",
			"10:C",
			"J:C",
			"Q:C",
			"K:C",
		},
	}
}
