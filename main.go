package main

import (
	"fmt"
	"time"
)

var name string
var strength uint8
var tenacity uint8
var luck uint8
const availableAttributePoints uint8 = 15

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

	// Keep asking until total points are valid
	for {
		strength = getAttribute("Strength", "Enter your strength:\n")
		tenacity = getAttribute("Tenacity", "Enter your tenacity:\n")
		luck = getAttribute("Luck", "Enter your luck:\n")

		total := strength + tenacity + luck

		fmt.Printf("\nTotal attribute points used: %d/%d\n",
			total, availableAttributePoints)

		if total == availableAttributePoints {
			break
		}

		if total > availableAttributePoints {
			fmt.Printf(
				"You used too many attribute points! Maximum allowed is %d.\nPlease re-enter your attributes.\n\n",
				availableAttributePoints,
			)
		}

		if total < availableAttributePoints {
			fmt.Printf(
				"You didn't asign enough attribute points! The quantity allowed is %d.\nPlease re-enter your attributes.\n\n",
				availableAttributePoints,
			)
		}
		
	}

	fmt.Println("Your adventure starts now, get ready...")

	// Waits three seconds to build suspense
	time.Sleep(3 * time.Second)

	fmt.Println("It's a sunny summer morning in what by all looks is the middle of a particularly overgrown fallow field, you are awaken from your slumber by a loud noise in the distance")

	var choice string

	fmt.Println("Do you go left or right?")
	fmt.Scanln(&choice)

	if choice == "left" {
		decision1Left()
	} else if choice == "right" {
		decision1Right()
	} else {
		// TODO Re-ask the player
		fmt.Println("You hesitate, unable to choose.")
	}

}

func decision1Left() {
	fmt.Println("You head left toward the sound...")
}

func decision1Right() {
	fmt.Println("You head right into the tall grass...")
}