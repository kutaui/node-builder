package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type Project struct {
	ProjectName string
	PackageMng  string
	Framework   string
	Database    string
	Orm         string
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Node.js project structured for backend development",
	Long:  `Create a Node.js project structured for backend development and customized to your needs`,
	Run: func(cmd *cobra.Command, args []string) {
		project := Project{
			ProjectName: inputProjectName(),
			PackageMng:  selectPackageManager(),
			Framework:   selectBackendFramework(),
			Database:    selectDatabase(),
			Orm:         selectORM(),
		}

		copyPackageJson(project.ProjectName)

		installDependencies(project)
	},
}

func execCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

func installBasePackages(project Project, installCmd, saveDev string) error {
	return execCommand(project.PackageMng, installCmd, saveDev, "typescript@latest", "ts-node@latest", "@types/node@latest")
}

func installFramework(project Project, installCmd string) error {
	return execCommand(project.PackageMng, installCmd, project.Framework)
}

func installDependencies(project Project) error {
	installCmd := "add"
	saveDev := "--dev"
	if project.PackageMng == "npm" {
		installCmd = "install"
		saveDev = "--save-dev"
	} else if project.PackageMng == "pnpm" {
		saveDev = "--save-dev"
	}

	if err := installBasePackages(project, installCmd, saveDev); err != nil {
		return fmt.Errorf("error installing base packages: %w", err)
	}

	if err := installFramework(project, installCmd); err != nil {
		return fmt.Errorf("error installing framework: %w", err)
	}

	if err := installDatabase(project, installCmd, saveDev); err != nil {
		return fmt.Errorf("error installing database packages: %w", err)
	}

	return nil
}

func installDatabase(project Project, installCmd, saveDev string) error {
	var packages []string

	switch project.Orm {
	case "drizzle":
		packages = append(packages, "drizzle-orm")
		switch project.Database {
		case "postgresql":
			packages = append(packages, "pg")
		case "mysql":
			packages = append(packages, "mysql2")
		case "sqlite":
			packages = append(packages, "better-sqlite3")
		default:
			return fmt.Errorf("unsupported database for drizzle: %s", project.Database)
		}
		if err := execCommand(project.PackageMng, append([]string{installCmd}, packages...)...); err != nil {
			return fmt.Errorf("error installing drizzle and database driver: %w", err)
		}
		return execCommand(project.PackageMng, installCmd, saveDev, "drizzle-kit")

	case "prisma":
		return execCommand(project.PackageMng, installCmd, "prisma")

	case "sequelize":
		packages = append(packages, "sequelize")
		switch project.Database {
		case "postgresql":
			packages = append(packages, "pg", "pg-hstore")
		case "mysql":
			packages = append(packages, "mysql2")
		case "sqlite":
			packages = append(packages, "sqlite3")
		default:
			return fmt.Errorf("unsupported database for sequelize: %s", project.Database)
		}
		return execCommand(project.PackageMng, append([]string{installCmd}, packages...)...)

	default:
		return fmt.Errorf("unsupported ORM: %s", project.Orm)
	}
}

func copyPackageJson(projectName string) {
	packageJsonFile, err := os.Create(filepath.Join(".", "package.json"))
	if err != nil {
		cobra.CheckErr(err)
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer packageJsonFile.Close()

	srcFile, err := os.Open(filepath.Join("cmd", "template", "package.json"))
	if err != nil {
		cobra.CheckErr(err)
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer srcFile.Close()

	_, err = io.Copy(packageJsonFile, srcFile)
	if err != nil {
		cobra.CheckErr(err)
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	err = packageJsonFile.Sync()
	if err != nil {
		cobra.CheckErr(err)
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	cmd := exec.Command("npm", "pkg", "set", "name="+projectName)
	_, err = cmd.Output()
	if err != nil {
		cobra.CheckErr(err)
		fmt.Println("could not run command: ", err)
		os.Exit(1)
	}
}

func inputProjectName() string {
	validate := func(input string) error {
		length := len(input)
		if length < 3 {
			return errors.New("Project name should be more than 3 characters")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Project Name",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	return result
}

func selectPackageManager() string {
	pkgmanagers := []string{"npm", "yarn", "pnpm", "bun"}
	index := -1
	var result string
	var err error

	for index < 0 {
		prompt := promptui.Select{
			Label: "Select your package manager",
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

	return result
}

func selectBackendFramework() string {
	frameworks := []string{"express", "fastify", "koa"}
	index := -1
	var result string
	var err error

	for index < 0 {
		prompt := promptui.Select{
			Label: "Select your NodeJS framework",
			Items: frameworks,
		}

		index, result, err = prompt.Run()

		if index == -1 {
			frameworks = append(frameworks, result)
		}
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	return result
}

func selectDatabase() string {
	databases := []string{"mysql", "postgresql", "sqlite"}
	index := -1
	var result string
	var err error

	for index < 0 {
		prompt := promptui.Select{
			Label: "Select your database",
			Items: databases,
		}

		index, result, err = prompt.Run()

		if index == -1 {
			databases = append(databases, result)
		}
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	return result
}

func selectORM() string {
	orms := []string{"prisma", "drizzle", "typeorm", "sequelize"}
	index := -1
	var result string
	var err error

	for index < 0 {
		prompt := promptui.Select{
			Label: "Select your ORM",
			Items: orms,
		}

		index, result, err = prompt.Run()

		if index == -1 {
			orms = append(orms, result)
		}
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	return result
}
