package cmd

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func createProjectStructure(project Project) error {
	fmt.Println("Starting createProjectStructure...")

	if err := createFolderIfNotExists("src"); err != nil {
		return fmt.Errorf("error creating src folder: %v", err)
	}

	srcFolders := []string{"controllers", "routes", "services", "utils", "models", "middlewares", "config"}
	for _, folder := range srcFolders {
		fmt.Printf("Attempting to create folder: src/%s\n", folder)
		if err := createFolderIfNotExists(filepath.Join("src", folder)); err != nil {
			return fmt.Errorf("error creating folder src/%s: %v", folder, err)
		}
	}
	fmt.Println("Finished creating folders.")

	fmt.Printf("Reading template file for %s framework...\n", project.Framework)
	templatePath := fmt.Sprintf("template/main/%s.js", project.Framework)
	templateContent, err := fs.ReadFile(templateFS, templatePath)
	if err != nil {
		fmt.Printf("Warning: Template file for %s framework not found. Skipping main file creation.\n", project.Framework)
	} else {
		fmt.Println("Creating src/main.ts file...")
		if err := createFileWithConfirmation("src/main.ts", templateContent); err != nil {
			return err
		}
	}

	if project.Orm == "drizzle" {
		fmt.Println("Setting up Drizzle ORM configuration...")
		if err := createFileFromTemplate(filepath.Join("src", "config", "db.ts"), fmt.Sprintf("template/config/db/drizzle-%s.js", project.Database)); err != nil {
			return err
		}
	}

	if project.Orm == "sequelize" {
		fmt.Println("Setting up Sequelize ORM configuration...")
		if err := createFileFromTemplate(filepath.Join("src", "config", "db.ts"), "template/config/db/sequelize.js"); err != nil {
			return err
		}
	}

	if project.Orm == "typeorm" {
		fmt.Println("Setting up TypeORM configuration...")
		if err := createTypeORMConfig(project); err != nil {
			return err
		}

		// Create entity, migration, and subscriber folders inside src
		for _, folder := range []string{"entities", "migrations", "subscribers"} {
			if err := createFolderIfNotExists(filepath.Join("src", folder)); err != nil {
				return err
			}
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

	fmt.Println("Creating README.md file...")
	if err := createReadme(project); err != nil {
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

func createTypeORMConfig(project Project) error {
	configContent := fmt.Sprintf(`import { DataSource } from "typeorm"

export const AppDataSource = new DataSource({
    type: "%s",
    database: "your_database_name",
    entities: ["src/entities/**/*.ts"],
    migrations: ["src/migrations/**/*.ts"],
    subscribers: ["src/subscribers/**/*.ts"],
    synchronize: true,
})`, project.Database)

	return createFileWithConfirmation("src/data-source.ts", []byte(configContent))
}

func copyPackageJson(projectName string) error {
	if err := createFileFromTemplate("package.json", "template/package.json"); err != nil {
		return err
	}
	return execCommand("npm", "pkg", "set", "name="+projectName)
}
