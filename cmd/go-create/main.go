package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"go-create/internal/template"
)

const basePath = "/home/per/code"
const applicationName = "go-create"
const applicationVersion = "v1.2.1"

var reader *bufio.Reader
var force = false

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-f" {
		force = true
	}

	reader = bufio.NewReader(os.Stdin)

	_, _ = fmt.Fprintf(os.Stdout, "%s - %s\n", applicationName, applicationVersion)
	_, _ = fmt.Fprintf(os.Stdout, "Force mode (-f) : %v\n\n", force)

	// Get a project name
	input := getUserInput("Enter project name to create")
	projectName := strings.Replace(input, " ", "-", -1)
	projectPath := path.Join(basePath, projectName)

	// Make sure that it does not exist (unless force is true)
	if !force {
		if _, err := os.Stat(projectPath); err == nil {
			_, _ = fmt.Fprintf(os.Stdout, "A project with the path [%s] already exists!\n", projectPath)
			os.Exit(0)
		}
	}

	// Get a project description
	projectDescription := getUserInput("Enter project description")

	// Get the project type (template to use)
	projectTemplate, err := getTemplate()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Invalid project type!\n")
		os.Exit(0)
	}

	// Ask for confirmation
	message := fmt.Sprintf("Create project [%s]? (Y/n)", projectPath)
	input = getUserInput(message)
	if input != "" && !strings.HasPrefix(strings.ToUpper(input), "Y") {
		_, _ = fmt.Fprintf(os.Stdout, "Aborted by user...\n")
		os.Exit(0)
	}

	// Create project
	_, _ = fmt.Fprintf(os.Stdout, "Creating project '%s'...\n", projectName)

	project := &template.Project{
		Path:        projectPath,
		Name:        projectName,
		Description: projectDescription,
		Template:    projectTemplate,
	}
	project.Create()

	_, _ = fmt.Fprintf(os.Stdout, "Finished creating project '%s'...\n", projectName)
}

func getTemplate() (template.Type, error) {
	_, _ = fmt.Fprintf(os.Stdout, "Project types :\n")
	_, _ = fmt.Fprintf(os.Stdout, "  n) Normal\n")
	_, _ = fmt.Fprintf(os.Stdout, "  g) Gtk\n")
	_, _ = fmt.Fprintf(os.Stdout, "  p) P5\n")
	input := getUserInput("Select a project type")
	switch strings.ToLower(input[:1]) {
	case "n":
		return template.Normal, nil
	case "g":
		return template.GTK, nil
	case "p":
		return template.P5, nil
	default:
		err := errors.New(fmt.Sprintf("Invalid project type : %s\n", input))
		return template.Normal, err
	}
}

func getUserInput(msg string) string {
	_, _ = fmt.Fprintf(os.Stdout, "%s : ", msg)
	// ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return input[:len(input)-1]
}
