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

	if game.learnedMudLocation {
		fmt.Println("Do you...?: \n 1. Continue South for Mud \n 2. Continue into a forest that borders the fields")
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
		fmt.Println("Do you...?: \n 1. Continue South \n 2. Continue into a forest that borders the fields")
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

	// Used by Raven
	ironDagger, err := NewWeapon(
		"Iron Dagger",
		strength + 2,
		Melee,
		Bladed,
	)
	if err != nil {
		fmt.Println("Could not create Raven's iron dagger:", err)
		return
	}

	raven := Character{
		Name:             "Raven",
		Strength:         3,
		Dexterity:        3,
		Tenacity:         3,
		Agility:          3,
		Luck:             3,
		Charisma:         3,
		MeleeWeapon:      ironDagger,
		PossesivePronoun: "her",
	}

	if game.learnedMudLocation {
		fmt.Println("You decide to go further South towards Mud")

		time.Sleep(2 * time.Second)

		fmt.Println("After much walking you find a cluster of buildings, they are grouped in small sets, forming streets inbetween them with what seems to be a townsquare in the center of it all. You assume this is Mud and a sign soon confirms your suspicions")
	} else {
		fmt.Println("You decide to go further South")

		time.Sleep(2 * time.Second)

		fmt.Println("After much walking you find a cluster of buildings, they are grouped in small sets, forming streets inbetween them with what seems to be a townsquare in the center of it all. A sign lets you know that this village is known as \"Mud\"")
	
		game.learnedMudLocation = true
	}

	time.Sleep(time.Second)

	fmt.Println("From the corner of your eye you perceive a shadowy figure spy on you, hurrying to close the door of the building it is living in")

	// TODO elaborate

	// TODO add haggling | added this to silence the compiler
	resultConversationalChallengeRaven := conversationalChallenge(player, &raven)
	// Added this to silence the compiler
	fmt.Println(resultConversationalChallengeRaven)
}

func decisionSouthDecisionForest() {
	fmt.Println("You penetrate the thick vegatation")

	time.Sleep(2 * time.Second)

	// TODO elaborate
}

func askXavierWhereWeAre() {
	fmt.Println("You inquire about this place")

	time.Sleep(time.Second)
	
	fmt.Println("\"Have you struck your head with a rock or something?\" You do not know how to react towards this valid hypothesis about your current situation dressed up as a quip, you feel a sharp pain in the back of your head. \"This is the small town of Mud, part of the parish of Soulstar\"")

	time.Sleep(time.Second)
}

type GameState struct {
    learnedMudLocation bool
}

var game GameState

func askXavierForDirections() {
	fmt.Println("You ask Xavier for directions")

	time.Sleep(time.Second)

	fmt.Println("\"Mud is South of here if that's what you're asking\" he pauses for a second, \"so keep traveling in the direction you were going and you should be good\"")

	game.learnedMudLocation = true

	time.Sleep(time.Second)
}

func askXavierIfYouCanDrinkFromTheFaucet(player *Character, xavier *Character) {
	fmt.Println("\"Huh, I would let you... but I'm worried about the cursed, you know?\"")

	time.Sleep(time.Second)

	resultConversationalChallengeXavier := conversationalChallenge(player, xavier)
	
	switch resultConversationalChallengeXavier {
		case 2: // Complete success
			fmt.Println("The man allows you to drink; you take a refreshing sip. You have regenerated 5 hitpoints!")
			player.CurrentHitpoints += 5
	
		case 1: // Faux pas
			fmt.Println("The man hesitates but, after a while, allows you to drink; you take a refreshing sip. You have regenerated 5 hitpoints!")
			player.CurrentHitpoints += 5
	
		case 0: // No healing
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

// TODO implement weapons
func combat(player *Character, enemy *Character, distance int) {

	var chosenWeapon string
	var playerWeapon Weapon
	
	fmt.Println("Will you use melee or ranged weapons?")

	for {
		fmt.Scanln(&chosenWeapon)
    	chosenWeapon = strings.ToLower(strings.TrimSpace(chosenWeapon))

		if chosenWeapon == "1" || chosenWeapon == "melee" {
			
			time.Sleep(time.Second)

			fmt.Println("You have chosen to get close and personal")

			playerWeapon = player.MeleeWeapon

			break

		} else if chosenWeapon == "2" || chosenWeapon == "ranged" {

			time.Sleep(time.Second)

			fmt.Println("You have chosen to combat at a distance")

			playerWeapon = player.RangedWeapon

			break
			
		} else {
			fmt.Println("You hesitate, unable to choose.\n")
			fmt.Println("Please enter either \"1\" / \"melee\" or \"2\" / \"ranged\".\n")
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
		// TODO gender
		fmt.Printf("%s attacks %s with his %s!\n", first.Name, second.Name, firstWeapon.Name)

		var rollEvasionSecondCharacter uint8 = uint8(rand.Intn(100))
		var rollCriticalStrikeFirstCharacter uint8 = uint8(rand.Intn(100))
		var randomAdditionalDamageFirstAttack uint8 = uint8(rand.Intn(2))

		if (second.Agility >= rollEvasionSecondCharacter) {
			fmt.Printf("%s has evaded %s's attack!\n", second.Name, first.Name)
		} else { 
			if first.Luck >= rollCriticalStrikeFirstCharacter {
				// Critical strike
				second.CurrentHitpoints -= int((firstWeapon.Damage + randomAdditionalDamageFirstAttack) * 3)
				fmt.Printf("%s has struck %s critically!\n", first.Name, second.Name)
			} else {
				//Normal strike
				second.CurrentHitpoints -= int(firstWeapon.Damage + randomAdditionalDamageFirstAttack)
			}  
		}

		fmt.Printf("%s has %d hitpoints left\n", second.Name, second.CurrentHitpoints)

		if second.CurrentHitpoints <= 0 {
			fmt.Printf("%s is defeated!\n", second.Name)
			break
		}

		time.Sleep(time.Second)

		// Second character attacks
		// TODO gender
		fmt.Printf("%s attacks %s with his %s!\n", second.Name, first.Name, secondWeapon.Name)

		var rollEvasionFirstCharacter uint8 = uint8(rand.Intn(100))
		var rollCriticalStrikeSecondCharacter uint8 = uint8(rand.Intn(100))
		var randomAdditionalDamageSecondAttack uint8 = uint8(rand.Intn(2))

		if (first.Agility >= rollEvasionFirstCharacter) {
			fmt.Printf("%s has evaded %s's attack!\n", first.Name, second.Name)
		} else {
			if second.Luck >= rollCriticalStrikeSecondCharacter {
				// Critical strike
				first.CurrentHitpoints -= int((secondWeapon.Damage + randomAdditionalDamageSecondAttack) * 3)
				fmt.Printf("%s has struck %s critically!\n", second.Name, first.Name)
			} else {
				//Normal strike
				first.CurrentHitpoints -= int(secondWeapon.Damage + randomAdditionalDamageSecondAttack)
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