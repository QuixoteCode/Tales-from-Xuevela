package main

import (
	"fmt"
	"time"
	"strings"
	"math/rand"
	"os"
)

var name string
var strength uint8
var tenacity uint8
var agility uint8
var luck uint8
const availableAttributePoints uint8 = 20

// TODO differenciate between max hitpoints and current hitpoints
type Character struct {
	Name         string
	Strength     uint8
	Tenacity     uint8
	Agility      uint8
	Luck         uint8
	Hitpoints    int
}

// Checks whether the character attribute is 10 or lower and bigger than 0
func getAttribute(statName string, prompt string) uint8 {
	var value int

	for {
		fmt.Print(prompt)
		fmt.Scanln(&value)

		switch {
			case value >= 0 && value <= 10:
				return uint8(value)

			case value > 10:
				fmt.Printf("%s must be 10 or lower. Try again.\n", statName)

			case value < 0:
				fmt.Printf("%s cannot be lower than 0. Try again.\n", statName)
		}
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

	if player.Hitpoints <= 0 {
		fmt.Printf("%s is defeated!\n", player.Name)
		os.Exit(0)
	}
}

func decision1North(player *Character) {
	fmt.Println("You head North towards the sound...")

	time.Sleep(2 * time.Second)

	rat := Character{
		Name:      "Rat",
		Strength:  1,
		Tenacity:  2,
		Agility:   1,
		Luck:      1,
	}
	rat.Hitpoints = int(rat.Tenacity) * 5

	fmt.Println("A giant rat jumps out of the grass! It looks angry...")

	time.Sleep(2 * time.Second)

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

		var rollEvasionSecondCharacter uint8 = uint8(rand.Intn(100))
		var rollCriticalStrikeFirstCharacter uint8 = uint8(rand.Intn(100))
		var randomAdditionalDamageFirstAttack uint8 = uint8(rand.Intn(2))

		if (second.Agility >= rollEvasionSecondCharacter) {
			fmt.Printf("%s has evaded %s's attack!\n", second.Name, first.Name)
		} else { 
			if first.Luck >= rollCriticalStrikeFirstCharacter {
				// Critical strike
				second.Hitpoints -= int((first.Strength + randomAdditionalDamageFirstAttack) * 3)
				fmt.Printf("%s has struck %s critically!\n", first.Name, second.Name)
			} else {
				//Normal strike
				second.Hitpoints -= int(first.Strength + randomAdditionalDamageFirstAttack)
			}  
		}

		fmt.Printf("%s has %d HP left\n", second.Name, second.Hitpoints)

		if second.Hitpoints <= 0 {
			fmt.Printf("%s is defeated!\n", second.Name)
			break
		}

		time.Sleep(time.Second)

		// Second character attacks
		fmt.Printf("%s attacks %s!\n", second.Name, first.Name)

		var rollEvasionFirstCharacter uint8 = uint8(rand.Intn(100))
		var rollCriticalStrikeSecondCharacter uint8 = uint8(rand.Intn(100))
		var randomAdditionalDamageSecondAttack uint8 = uint8(rand.Intn(2))

		if (first.Agility >= rollEvasionFirstCharacter) {
			fmt.Printf("%s has evaded %s's attack!\n", first.Name, second.Name)
		} else {
			if second.Luck >= rollCriticalStrikeSecondCharacter {
				// Critical strike
				first.Hitpoints -= int((second.Strength + randomAdditionalDamageSecondAttack) * 3)
				fmt.Printf("%s has struck %s critically!\n", second.Name, first.Name)
			} else {
				//Normal strike
				first.Hitpoints -= int(second.Strength + randomAdditionalDamageSecondAttack)
			}
		}

		fmt.Printf("%s has %d HP left\n", first.Name, first.Hitpoints)

		if first.Hitpoints <= 0 {
			fmt.Printf("%s is defeated!\n", first.Name)
			break
		}

		time.Sleep(time.Second)
	}
}