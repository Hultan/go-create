package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"

	"go-create/internal/copyFile"
)

const basePath = "/home/per/code"
const templateBasePath = "/home/per/code/template/gtk-startup"
const applicationVersion = "1.0.0"

var reader *bufio.Reader

func main() {
	reader = bufio.NewReader(os.Stdin)

	// Get a project name
	input := getUserInput("Enter project name to create")
	projectName := strings.Replace(input, " ", "-", -1)
	projectPath := path.Join(basePath, projectName)

	// Make sure that it does not exist
	if _, err := os.Stat(projectPath); err == nil {
		fmt.Fprintf(os.Stdout, "A project with the path [%s] already exists!\n", projectPath)
		os.Exit(0)
	}

	// Get a project description
	desc := getUserInput("Enter project description")

	// Ask for confirmation
	msg := fmt.Sprintf("Create project [%s]? (Y/n)", projectPath)
	input = getUserInput(msg)
	if input != "" && !strings.HasPrefix(strings.ToUpper(input), "Y") {
		fmt.Fprintf(os.Stdout, "Aborted by user...\n")
		os.Exit(0)
	}

	// Create project
	fmt.Fprintf(os.Stdout, "Creating project '%s'...\n", projectName)

	createProjectFolders(projectPath, projectName)
	copyProjectFiles(projectPath, projectName, desc)
	goMod(projectPath, projectName)
	gitInit(projectPath)

	fmt.Printf("Finished creating project '%s'...\n", projectName)
}

func copyProjectFiles(projectPath, projectName, desc string) {
	// BASE FILES
	cfo := &copyFile.CopyFileOperation{
		From:        &copyFile.CopyFilePath{BasePath: templateBasePath},
		To:          &copyFile.CopyFilePath{BasePath: projectPath},
		ProjectName: projectName,
		Description: desc,
	}
	cfo.SetFileName(".gitignore")
	cfo.CopyFile()
	cfo.SetFileName("readme.md")
	cfo.CopyFile()

	// ASSETS
	cfo.SetRelativePath("assets")
	cfo.SetFileName("application.png")
	cfo.CopyFile()
	cfo.SetFileName("main.glade")
	cfo.CopyFile()

	// MAIN FILES
	cfo.SetFileName("main.go")
	cfo.From.RelativePath = "cmd/gtk-startup"
	cfo.To.RelativePath = fmt.Sprintf("cmd/%s", projectName)
	cfo.CopyFile()

	// INTERNAL FILES
	cfo.From.RelativePath = "internal/gtk-startup"
	cfo.To.RelativePath = fmt.Sprintf("internal/%s", projectName)
	cfo.SetFileName("mainForm.go")
	cfo.CopyFile()
	cfo.SetFileName("extraForm.go")
	cfo.CopyFile()
	cfo.SetFileName("dialog.go")
	cfo.CopyFile()
	cfo.SetFileName("aboutDialog.go")
	cfo.CopyFile()

	// RUN CONFIGURATION
	cfo.SetRelativePath(".run")
	cfo.From.FileName = "project-name.run.xml"
	cfo.To.FileName = fmt.Sprintf("%s.run.xml", projectName)
	cfo.CopyFile()
}

func goMod(projectPath string, projectName string) {
	fmt.Printf("Running : go mod init github.com/hultan/%s...\n", projectName)

	command := fmt.Sprintf("cd %s;go mod init github.com/hultan/%s", projectPath, projectName)
	cmd := exec.Command("bash", "-c", command)
	// Forces the new process to detach from the GitDiscover process
	// so that it does not die when GitDiscover dies
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to run : go mod init github.com/hultan/%s : %v", projectName, err)
	}
	err = cmd.Process.Release()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to release process (goMod) : %v", err)
	}

	fmt.Println(string(output))
}

func gitInit(projectPath string) {
	fmt.Println("Running : git init...")

	command := fmt.Sprintf("cd %s;git init", projectPath)
	cmd := exec.Command("bash", "-c", command)
	// Forces the new process to detach from the GitDiscover process
	// so that it does not die when GitDiscover dies
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to run : git init : %v", err)
	}
	err = cmd.Process.Release()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to release process (gitInit): %v", err)
	}

	fmt.Println(string(output))
}

func createProjectFolders(projectPath, projectName string) {
	// Create main project directory
	createFolder(projectPath)

	// Create project folders
	createFolder(path.Join(projectPath, "assets"))
	createFolder(path.Join(projectPath, "build"))
	createFolder(path.Join(projectPath, "cmd"))
	createFolder(path.Join(projectPath, "cmd", projectName))
	createFolder(path.Join(projectPath, "internal"))
	createFolder(path.Join(projectPath, "internal", projectName))
	createFolder(path.Join(projectPath, ".run"))
}

func createFolder(path string) {
	err := os.Mkdir(path, 0755)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to create folder (%s) : %v", path, err)
	}
}

func getUserInput(msg string) string {
	fmt.Fprintf(os.Stdout, "%s : ", msg)
	// ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return input[:len(input)-1]
}
