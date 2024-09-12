package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func checkError(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func execCommand(name string, args ...string) error {
	fmt.Printf("Executing command: %s %s\n", name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
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

func createReadme(project Project) error {
	readmeContent := `# ` + project.ProjectName + `

This project was generated using [Node Builder](https://github.com/kutaui/node-builder).

## 🚀 Quick Start

### Prerequisites

- Node.js (version 14 or higher)
- ` + project.PackageMng + ` (package manager)

### Installation

1. Clone the repository:
   ` + "```bash" + `
   git clone https://github.com/your-username/` + project.ProjectName + `.git
   cd ` + project.ProjectName + `
   ` + "```" + `

2. Install dependencies:
   ` + "```bash" + `
   ` + project.PackageMng + ` install
   ` + "```" + `

3. Set up your environment variables:
   ` + "```bash" + `
   cp .env.example .env
   ` + "```" + `
   Then, edit the ` + "`.env`" + ` file with your configuration.

### Running the Application

To start the development server:

` + "```bash" + `
` + project.PackageMng + ` run dev
` + "```" + `

The application will be available at ` + "`http://localhost:3000`" + ` (or the port specified in your environment variables).

## 🏗️ Project Structure

` + "```" + `
` + project.ProjectName + `/
├── src/
│   ├── controllers/    # Request handlers
│   ├── routes/         # Application routes
│   ├── services/       # Business logic
│   ├── utils/          # Utility functions
│   ├── models/         # Data models
│   ├── middlewares/    # Express middlewares
│   ├── config/         # Configuration files
│   └── main.ts         # Application entry point
├── .env                # Environment variables
├── .gitignore
├── README.md
└── package.json
` + "```" + `

## 🛠️ Built With

- [Node.js](https://nodejs.org/)
- [` + strings.Title(project.Framework) + `](` + getFrameworkURL(project.Framework) + `) - Web framework
- [` + strings.Title(project.Database) + `](` + getDatabaseURL(project.Database) + `) - Database
- [` + strings.Title(project.Orm) + `](` + getORMURL(project.Orm) + `) - ORM

## 📚 Additional Resources

- [Express.js Documentation](https://expressjs.com/)
- [Node.js Documentation](https://nodejs.org/en/docs/)

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

🌟 Happy coding! If you have any questions or run into issues, please open an issue on the [Node Builder repository](https://github.com/kutaui/node-builder/issues).
`

	return createFileWithConfirmation("README.md", []byte(readmeContent))
}

func getFrameworkURL(framework string) string {
	switch framework {
	case "express":
		return "https://expressjs.com/"
	case "fastify":
		return "https://www.fastify.io/"
	default:
		return ""
	}
}

func getDatabaseURL(database string) string {
	switch database {
	case "mysql":
		return "https://www.mysql.com/"
	case "postgresql":
		return "https://www.postgresql.org/"
	case "sqlite":
		return "https://www.sqlite.org/"
	default:
		return ""
	}
}

func getORMURL(orm string) string {
	switch orm {
	case "prisma":
		return "https://www.prisma.io/"
	case "drizzle":
		return "https://github.com/drizzle-team/drizzle-orm"
	case "typeorm":
		return "https://typeorm.io/"
	case "sequelize":
		return "https://sequelize.org/"
	default:
		return ""
	}
}
