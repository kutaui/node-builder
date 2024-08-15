package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

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
	Run:   createProject,
}

func createProject(cmd *cobra.Command, args []string) {
	project := Project{
		ProjectName: promptInput("Project Name", validateProjectName),
		PackageMng:  promptSelect("Select your package manager", []string{"npm", "yarn", "pnpm", "bun"}),
		Framework:   promptSelect("Select your NodeJS framework", []string{"express", "fastify"}),
		Database:    promptSelect("Select your database", []string{"mysql", "postgresql", "sqlite"}),
		Orm:         promptSelect("Select your ORM", []string{"prisma", "drizzle", "typeorm", "sequelize"}),
	}

	createProjectStructure(project)
	installDependencies(project)
}

func createProjectStructure(project Project) {
	createFolders([]string{"controller", "route", "service", "util", "model", "middleware", "config"})
	createFileFromTemplate("main.ts", filepath.Join("cmd", "template", "main", project.Framework+".js"))
	if project.Orm == "drizzle" {
		createFileFromTemplate(filepath.Join("config", "db.ts"), filepath.Join("cmd", "template", "config", "db", "drizzle-"+project.Database+".js"))
	}
	copyPackageJson(project.ProjectName)
}

func installDependencies(project Project) {
	installCmd, saveDev := getPackageManagerCommands(project.PackageMng)
	checkError(installBasePackages(project.PackageMng, installCmd, saveDev))
	checkError(installFramework(project.PackageMng, installCmd, project.Framework))
	checkError(installDatabase(project))
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

func copyPackageJson(projectName string) {
	createFileFromTemplate("package.json", filepath.Join("cmd", "template", "package.json"))
	execCommand("npm", "pkg", "set", "name="+projectName)
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
	return execCommand(packageMng, installCmd, saveDev, "typescript@latest", "ts-node@latest", "@types/node@latest")
}

func installFramework(packageMng, installCmd, framework string) error {
	return execCommand(packageMng, installCmd, framework)
}

func createFileFromTemplate(destPath, srcPath string) {
	destFile, err := os.Create(destPath)
	checkError(err)
	defer destFile.Close()

	srcFile, err := os.Open(srcPath)
	checkError(err)
	defer srcFile.Close()

	_, err = io.Copy(destFile, srcFile)
	checkError(err)
	checkError(destFile.Sync())
}

func createFolders(folders []string) {
	for _, folder := range folders {
		checkError(os.Mkdir(folder, 0755))
	}
}
