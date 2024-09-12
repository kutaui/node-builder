package cmd

import "fmt"

func installDependencies(project Project) error {
	installCmd, saveDev := getPackageManagerCommands(project.PackageMng)

	fmt.Println("Installing base packages...")
	if err := installBasePackages(project.PackageMng, installCmd, saveDev); err != nil {
		return fmt.Errorf("error installing base packages: %v", err)
	}

	fmt.Println("Installing framework...")
	if err := installFramework(project.PackageMng, installCmd, project.Framework); err != nil {
		return fmt.Errorf("error installing framework: %v", err)
	}

	fmt.Println("Installing database and ORM...")
	if err := installDatabase(project); err != nil {
		return fmt.Errorf("error installing database and ORM: %v", err)
	}

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
		"typeorm": {
			"mysql":      "mysql2",
			"postgresql": "pg",
			"sqlite":     "sqlite3",
		},
	}

	installCmd, _ := getPackageManagerCommands(project.PackageMng)

	if project.Orm == "drizzle" || project.Orm == "sequelize" || project.Orm == "typeorm" {
		dbPackage := dbPackages[project.Orm][project.Database]
		return execCommand(project.PackageMng, installCmd, project.Orm, dbPackage)
	}
	if project.Orm == "prisma" {
		return execCommand(project.PackageMng, installCmd, "prisma")
	}
	return fmt.Errorf("unsupported ORM: %s", project.Orm)
}

func getPackageManagerCommands(packageMng string) (installCmd, saveDev string) {
	switch packageMng {
	case "npm":
		return "install", "--save-dev"
	case "yarn":
		return "add", "--dev"
	case "pnpm":
		return "add", "--save-dev"
	case "bun":
		return "add", "--dev"
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
