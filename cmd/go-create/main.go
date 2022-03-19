package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	"go-create/internal/template"
)

const basePath = "/home/per/code"
const applicationName = "go-create"
const applicationVersion = "v1.1.0"

var reader *bufio.Reader

func main() {
	reader = bufio.NewReader(os.Stdin)

	_, _ = fmt.Fprintf(os.Stdout, "%s - %s\n\n", applicationName, applicationVersion)

	// Get a project name
	input := getUserInput("Enter project name to create")
	projectName := strings.Replace(input, " ", "-", -1)
	projectPath := path.Join(basePath, projectName)

	// Make sure that it does not exist
	if _, err := os.Stat(projectPath); err == nil {
		_, _ = fmt.Fprintf(os.Stdout, "A project with the path [%s] already exists!\n", projectPath)
		os.Exit(0)
	}

	// Get a project description
	projectDescription := getUserInput("Enter project description")

	// TODO : Choose a template
	projectTemplate := template.GTK

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

func getUserInput(msg string) string {
	_, _ = fmt.Fprintf(os.Stdout, "%s : ", msg)
	// ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return input[:len(input)-1]
}
