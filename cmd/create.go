package cmd

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Node.js project structured for backend development",
	Long:  `Create a Node.js project structured for backend development`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Choose the package manager you want to use")
		packageMng := selectPackageManager()

	},
}

func selectPackageManager() string {
	pkgmanagers := []string{"npm", "yarn", "pnpm", "bun"}
	index := -1
	var result string
	var err error

	for index < 0 {
		prompt := promptui.Select{
			Label: "Select package manager",
			Items: pkgmanagers,
		}

		index, result, err = prompt.Run()

		if index == -1 {
			pkgmanagers = append(pkgmanagers, result)
		}
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Input: %s\n", result)

	return result
}
