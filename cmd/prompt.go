package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

func promptInput(label string, validate promptui.ValidateFunc) string {
	prompt := promptui.Prompt{Label: label, Validate: validate}
	result, err := prompt.Run()
	checkError(err)
	return result
}

func promptSelect(label string, items []string) string {
	prompt := promptui.Select{Label: label, Items: items}
	_, result, err := prompt.Run()
	checkError(err)
	return result
}

func promptConfirm(prompt string) bool {
	spinnerInstance.Stop()
	defer spinnerInstance.Start()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s (y/n): ", prompt)
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return false
		}
		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
		fmt.Println("Please answer with 'y' or 'n'")
	}
}
