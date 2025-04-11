package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Print("Pokedex > ")
        scanner.Scan()
        data := scanner.Text()
        cleaned := cleanInput(data)
        if cleaned[0] == "exit" {
            commandExit()
        }
    }

}

func commandExit() error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil
}

func cleanInput(text string) []string {
	var cleaned []string
	trimmed := strings.TrimSpace(text)
	words := strings.SplitSeq(trimmed, " ")

	for word := range words {
		if word != "" {
			cleaned = append(cleaned, strings.ToLower(word))
		} else {
			fmt.Println("skipping empty word")
		}
	}
	return cleaned
}
