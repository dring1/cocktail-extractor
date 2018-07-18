package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dring1/cocktailio/models"
)

func main() {

	inputFilePathPtr := flag.String("filepath", "", "Filepath to list of tasty cocktails! 🍸🍸")
	outputFilePathPtr := flag.String("output", "", "Output filepath for list of tasty cocktails! 🍸🍸")

	flag.Parse()
	cocktails, err := ParseFile(*inputFilePathPtr)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create(*outputFilePathPtr)
	if err != nil {
		log.Fatal(err)
	}
	if *outputFilePathPtr == "" {
		return
	}
	w := bufio.NewWriter(f)
	json.NewEncoder(w).Encode(cocktails)
}

func ParseFile(filePath string) ([]*models.Cocktail, error) {
	cocktailRecipes := []*models.Cocktail{}
	existingNames := make(map[string]interface{})

	// read file line by line
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	// for each line
	scanner := bufio.NewScanner(f)
	currentCocktail := &models.Cocktail{}
	for scanner.Scan() {
		rawText := scanner.Text()

		// if is digit start of a new cocktail
		// else apart of the recipe and push into ingredients or ingredient
		if isDigit(rawText[0]) {
			// new cocktail
			if currentCocktail.Name != "" {
				cocktailRecipes = append(cocktailRecipes, currentCocktail)
			}
			currentCocktail = &models.Cocktail{}
			name := strings.SplitAfter(string(rawText), ".")
			if len(name) < 2 {
				return nil, fmt.Errorf("%s too short", name)
			}
			currentCocktail.Name = strings.TrimSpace(name[1])
			cName := currentCocktail.Name
			if _, ok := existingNames[cName]; ok {
				log.Printf("Found duplicate: %s", cName)
			} else {
				existingNames[cName] = nil
			}
			//fmt.Printf("%s\n", currentCocktail.Name)
		} else if isInstruction(rawText[0]) {
			instruction := strings.SplitN(string(rawText), ".", 2)
			if len(instruction) < 2 {
				return nil, fmt.Errorf("instruction:\n%s\n too short", instruction)
			}
			currentCocktail.Instructions = append(currentCocktail.Instructions, strings.TrimSpace(instruction[1]))
			//fmt.Printf("\tinstruction:\t%s\n", instruction)
		} else {
			ingredient := strings.SplitN(string(rawText), ".", 2)
			if len(ingredient) < 2 {
				return nil, fmt.Errorf("ingredient:\n%s\n too short", ingredient)
			}
			currentCocktail.Ingredients = append(currentCocktail.Ingredients, &models.Ingredient{Name: strings.TrimSpace(ingredient[1])})
			//fmt.Printf("\tingredient:\t%s\n", ingredient)
		}
	}
	// if the first char is a number

	// read until a .
	// assume lines containing integers are each a new recipe
	// find the next line start with an integer, mark position as start of next cocktail
	// everything in between is either ingredients or a
	// put cocktails into bucket
	return cocktailRecipes, nil
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isInstruction(ch byte) bool {
	return ch == 'i' || ch == 'v'
}
