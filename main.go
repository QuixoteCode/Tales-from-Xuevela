package main

import (
	"fmt"
	"time"
	"strings"
	"math/rand"
	"os"
	"bufio"
)

var name string
var strength uint8
var tenacity uint8
var agility uint8
var luck uint8
var charisma uint8
const availableAttributePoints uint8 = 25

type Character struct {
	Name                    string
	Strength                uint8
	Tenacity                uint8
	Agility                 uint8
	Luck                    uint8
	Charisma                uint8
	MaxHitpoints            int
	CurrentHitpoints        int
	Experience              int
	Level                   int
	ExperienceToNextLevel   int
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
		charisma = getAttribute("Charisma", "Enter your charisma:\n")

		total := strength + tenacity + agility + luck + charisma

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
		Name:                   name,
		Strength:               strength,
		Tenacity:               tenacity,
		Agility:                agility,
		Luck:                   luck,
		Charisma:               charisma,
		MaxHitpoints:           int(tenacity) * 5,
		CurrentHitpoints:       (int(tenacity) * 5) / 2,
		Experience:             0,
		Level:                  1,
		ExperienceToNextLevel:  50,
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
			decisionNorth(&player)
			break
		} else if choice == "south" {
			decisionSouth(&player)
			break
		} else {
			fmt.Println("You hesitate, unable to choose.\n")
			fmt.Println("Please enter either North or South.\n")
		}

	}

	if player.CurrentHitpoints <= 0 {
		fmt.Printf("%s is defeated!\n", player.Name)
		os.Exit(0)
	}
}

func decisionNorth(player *Character) {
	fmt.Println("You head North towards the sound...")

	time.Sleep(2 * time.Second)

	rat := Character{
		Name:      "Rat",
		Strength:  1,
		Tenacity:  2,
		Agility:   1,
		Luck:      1,
		Charisma:  1,
	}
	rat.MaxHitpoints = int(rat.Tenacity) * 5
	rat.CurrentHitpoints = (int(rat.Tenacity) * 5)

	fmt.Println("A giant rat jumps out of the grass! It looks angry...")

	time.Sleep(2 * time.Second)

	combat(player, &rat)
}

func decisionSouth(player *Character) {
	fmt.Println("You head South into the tall grass...")

	time.Sleep(2 * time.Second)

	xavier := Character{
		Name:      "Xavier",
		Strength:  4,
		Tenacity:  4,
		Agility:   4,
		Luck:      4,
		Charisma:  2,
	}

	fmt.Println("You find a tall bald man, almost as big as the cow next to him, he's resting next to a water faucet. The bovine is slowly but firmly carrying a plow with its yoke, tilling a patch of land. The bald man notices your presence and looks at you, wide-eyed, clearly not recognising you")
	
	time.Sleep(time.Second)
	
	fmt.Println("Do you...?: \n 1. Ask where are we \n 2. Say nothing")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		choice := strings.ToLower(strings.TrimSpace(scanner.Text()))
	
		if choice == "1" || choice == "ask where are we" {

			askXavierWhereWeAre()

			fmt.Println("Do you...?: \n 1. Ask for directions \n 2. Ask if you can drink from the faucet")

			for {
				scanner.Scan()
				choiceA := strings.ToLower(strings.TrimSpace(scanner.Text()))

				if choiceA == "1" || choiceA == "ask for directions" {
					
					askXavierForDirections()

					break

				} else if choiceA == "2" || choiceA == "ask if you can drink from the faucet" {

					askXavierIfYouCanDrinkFromTheFaucet(player, &xavier)

					break
					
				} else {
					fmt.Println("You hesitate, unable to choose.\n")
					fmt.Println("Please enter either \"1\" / \"ask for directions\" or \"2\" / \"ask if you can drink from the faucet\".\n")
				}

			}

			break

		} else if choice == "2" || choice == "say nothing" {

			sayNothingToXavier()

			fmt.Println("Do you...?: \n 1. Ask for directions \n 2. Ask if you can drink from the faucet \n 3. ask where are we \n 4. Continue to say nothing")

			sayNothingLoop:
				for {
					scanner.Scan()
					choiceB := strings.ToLower(strings.TrimSpace(scanner.Text()))
			
					switch choiceB {
					case "1", "ask for directions":
						askXavierForDirections()
						break sayNothingLoop
			
					case "2", "ask if you can drink from the faucet":
						askXavierIfYouCanDrinkFromTheFaucet(player, &xavier)
						break sayNothingLoop
			
					case "3", "ask where are we":
						askXavierWhereWeAre()
						break sayNothingLoop
			
					case "4", "continue to say nothing":
						continueToSayNothingToXavier()
						break sayNothingLoop
			
					default:
						fmt.Println("You hesitate, unable to choose.\n")
						fmt.Println("Please enter either \"1\" / \"ask for directions\", \"2\" / \"ask if you can drink from the faucet\", \"3\" / \"ask where are we\" or \"4\" / \"continue to say nothing\".\n")
					}
				}

			break

		} else {

			fmt.Println("You hesitate, unable to choose.\n")
			fmt.Println("Please enter either \"1\" / \"ask where are we\" or \"2\" / \"say nothing\".\n")

		}
	}

	if game.learnedMudLocationFromXavier {
		fmt.Println("Do you...?: \n 1. Continue South for Mud \n 2. Continue into a forest that borders the fields")
		scanner.Scan()
		choiceDecision := strings.ToLower(strings.TrimSpace(scanner.Text()))
		for {
			if choiceDecision == "1" || choiceDecision == "continue south for mud" {
				decisionSouthDecisionSouth()
				break
			} else if choiceDecision == "2" || choiceDecision == "continue into a forest that borders the fields" {
				decisionSouthDecisionForest()
				break
			}
		}
	} else {
		fmt.Println("Do you...?: \n 1. Continue South \n 2. Continue into a forest that borders the fields")
		scanner.Scan()
		choiceDecision := strings.ToLower(strings.TrimSpace(scanner.Text()))
		for {
			if choiceDecision == "1" || choiceDecision == "continue south" {
				decisionSouthDecisionSouth()
				break
			} else if choiceDecision == "2" || choiceDecision == "continue into a forest that borders the fields" {
				decisionSouthDecisionForest()
			}
		}
	}
}

func decisionSouthDecisionSouth() {
	time.Sleep(time.Second)

	// TODO change this println if the player knows about Mud
	fmt.Println("You decide to go further South")
}

func decisionSouthDecisionForest() {
	time.Sleep(time.Second)

	fmt.Println("You penetrate the thick vegatation")
}

func askXavierWhereWeAre() {
	fmt.Println("You inquire about this place")

	time.Sleep(time.Second)
	
	fmt.Println("\"Have you struck your head with a rock or something?\" You do not know how to react towards this valid hypothesis about your current situation dressed up as a quip, you feel a sharp pain in the back of your head. \"This is the small town of Mud, part of the parish of Soulstar\"")

	time.Sleep(time.Second)
}

type GameState struct {
    learnedMudLocationFromXavier bool
}

var game GameState

func askXavierForDirections() {
	fmt.Println("You ask Xavier for directions")

	time.Sleep(time.Second)

	fmt.Println("\"Mud is South of here if that's what you're asking\" he pauses for a second, \"so keep traveling in the direction you were going and you should be good\"")

	game.learnedMudLocationFromXavier = true

	time.Sleep(time.Second)
}

func askXavierIfYouCanDrinkFromTheFaucet(player *Character, xavier *Character) {
	fmt.Println("\"Huh, I would let you... but I'm worried about the cursed, you know?\"")

	time.Sleep(time.Second)

	resultConversationalChallengeXavier := conversationalChallenge(player, xavier)
	
	switch resultConversationalChallengeXavier {
		// complete success
		case 2:
			fmt.Println("The man allows you to drink; you take a refreshing sip. You have regenerated 5 hitpoints!")
			player.CurrentHitpoints += 5
	
		// faux pas
		case 1:
			fmt.Println("The man hesitates but, after a while, allows you to drink; you take a refreshing sip. You have regenerated 5 hitpoints!")
			player.CurrentHitpoints += 5
	
		// no healing
		case 0:
			fmt.Println("\"Sorry fella, can't do. \"Better safe than sorry\" as they say\"")
	}

	time.Sleep(time.Second)
}

func sayNothingToXavier() {
	fmt.Println("You say nothing")

	time.Sleep(time.Second)

	fmt.Println("The man coughs, not knowing where to look exactly. So, huh... what leads you to this place?")

	time.Sleep(time.Second)
}

func continueToSayNothingToXavier() {
	fmt.Println("You continue to say nothing")

	time.Sleep(time.Second)

	fmt.Println("The man in front of you decides to shift his focus away from you, continuing with what he was doing previous to your presence in that place")

	time.Sleep(time.Second)
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

	for player.CurrentHitpoints > 0 && enemy.CurrentHitpoints > 0 {

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
				second.CurrentHitpoints -= int((first.Strength + randomAdditionalDamageFirstAttack) * 3)
				fmt.Printf("%s has struck %s critically!\n", first.Name, second.Name)
			} else {
				//Normal strike
				second.CurrentHitpoints -= int(first.Strength + randomAdditionalDamageFirstAttack)
			}  
		}

		fmt.Printf("%s has %d hitpoints left\n", second.Name, second.CurrentHitpoints)

		if second.CurrentHitpoints <= 0 {
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
				first.CurrentHitpoints -= int((second.Strength + randomAdditionalDamageSecondAttack) * 3)
				fmt.Printf("%s has struck %s critically!\n", second.Name, first.Name)
			} else {
				//Normal strike
				first.CurrentHitpoints -= int(second.Strength + randomAdditionalDamageSecondAttack)
			}
		}

		fmt.Printf("%s has %d hitpoints left\n", first.Name, first.CurrentHitpoints)

		if first.CurrentHitpoints <= 0 {
			fmt.Printf("%s is defeated!\n", first.Name)
			break
		}

		time.Sleep(time.Second)
	}

	player.addExperience((int(enemy.Tenacity) * 10) + (int(enemy.Strength) * 10))
}

/*
	If it returns "2" the player has succeeded 
	If it returns "1" the player has performed somewhat acceptably but commited a faux pas
	If it returns "0" the player has completely failed
*/
func conversationalChallenge(player *Character, enemy *Character) uint8 {

	var charismaRoll = rand.Intn(100)
	threshold := int(player.Charisma) - int(enemy.Charisma)/2

	var fauxPasRoll = rand.Intn(100)

	time.Sleep(time.Second)

	fmt.Println("You have engaged in a dialectical encounter against", enemy.Name)

	time.Sleep(time.Second)

    // Complete success
    if charismaRoll < threshold {
		fmt.Println("You stand on the shoulders of giants and you have managed to convince your interlocutor")
		time.Sleep(time.Second)
		player.addExperience(int(enemy.Charisma) * 10)
        return 2
    }

    // The player's luck and agility allows him to savage the situation and make a faux pas instead of utterly failing
    salvage := int(player.Luck) + int(player.Agility)
    if fauxPasRoll < salvage {
		fmt.Println("You have awkwardly fumbled your way into a faux pas, it could be worse, but it certainly could be better")
		player.addExperience(int(enemy.Charisma) * 5)
		// Longer wait for further awkwardness
		time.Sleep(3 * time.Second)
        return 1
    }

    // Complete failure
	fmt.Println("You utterly failed at trying to convince your interlocutor of anything, you may as well talk to a wall")
	time.Sleep(time.Second)
    return 0

}

func (player *Character) addExperience(ExperienceToAdd int) {
	time.Sleep(time.Second)

	// Adding experience
	player.Experience += ExperienceToAdd
	fmt.Printf("You receive %d experience points\n", ExperienceToAdd)

	time.Sleep(time.Second)

	// Leveling up and the required experience for next level up
	for player.Experience >= player.ExperienceToNextLevel {
		player.Experience -= player.ExperienceToNextLevel
		player.Level++
		player.ExperienceToNextLevel = 50 * player.Level
		fmt.Println("You have leveled up!")
		time.Sleep(time.Second)
		fmt.Println("You are now level", player.Level, "and you need", player.ExperienceToNextLevel, "experience to level up again")
		time.Sleep(time.Second)
	}
}