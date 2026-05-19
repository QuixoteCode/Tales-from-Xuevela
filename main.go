package main

import (
	“fmt”
	“time”
)

var name string
var strength uint8
var tenacity uint8
var luck uint8

// Checks whether the character attribute is 10 or lower
func getAttribute(statName string, prompt string) uint8 {
	var value int

	for {
		fmt.Print(prompt)
		fmt.Scanln(&value)

		if value >= 0 && value <= 10 {
			return uint8(value)
		}

		fmt.Printf("%s must be 10 or lower. Try again.\n", statName)
	}
}

func main() {
	fmt.Print("Enter your name:\n")
	fmt.Scanln(&name)
	fmt.Println("Hello,", name)

	strength = getAttribute("Strength", "Enter your strength:\n")
	fmt.Println("Your strength is:", strength)

	tenacity = getAttribute("Tenacity", "Enter your tenacity:\n")
	fmt.Println("Your tenacity is:", tenacity)

	luck = getAttribute("Luck", "Enter your luck:\n")
	fmt.Println("Your luck is:", luck)

	fmt.Println("Your adventure starts now, get ready...")

	// Waits three seconds to build suspense
	time.Sleep(3 * time.Second)

	fmt.Println("It's a sunny summer morning in what by all looks is the middle of a particularly overgrown fallow field, you are awaken from your slumber by a loud noise in the distance")
}