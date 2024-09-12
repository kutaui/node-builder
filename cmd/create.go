package cmd

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
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

	fmt.Printf("Debug: Project details collected: %+v\n", project)

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
		fmt.Println("You may need to run the installation manually.")
	} else {
		fmt.Println("Dependencies installed successfully.")
	}

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

	if project.Orm == "prisma" {
		fmt.Println("\nSince you chose Prisma as your ORM, here's some helpful documentation:")
		fmt.Println("- Prisma Quickstart: https://www.prisma.io/docs/getting-started/quickstart")
	} else if project.Orm == "typeorm" {
		fmt.Println("\nSince you chose TypeORM as your ORM, here's some helpful documentation:")
		fmt.Println("- TypeORM Getting Started: https://typeorm.io/#quick-start")
		fmt.Println("- TypeORM Documentation: https://typeorm.io/")
	}
}
