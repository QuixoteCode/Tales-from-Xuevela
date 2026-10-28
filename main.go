package main

import (
	"fmt"
	"time"
	"strings"
	"math/rand"
	"os"
	"bufio"
	"runtime"
	"Tales-from-Xuevela-main/display"
)

type WeaponType int

const (
	Melee WeaponType = iota
	Ranged
)

type WeaponSubtype int

const (
	// Melee subtypes
	Bladed WeaponSubtype = iota
	Blunt

	// Ranged subtypes
	Shooting
	Throwing
)

type Weapon struct {
	Name    string
	Damage  uint8
	Type    WeaponType
	Subtype WeaponSubtype
}

func (w *Weapon) Valid() bool {
	switch w.Type {
	case Melee:
		return w.Subtype == Bladed || w.Subtype == Blunt
	case Ranged:
		return w.Subtype == Shooting || w.Subtype == Throwing
	}
	return false
}

func NewWeapon(name string, damage uint8, t WeaponType, s WeaponSubtype) (Weapon, error) {
	w := Weapon{
		Name:    name,
		Damage:  damage,
		Type:    t,
		Subtype: s,
	}

	if !w.Valid() {
		return Weapon{}, fmt.Errorf(
			"invalid weapon: type %v with subtype %v",t,s,
		)
	}

	return w, nil
}

var playerName string
var playerPossesivePronoun string
var playerStrength uint8
var playerDexterity uint8
var playerTenacity uint8
var playerAgility uint8
var playerLuck uint8
var playerCharisma uint8
const availableAttributePoints uint8 = 30

type Character struct {
	Name                    string
	PossesivePronoun        string
	Strength                uint8
	Dexterity               uint8
	Tenacity                uint8
	Agility                 uint8
	Luck                    uint8
	Charisma                uint8
	MaxHitpoints            int
	CurrentHitpoints        int
	Experience              int
	Level                   int
	ExperienceToNextLevel   int
    MeleeWeapon             Weapon
    RangedWeapon            Weapon
	Inventory               Inventory
}

func (w Weapon) weaponIsEmpty() bool {
    return w.Name == "" && w.Damage == 0
}

// TODO this function seems to make the window be forever stuck
func getAttribute(statName string, prompt string) uint8 {
	var value int

	// Checks whether the character attribute is 10 or lower and 0 or higher
	for {
		// TODO change for window-mac.go's functions
		display.Printdisplay(prompt)
		fmt.Scanln(&value)

		switch {
			case value >= 0 && value <= 10:
				return uint8(value)

			case value > 10:
				display.Printdisplay(fmt.Sprintf("%s must be 10 or lower. Try again.\n", statName))

			case value < 0:
				display.Printdisplay(fmt.Sprintf("%s cannot be lower than 0. Try again.\n", statName))
		}
	}
}

type Item struct {
	Name     string
	Quantity int
}

type Inventory struct {
	Items []Item
}

// It increases the quantity of an existing item or adds a new item to the inventory if it doesnt exist
func (inventory *Inventory) Add(name string, quantity int) {
	for i := range inventory.Items {
		if inventory.Items[i].Name == name {
			inventory.Items[i].Quantity += quantity
			return
		}
	}

	inventory.Items = append(inventory.Items, Item{
		Name:     name,
		Quantity: quantity,
	})
}

// It removes a specified quantity of an item and returns true if successful or false if the item is missing or there arent enough units
func (inventory *Inventory) Remove(name string, quantity int) bool {
	for i := range inventory.Items {
		if inventory.Items[i].Name == name {
			if inventory.Items[i].Quantity < quantity {
				return false
			}

			inventory.Items[i].Quantity -= quantity

			if inventory.Items[i].Quantity == 0 {
				inventory.Items = append(
					inventory.Items[:i],
					inventory.Items[i+1:]...,
				)
			}

			return true
		}
	}

	return false
}

// Searches the inventory for an item by name
func (inventory *Inventory) Find(name string) *Item {
	for i := range inventory.Items {
		if inventory.Items[i].Name == name {
			return &inventory.Items[i]
		}
	}

	return nil
}

func main() {
    // Keep this goroutine on the main OS thread.
    // Cocoa UI initialization should happen on that thread.
    runtime.LockOSThread()

    // Create the native window.
    display.OpenWindow()
	
	display.Printdisplay("Enter your name:")
	playerName = display.Input("Enter your name:")
	display.Printdisplay("Hello, " + playerName)

	display.Printdisplay("Are you male or female:")
	var selectedGender string

	// Keep asking until the user responds male or female
	for {
		selectedGender = strings.ToLower(display.Input("Are you male or female:"))

		if selectedGender == "male" {
			playerPossesivePronoun = "his"
			break
		} else if selectedGender == "female" {
			playerPossesivePronoun = "her"
			break
		}

		display.Printdisplay("Please, introduce either \"male\" or \"female\":")
	}

	// Keep asking until total points are valid
	for {
		playerStrength = getAttribute("Strength", "Enter your strength:")
		playerDexterity = getAttribute("Dexterity", "Enter your dexterity:")
		playerTenacity = getAttribute("Tenacity", "Enter your tenacity:")
		playerAgility = getAttribute("Agility", "Enter your agility:")
		playerLuck = getAttribute("Luck", "Enter your luck:")
		playerCharisma = getAttribute("Charisma", "Enter your charisma:")

		total := playerStrength + playerDexterity + playerTenacity + playerAgility + playerLuck + playerCharisma

		display.Printdisplay(fmt.Sprintf("\nTotal attribute points used: %d/%d\n",
			total, availableAttributePoints))

		if total == availableAttributePoints {
			break
		}

		if total > availableAttributePoints {
			display.Printdisplay(fmt.Sprintf(
				"You used too many attribute points! Maximum allowed is %d.\nPlease re-enter your attributes.\n\n",
				availableAttributePoints,
			))
		}

		if total < availableAttributePoints {
			display.Printdisplay(fmt.Sprintf(
				"You didn't asign enough attribute points! The quantity allowed is %d.\nPlease re-enter your attributes.\n\n",
				availableAttributePoints,
			))
		}
	}

	defaultMeleeWeapon, err := NewWeapon(
		"Stick",
		playerStrength+1,
		Melee,
		Blunt,
	)
	if err != nil {
		display.Printdisplay("Could not create default melee weapon: " + err.Error())
		return
	}

	defaultRangedWeapon, err := NewWeapon(
		"Pebbles",
		playerDexterity+1,
		Ranged,
		Throwing,
	)
	if err != nil {
		display.Printdisplay("Could not create default ranged weapon: " + err.Error())
		return
	}

	player := Character{
		Name:                    playerName,
		PossesivePronoun:        playerPossesivePronoun,
		Strength:                playerStrength,
		Dexterity:               playerDexterity,
		Tenacity:                playerTenacity,
		Agility:                 playerAgility,
		Luck:                    playerLuck,
		Charisma:                playerCharisma,
		MaxHitpoints:            int(playerTenacity) * 5,
		CurrentHitpoints:        (int(playerTenacity) * 5) / 2,
		Experience:              0,
		Level:                   1,
		ExperienceToNextLevel:   50,
		MeleeWeapon:             defaultMeleeWeapon,
		RangedWeapon:            defaultRangedWeapon,
	}

	display.Printdisplay("Your adventure starts now, get ready...")

	// Waits three seconds to build suspense
	time.Sleep(3 * time.Second)

	display.Printdisplay("It's a sunny summer morning in what by all looks is the middle of a particularly overgrown fallow field, you are awaken from your slumber by a loud noise in the distance coming from the North")

	var choice string
	display.Printdisplay("Do you go North or South?")

	for {
		// We convert the input to lowercase to avoid any issues
		choice = strings.ToLower(display.Input("Do you go North or South?"))

		if choice == "north" {
			decisionNorth(&player)
			break
		} else if choice == "south" {
			decisionSouth(&player)
			break
		} else {
			display.Printdisplay("You hesitate, unable to choose.\n")
			display.Printdisplay("Please enter either North or South.\n")
		}
	}

	if player.CurrentHitpoints <= 0 {
		display.Printdisplay(fmt.Sprintf("%s is defeated!\n", player.Name))
		os.Exit(0)
	}

	// Start Cocoa's event loop to keep the window responsive.
	display.Run()
}

func decisionNorth(player *Character) {
	const ratStrength = 1

	// Used by the rat
	claws, err := NewWeapon(
		"Claws",
		ratStrength,
		Melee,
		Bladed,
	)
	if err != nil {
		// Make sure to keep it as a fmt.Println() and not replace it with display.Printdisplay()
		fmt.Println("Could not create rat's claws:", err)
		return
	}

	rat := Character{
		Name:             "Rat",
		PossesivePronoun: "its",
		Strength:         ratStrength,
		Tenacity:         2,
		Agility:          1,
		Luck:             1,
		Charisma:         1,
		MeleeWeapon:      claws,
	}
	rat.MaxHitpoints = int(rat.Tenacity) * 5
	rat.CurrentHitpoints = (int(rat.Tenacity) * 5)

	display.Printdisplay("You head North towards the sound...")

	time.Sleep(2 * time.Second)

	display.Printdisplay("A giant rat jumps out of the grass! It looks angry...")

	time.Sleep(2 * time.Second)

	combat(player, &rat, 2)

	time.Sleep(2 * time.Second)

	display.Printdisplay("Phew! That was close")

	time.Sleep(2 * time.Second)

	display.Printdisplay("In front of you, a dense forest raises, the rays of the sun struggle to pass through the unforgiving mass of leaves")

	time.Sleep(2 * time.Second)

	display.Printdisplay("Something on the floor catches a glimpse of your eye. You bend your back to crouch down and look at it more closely")

	time.Sleep(2 * time.Second)

	// Vara = roughly 0.8359 meters, 12 varas = roughly 10 meters, at 6 meters or taller virtually all falls are dangerous (no matter how well one lands), this amount of rope will allow the player to realistically get over a dangerous fall
	display.Printdisplay("It is a reel of rope, you feel like it will be useful later on, it has been added to your inventory, it seems to be roughly 12 varas long")

	player.Inventory.Add("Rope", 12)

	time.Sleep(2 * time.Second)

	display.Printdisplay("You continue walking and the forest gets denser and denser, you eventually reach a small clearing, there you perceive a well-trodden path through the grass which is surrounded by broken crates and barrels deformed as if they had imploded and all manner of junk")

	time.Sleep(2 * time.Second)

	display.Printdisplay("Do you continue walking in the same direction or deviate towards the well-trodden path?")

	// Needed because of multiple words like in "ask where are we", fmt.Scanln only reads a single value
	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		choice := strings.ToLower(strings.TrimSpace(scanner.Text()))
	
		if choice == "1" || choice == "continue walking in the same direction" {

			decisionNorthSameDirection()

			break

		} else if choice == "2" || choice == "deviate towards the well-trodden path" {

			decisionNorthPath()

			break

		} else {

			display.Printdisplay("You hesitate, unable to choose.\n")
			display.Printdisplay("Please enter either \"1\" / \"continue walking in the same direction\" or \"2\" / \"deviate towards the well-trodden path\".\n")

		}
	}

}


func decisionNorthSameDirection() {

}

func decisionNorthPath() {

}

func decisionSouth(player *Character) {
	const xavierStrength = 4

	// Used by Xavier
	hoe, err := NewWeapon(
		"Hoe",
		xavierStrength + 1,
		Melee,
		Blunt,
	)
	if err != nil {
		// Make sure to keep it as a fmt.Println() and not replace it with display.Printdisplay()
		fmt.Println("Could not create Xavier's hoe:", err)
		return
	}

	xavier := Character{
		Name:             "Xavier",
		PossesivePronoun: "his",
		Strength:         4,
		Dexterity:        4,
		Tenacity:         4,
		Agility:          4,
		Luck:             4,
		Charisma:         2,
		MeleeWeapon:      hoe,
	}

	display.Printdisplay("You head South into the tall grass...")

	time.Sleep(2 * time.Second)

	display.Printdisplay("You find a tall bald man, almost as big as the cow next to him, he's resting next to a water faucet. The bovine is slowly but firmly carrying a plow with its yoke, tilling a patch of land. The bald man notices your presence and looks at you, wide-eyed, clearly not recognising you")
	
	time.Sleep(time.Second)
	
	display.Printdisplay("Do you...?: \n 1. Ask where are we \n 2. Say nothing")

	// Needed because of multiple words like in "ask where are we", fmt.Scanln only reads a single value
	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		choice := strings.ToLower(strings.TrimSpace(scanner.Text()))
	
		if choice == "1" || choice == "ask where are we" {

			askXavierWhereWeAre()

			display.Printdisplay("Do you...?: \n 1. Ask for directions \n 2. Ask if you can drink from the faucet")

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
					display.Printdisplay("You hesitate, unable to choose.\n")
					display.Printdisplay("Please enter either \"1\" / \"ask for directions\" or \"2\" / \"ask if you can drink from the faucet\".\n")
				}

			}

			break

		} else if choice == "2" || choice == "say nothing" {

			sayNothingToXavier()

			display.Printdisplay("Do you...?: \n 1. Ask for directions \n 2. Ask if you can drink from the faucet \n 3. ask where are we \n 4. Continue to say nothing")

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
						display.Printdisplay("You hesitate, unable to choose.\n")
						display.Printdisplay("Please enter either \"1\" / \"ask for directions\", \"2\" / \"ask if you can drink from the faucet\", \"3\" / \"ask where are we\" or \"4\" / \"continue to say nothing\".\n")
					}
				}

			break

		} else {

			display.Printdisplay("You hesitate, unable to choose.\n")
			display.Printdisplay("Please enter either \"1\" / \"ask where are we\" or \"2\" / \"say nothing\".\n")

		}
	}

	if game.learnedMudLocation {
		display.Printdisplay("Do you...?: \n 1. Continue South for Mud \n 2. Continue into a forest that borders the fields")
		scanner.Scan()
		choiceDecision := strings.ToLower(strings.TrimSpace(scanner.Text()))
		for {
			if choiceDecision == "1" || choiceDecision == "continue south for mud" {
				decisionSouthDecisionSouth(player)
				break
			} else if choiceDecision == "2" || choiceDecision == "continue into a forest that borders the fields" {
				decisionSouthDecisionForest()
				break
			}
		}
	} else {
		display.Printdisplay("Do you...?: \n 1. Continue South \n 2. Continue into a forest that borders the fields")
		scanner.Scan()
		choiceDecision := strings.ToLower(strings.TrimSpace(scanner.Text()))
		for {
			if choiceDecision == "1" || choiceDecision == "continue south" {
				decisionSouthDecisionSouth(player)
				break
			} else if choiceDecision == "2" || choiceDecision == "continue into a forest that borders the fields" {
				decisionSouthDecisionForest()
			}
		}
	}
}

func decisionSouthDecisionSouth(player *Character) {
	const ravenStrength = 3

	// Used by Raven
	ironDagger, err := NewWeapon(
		"Iron Dagger",
		ravenStrength + 2,
		Melee,
		Bladed,
	)
	if err != nil {
		// Make sure to keep it as a fmt.Println() and not replace it with display.Printdisplay()
		fmt.Println("Could not create Raven's iron dagger:", err)
		return
	}

	raven := Character{
		Name:             "Raven",
		PossesivePronoun: "her",
		Strength:         3,
		Dexterity:        3,
		Tenacity:         3,
		Agility:          3,
		Luck:             3,
		Charisma:         3,
		MeleeWeapon:      ironDagger,
	}

	if game.learnedMudLocation {
		display.Printdisplay("You decide to go further South towards Mud")

		time.Sleep(2 * time.Second)

		display.Printdisplay("After much walking you find a cluster of buildings, they are grouped in small sets, forming streets inbetween them with what seems to be a townsquare in the center of it all. You assume this is Mud and a sign soon confirms your suspicions")
	} else {
		display.Printdisplay("You decide to go further South")

		time.Sleep(2 * time.Second)

		display.Printdisplay("After much walking you find a cluster of buildings, they are grouped in small sets, forming streets inbetween them with what seems to be a townsquare in the center of it all. A sign lets you know that this village is known as \"Mud\"")
	
		game.learnedMudLocation = true
	}

	time.Sleep(time.Second)

	display.Printdisplay("From the corner of your eye you perceive a shadowy figure spy on you, hurrying to close the door of the building it is living in")

	// TODO elaborate

	// TODO add haggling | added this to silence the compiler
	resultConversationalChallengeRaven := conversationalChallenge(player, &raven)
	// Added this to silence the compiler
	fmt.Println(resultConversationalChallengeRaven)
}

func decisionSouthDecisionForest() {
	display.Printdisplay("You penetrate the thick vegatation")

	time.Sleep(2 * time.Second)

	// TODO elaborate
}

func askXavierWhereWeAre() {
	display.Printdisplay("You inquire about this place")

	time.Sleep(time.Second)
	
	display.Printdisplay("\"Have you struck your head with a rock or something?\" You do not know how to react towards this valid hypothesis about your current situation dressed up as a quip, you feel a sharp pain in the back of your head. \"This is the small town of Mud, part of the parish of Soulstar\"")

	time.Sleep(time.Second)
}

type GameState struct {
    learnedMudLocation bool
}

var game GameState

func askXavierForDirections() {
	display.Printdisplay("You ask Xavier for directions")

	time.Sleep(time.Second)

	display.Printdisplay("\"Mud is South of here if that's what you're asking\" he pauses for a second, \"so keep traveling in the direction you were going and you should be good\"")

	game.learnedMudLocation = true

	time.Sleep(time.Second)
}

func askXavierIfYouCanDrinkFromTheFaucet(player *Character, xavier *Character) {
	display.Printdisplay("\"Huh, I would let you... but I'm worried about the cursed, you know?\"")

	time.Sleep(time.Second)

	resultConversationalChallengeXavier := conversationalChallenge(player, xavier)
	
	switch resultConversationalChallengeXavier {
		case 2: // Complete success
			display.Printdisplay("The man allows you to drink; you take a refreshing sip. You have regenerated 5 hitpoints!")
			player.CurrentHitpoints += 5
	
		case 1: // Faux pas
			display.Printdisplay("The man hesitates but, after a while, allows you to drink; you take a refreshing sip. You have regenerated 5 hitpoints!")
			player.CurrentHitpoints += 5
	
		case 0: // No healing
			display.Printdisplay("\"Sorry fella, can't do. \"Better safe than sorry\" as they say\"")
	}

	time.Sleep(time.Second)
}

func sayNothingToXavier() {
	display.Printdisplay("You say nothing")

	time.Sleep(time.Second)

	display.Printdisplay("The man coughs, not knowing where to look exactly. So, huh... what leads you to this place?")

	time.Sleep(time.Second)
}

func continueToSayNothingToXavier() {
	display.Printdisplay("You continue to say nothing")

	time.Sleep(time.Second)

	display.Printdisplay("The man in front of you decides to shift his focus away from you, continuing with what he was doing previous to your presence in that place")

	time.Sleep(time.Second)
}

func combat(player *Character, enemy *Character, distance int) {

	var chosenWeapon string
	var playerWeapon Weapon
	
	display.Printdisplay("Will you use melee or ranged weapons?")

	for {
		fmt.Scanln(&chosenWeapon)
    	chosenWeapon = strings.ToLower(strings.TrimSpace(chosenWeapon))

		if chosenWeapon == "1" || chosenWeapon == "melee" {
			
			time.Sleep(time.Second)

			display.Printdisplay("You have chosen to get close and personal")

			playerWeapon = player.MeleeWeapon

			break

		} else if chosenWeapon == "2" || chosenWeapon == "ranged" {

			time.Sleep(time.Second)

			display.Printdisplay("You have chosen to combat at a distance")

			playerWeapon = player.RangedWeapon

			break
			
		} else {
			display.Printdisplay("You hesitate, unable to choose.\n")
			display.Printdisplay("Please enter either \"1\" / \"melee\" or \"2\" / \"ranged\".\n")
		}

	}

	var enemyWeapon Weapon

	switch {
		case enemy.MeleeWeapon.weaponIsEmpty():
    		enemyWeapon = enemy.RangedWeapon
		case enemy.RangedWeapon.weaponIsEmpty():
    		enemyWeapon = enemy.MeleeWeapon
		case distance <= 5:
    		enemyWeapon = enemy.MeleeWeapon
		default:
    		enemyWeapon = enemy.RangedWeapon
	}

	// Determine turn order (initiative) based on agility
	var first *Character
	var second *Character
	var firstWeapon Weapon
	var secondWeapon Weapon

	if player.Agility >= enemy.Agility {
		first = player
		second = enemy
		firstWeapon = playerWeapon
		secondWeapon = enemyWeapon
	} else {
		first = enemy
		second = player
		firstWeapon = enemyWeapon
		secondWeapon = playerWeapon
	}

	for player.CurrentHitpoints > 0 && enemy.CurrentHitpoints > 0 {

		// First character attacks
		display.Printdisplay(fmt.Sprintf("%s attacks %s with %s %s!\n", first.Name, second.Name, first.PossesivePronoun, firstWeapon.Name,))
		var rollEvasionSecondCharacter uint8 = uint8(rand.Intn(100))
		var rollCriticalStrikeFirstCharacter uint8 = uint8(rand.Intn(100))
		var randomAdditionalDamageFirstAttack uint8 = uint8(rand.Intn(2))

		if (second.Agility >= rollEvasionSecondCharacter) {
			display.Printdisplay(fmt.Sprintf("%s has evaded %s's attack!\n", second.Name, first.Name))
		} else { 
			if first.Luck >= rollCriticalStrikeFirstCharacter {
				// Critical strike
				second.CurrentHitpoints -= int((firstWeapon.Damage + randomAdditionalDamageFirstAttack) * 3)
				display.Printdisplay(fmt.Sprintf("%s has struck %s critically!\n", first.Name, second.Name))
			} else {
				//Normal strike
				second.CurrentHitpoints -= int(firstWeapon.Damage + randomAdditionalDamageFirstAttack)
			}  
		}

		display.Printdisplay(fmt.Sprintf("%s has %d hitpoints left\n", second.Name, second.CurrentHitpoints))

		if second.CurrentHitpoints <= 0 {
			display.Printdisplay(fmt.Sprintf("%s is defeated!\n", second.Name))
			break
		}

		time.Sleep(time.Second)

		// Second character attacks
		display.Printdisplay(fmt.Sprintf("%s attacks %s with %s %s!\n", second.Name, first.Name, second.PossesivePronoun, secondWeapon.Name))

		var rollEvasionFirstCharacter uint8 = uint8(rand.Intn(100))
		var rollCriticalStrikeSecondCharacter uint8 = uint8(rand.Intn(100))
		var randomAdditionalDamageSecondAttack uint8 = uint8(rand.Intn(2))

		if (first.Agility >= rollEvasionFirstCharacter) {
			display.Printdisplay(fmt.Sprintf("%s has evaded %s's attack!\n", first.Name, second.Name))
		} else {
			if second.Luck >= rollCriticalStrikeSecondCharacter {
				// Critical strike
				first.CurrentHitpoints -= int((secondWeapon.Damage + randomAdditionalDamageSecondAttack) * 3)
				display.Printdisplay(fmt.Sprintf("%s has struck %s critically!\n", second.Name, first.Name))
			} else {
				//Normal strike
				first.CurrentHitpoints -= int(secondWeapon.Damage + randomAdditionalDamageSecondAttack)
			}
		}

		display.Printdisplay(fmt.Sprintf("%s has %d hitpoints left\n", first.Name, first.CurrentHitpoints))

		if first.CurrentHitpoints <= 0 {
			display.Printdisplay(fmt.Sprintf("%s is defeated!\n", first.Name))
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

	display.Printdisplay(fmt.Sprintf("You have engaged in a dialectical encounter against %s", enemy.Name))

	time.Sleep(time.Second)

    // Complete success
    if charismaRoll < threshold {
		display.Printdisplay("You stand on the shoulders of giants and you have managed to convince your interlocutor")
		time.Sleep(time.Second)
		player.addExperience(int(enemy.Charisma) * 10)
        return 2
    }

    // The player's luck and agility allows him to savage the situation and make a faux pas instead of utterly failing
    salvage := int(player.Luck) + int(player.Agility)
    if fauxPasRoll < salvage {
		display.Printdisplay("You have awkwardly fumbled your way into a faux pas, it could be worse, but it certainly could be better")
		player.addExperience(int(enemy.Charisma) * 5)
		// Longer wait for further awkwardness
		time.Sleep(3 * time.Second)
        return 1
    }

    // Complete failure
	display.Printdisplay("You utterly failed at trying to convince your interlocutor of anything, you may as well talk to a wall")
	time.Sleep(time.Second)
    return 0

}

func (player *Character) addExperience(ExperienceToAdd int) {
	time.Sleep(time.Second)

	// Adding experience
	player.Experience += ExperienceToAdd
	display.Printdisplay(fmt.Sprintf("You receive %d experience points\n", ExperienceToAdd))

	time.Sleep(time.Second)

	// Leveling up and the required experience for next level up
	for player.Experience >= player.ExperienceToNextLevel {
		player.Experience -= player.ExperienceToNextLevel
		player.Level++
		player.ExperienceToNextLevel = 50 * player.Level
		display.Printdisplay("You have leveled up!")
		time.Sleep(time.Second)
		display.Printdisplay(fmt.Sprintf("You are now level %d, and you need %d experience to level up again.\n", player.Level, player.ExperienceToNextLevel,))
		time.Sleep(time.Second)
	}
}