package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/manifoldco/promptui"
)

func checkError(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func execCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	_, err := cmd.Output()
	return err
}

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

func validateProjectName(input string) error {
	if input == "." {
		return nil
	}

	if len(input) < 3 {
		return errors.New("Project name should be more than 3 characters")
	}
	return nil
}
