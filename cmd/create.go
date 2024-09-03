package cmd

import (
	"bufio"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

//go:embed template
var templateFS embed.FS

type Project struct {
	ProjectName string
	PackageMng  string
	Framework   string
	Database    string
	Orm         string
	UseEslint   bool
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Node.js project structured for backend development",
	Long:  `Create a Node.js project structured for backend development and customized to your needs`,
	Run:   createProject,
}

var spinnerInstance *spinner.Spinner

func init() {
	spinnerInstance = spinner.New(spinner.CharSets[9], 100*time.Millisecond)
}

func createProject(cmd *cobra.Command, args []string) {
	spinnerInstance.Start()

	project := Project{
		ProjectName: promptInput("Project Name (enter '.' for current folder)", validateProjectName),
		PackageMng:  promptSelect("Select your package manager", []string{"npm", "yarn", "pnpm", "bun"}),
		Framework:   promptSelect("Select your NodeJS framework", []string{"express", "fastify"}),
		Database:    promptSelect("Select your database", []string{"mysql", "postgresql", "sqlite"}),
		Orm:         promptSelect("Select your ORM", []string{"prisma", "drizzle", "typeorm", "sequelize"}),
		UseEslint:   promptConfirm("Do you want to use ESLint?"),
	}

	if project.ProjectName == "." {
		currentDir, err := os.Getwd()
		if err != nil {
			spinnerInstance.Stop()
			fmt.Printf("Error getting current directory: %v\n", err)
			return
		}
		project.ProjectName = filepath.Base(currentDir)
	} else {
		err := os.Mkdir(project.ProjectName, 0755)
		if err != nil {
			spinnerInstance.Stop()
			fmt.Printf("Error creating project directory: %v\n", err)
			return
		}
		err = os.Chdir(project.ProjectName)
		if err != nil {
			spinnerInstance.Stop()
			fmt.Printf("Error changing to project directory: %v\n", err)
			return
		}
	}

	spinnerInstance.Suffix = " Creating project structure..."
	fmt.Println("Starting to create project structure...")
	err := createProjectStructure(project)
	if err != nil {
		spinnerInstance.Stop()
		fmt.Printf("Error creating project structure: %v\n", err)
		return
	}
	fmt.Println("Project structure created successfully.")

	spinnerInstance.Suffix = " Installing dependencies..."
	fmt.Println("Starting to install dependencies...")
	err = installDependencies(project)
	if err != nil {
		spinnerInstance.Stop()
		fmt.Printf("Error installing dependencies: %v\n", err)
		return
	}
	fmt.Println("Dependencies installed successfully.")

	if project.UseEslint {
		spinnerInstance.Suffix = " Setting up ESLint..."
		fmt.Println("Setting up ESLint...")
		err = setupEslint(project)
		if err != nil {
			spinnerInstance.Stop()
			fmt.Printf("Error setting up ESLint: %v\n", err)
			return
		}
		fmt.Println("ESLint setup completed successfully.")
	}

	spinnerInstance.Stop()
	fmt.Println("Project created successfully!")
}

func createProjectStructure(project Project) error {
	fmt.Println("Starting createProjectStructure...")

	folders := []string{"controller", "route", "service", "util", "model", "middleware", "config"}
	for _, folder := range folders {
		fmt.Printf("Attempting to create folder: %s\n", folder)
		if err := createFolderIfNotExists(folder); err != nil {
			return fmt.Errorf("error creating folder %s: %v", folder, err)
		}
	}
	fmt.Println("Finished creating folders.")

	fmt.Printf("Reading template file for %s framework...\n", project.Framework)
	templatePath := fmt.Sprintf("template/main/%s.js", project.Framework)
	templateContent, err := fs.ReadFile(templateFS, templatePath)
	if err != nil {
		fmt.Printf("Warning: Template file for %s framework not found. Skipping main file creation.\n", project.Framework)
	} else {
		fmt.Println("Creating main.ts file...")
		if err := createFileWithConfirmation("main.ts", templateContent); err != nil {
			return err
		}
	}

	if project.Orm == "drizzle" {
		fmt.Println("Setting up Drizzle ORM configuration...")
		if err := createFileFromTemplate(filepath.Join("config", "db.ts"), fmt.Sprintf("template/config/db/drizzle-%s.js", project.Database)); err != nil {
			return err
		}
	}

	fmt.Println("Creating .env file...")
	if err := createEmptyFile(".env"); err != nil {
		return err
	}

	fmt.Println("Creating .gitignore file...")
	if err := createFileFromTemplate(".gitignore", "template/gitignore"); err != nil {
		return err
	}

	fmt.Println("Copying package.json...")
	if err := copyPackageJson(project.ProjectName); err != nil {
		return err
	}

	fmt.Println("Project structure creation completed.")
	return nil
}

func createFolderIfNotExists(folder string) error {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err := os.Mkdir(folder, 0755)
		if err != nil {
			return err
		}
		fmt.Printf("Created folder '%s'\n", folder)
	} else if err != nil {
		return err
	} else {
		fmt.Printf("Folder '%s' already exists. Skipping creation.\n", folder)
	}
	return nil
}

func createFileWithConfirmation(filename string, content []byte) error {
	fmt.Printf("Attempting to create file: %s\n", filename)
	if _, err := os.Stat(filename); err == nil {
		fmt.Printf("File %s already exists.\n", filename)
		spinnerInstance.Stop()
		if !confirmOverwrite(filename) {
			spinnerInstance.Start()
			fmt.Printf("Skipping creation of %s\n", filename)
			return nil
		}
		spinnerInstance.Start()
		fmt.Printf("Overwriting %s\n", filename)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("error checking file %s: %v", filename, err)
	}

	fmt.Printf("Writing content to %s\n", filename)
	err := os.WriteFile(filename, content, 0644)
	if err != nil {
		return fmt.Errorf("error creating %s: %v", filename, err)
	}
	fmt.Printf("Successfully created %s\n", filename)
	return nil
}

func confirmOverwrite(filename string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("File %s already exists. Do you want to overwrite it? (y/n): ", filename)
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

func installDependencies(project Project) error {
	installCmd, saveDev := getPackageManagerCommands(project.PackageMng)
	checkError(installBasePackages(project.PackageMng, installCmd, saveDev))

	checkError(installFramework(project.PackageMng, installCmd, project.Framework))

	checkError(installDatabase(project))

	return nil
}

func installDatabase(project Project) error {
	dbPackages := map[string]map[string]string{
		"drizzle": {
			"mysql":      "mysql2",
			"postgresql": "pg",
			"sqlite":     "better-sqlite3",
		},
		"sequelize": {
			"mysql":      "mysql2",
			"postgresql": "pg pg-hstore",
			"sqlite":     "sqlite3",
		},
	}

	if project.Orm == "drizzle" || project.Orm == "sequelize" {
		dbPackage := dbPackages[project.Orm][project.Database]
		return execCommand(project.PackageMng, "install", project.Orm, dbPackage)
	}
	if project.Orm == "prisma" {
		return execCommand(project.PackageMng, "install", "prisma")
	}
	return fmt.Errorf("unsupported ORM: %s", project.Orm)
}

func copyPackageJson(projectName string) error {
	if err := createFileFromTemplate("package.json", "template/package.json"); err != nil {
		return err
	}
	return execCommand("npm", "pkg", "set", "name="+projectName)
}

func getPackageManagerCommands(packageMng string) (installCmd, saveDev string) {
	switch packageMng {
	case "npm":
		return "install", "--save-dev"
	case "yarn":
		return "add", "--dev"
	case "pnpm":
		return "add", "--save-dev"
	default:
		return "install", "--save-dev"
	}
}

func installBasePackages(packageMng, installCmd, saveDev string) error {
	return execCommand(packageMng, installCmd, saveDev, "typescript@latest", "ts-node@latest", "@types/node@latest", "prettier")
}

func installFramework(packageMng, installCmd, framework string) error {
	return execCommand(packageMng, installCmd, framework)
}

func createFileFromTemplate(destPath, srcPath string) error {
	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	srcContent, err := templateFS.ReadFile(srcPath)
	if err != nil {
		return err
	}

	_, err = destFile.Write(srcContent)
	return err
}

func createEmptyFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating %s: %v", filename, err)
	}
	defer file.Close()
	fmt.Printf("Created empty file: %s\n", filename)
	return nil
}

func setupEslint(project Project) error {
	installCmd, saveDev := getPackageManagerCommands(project.PackageMng)
	checkError(execCommand(project.PackageMng, installCmd, saveDev, "eslint", "@typescript-eslint/parser", "@typescript-eslint/eslint-plugin"))

	eslintConfig := `{
		"parser": "@typescript-eslint/parser",
		"plugins": ["@typescript-eslint"],
		"extends": [
			"eslint:recommended",
			"plugin:@typescript-eslint/recommended"
		]
	}`

	return createFileWithConfirmation(".eslintrc.json", []byte(eslintConfig))
}

func promptConfirm(prompt string) bool {
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
