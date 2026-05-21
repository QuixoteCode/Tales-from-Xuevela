package main

import (
	"fmt"
	"time"
	"strings"
)

var name string
var strength uint8
var tenacity uint8
var luck uint8
const availableAttributePoints uint8 = 15

type Character struct {
	Name      string
	Strength  uint8
	Tenacity  uint8
	Luck      uint8
	Hitpoints    uint8
}

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

	player := Character{
		Name:     name,
		Strength: strength,
		Tenacity: tenacity,
		Luck:     luck,
		Hitpoints:   int(tenacity) * 5,
	}

	fmt.Println("Your adventure starts now, get ready...")

	// Waits three seconds to build suspense
	time.Sleep(3 * time.Second)

	fmt.Println("It's a sunny summer morning in what by all looks is the middle of a particularly overgrown fallow field, you are awaken from your slumber by a loud noise in the distance coming from the North")

	var choice string

	fmt.Println("Do you go North or South?")

	for {
		fmt.Scanln(&choice)

		// We convert the input to lowercase to avoid any issues
		choice = strings.ToLower(choice)

		if choice == "north" {
			decision1North()
			break
		} else if choice == "south" {
			decision1South()
			break
		} else {
			fmt.Println("You hesitate, unable to choose.\n")
			fmt.Println("Please enter either North or South.\n")
		}

	}
}

func decision1North() {
	fmt.Println("You head North towards the sound...")

	rat := Character{
		Name:     "Rat",
		Strength: 1,
		Tenacity: 2,
		Luck:     1,
		Health:   10,
	}

	fmt.Println("A giant rat jumps out of the grass! It looks angry...")
	
	combat(player, &rat)
}

func decision1South() {
	fmt.Println("You head South into the tall grass...")
}

// TODO combat
func combat(player *Character, enemy *Character) { }