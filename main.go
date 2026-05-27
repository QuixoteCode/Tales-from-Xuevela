package main

import (
	"fmt"
	"time"
	"strings"
	"math/rand"
)

// TODO make it so attributes can never be 0
var name string
var strength uint8
var tenacity uint8
var agility uint8
var luck uint8
const availableAttributePoints uint8 = 20

type Character struct {
	Name         string
	Strength     uint8
	Tenacity     uint8
	Agility      uint8
	Luck         uint8
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
		agility = getAttribute("Agility", "Enter your agility:\n")
		luck = getAttribute("Luck", "Enter your luck:\n")

		total := strength + tenacity + agility + luck

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
		Name:        name,
		Strength:    strength,
		Tenacity:    tenacity,
		Agility:     agility,
		Luck:        luck,
		Hitpoints:   uint8(tenacity) * 5,
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
			decision1North(&player)
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

func decision1North(player *Character) {
	fmt.Println("You head North towards the sound...")

	rat := Character{
		Name:      "Rat",
		Strength:  1,
		Tenacity:  2,
		Agility:   1,
		Luck:      1,
	}
	rat.Hitpoints = rat.Tenacity * 5

	fmt.Println("A giant rat jumps out of the grass! It looks angry...")

	combat(player, &rat)
}

func decision1South() {
	fmt.Println("You head South into the tall grass...")
}

func combat(player *Character, enemy *Character) {

	// Determine turn order (initiative) based on agility
	var first *Character
	var second *Character

	if player.Agility >= enemy.Agility {
		first = player
		second = enemy
	} else {
		first = enemy
		second = player
	}

	for player.Hitpoints > 0 && enemy.Hitpoints > 0 {

		// First character attacks
		fmt.Printf("%s attacks %s!\n", first.Name, second.Name)

		// TODO fix uint8 vs. int conflict
		var rollEvasionSecondCharacter uint8 = rand.Intn(100)
		var rollCriticalStrikeFirstCharacter uint8 = rand.Intn(100)

		if (second.Agility >= rollEvasionSecondCharacter) {
			fmt.Printf("%s has evaded %d's attack!\n", second.Name, first.Name)
		} else { 
			if first.Luck >= rollCriticalStrikeFirstCharacter {
				// Critical strike
				second.Hitpoints -= first.Strength * 3
				fmt.Printf("%s has struck %d critically!\n", first.Name, second.Name)
			} else {
				//Normal strike
				second.Hitpoints -= first.Strength
			}  
		}

		fmt.Printf("%s has %d HP left\n", second.Name, second.Hitpoints)

		if second.Hitpoints <= 0 {
			fmt.Printf("%s is defeated!\n", second.Name)
			break
		}

		// TODO apply same logic as when the first character attacks
		// Second character attacks
		fmt.Printf("%s attacks %s!\n", second.Name, first.Name)
		first.Hitpoints -= second.Strength

		fmt.Printf("%s has %d HP left\n", first.Name, first.Hitpoints)

		if first.Hitpoints <= 0 {
			fmt.Printf("%s is defeated!\n", first.Name)
			break
		}
	}
}