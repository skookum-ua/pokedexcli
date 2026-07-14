package main
import (
	"fmt"
	"strings"
	"os"
	"bufio"
)
func main(){
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		ok := scanner.Scan()
		if !ok {
			break
		}
		userInput := scanner.Text()
		userCleanInput := cleanInput(userInput)
		if len(userCleanInput) != 0{
			fmt.Printf("Your command was: %s\n", userCleanInput[0] )
			
		}
	}
}
func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}