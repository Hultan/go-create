package template

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"syscall"

	"go-create/internal/copyFile"
)

const p5BasePath = "/home/per/code/templates/p5"

type p5Template struct {
	project *Project
}

func (t *p5Template) create() {
	t.createProjectFolders()
	t.copyProjectFiles()
	t.goMod()
	t.gitInit()
}

func (t *p5Template) createProjectFolders() {
	// Create main project directory
	createFolder(t.project.Path)

	// Create project folders
	createFolder(path.Join(t.project.Path, "assets"))
	createFolder(path.Join(t.project.Path, "build"))
	createFolder(path.Join(t.project.Path, "cmd"))
	createFolder(path.Join(t.project.Path, "cmd", t.project.Name))
	createFolder(path.Join(t.project.Path, ".run"))
}

func (t *p5Template) copyProjectFiles() {
	// BASE FILES
	cfo := &copyFile.CopyFileOperation{
		From:        &copyFile.CopyFilePath{BasePath: p5BasePath},
		To:          &copyFile.CopyFilePath{BasePath: t.project.Path},
		ProjectName: t.project.Name,
		Description: t.project.Description,
	}
	cfo.CopyFile(".gitignore")
	cfo.CopyFile("readme.md")

	// ASSETS
	cfo.SetRelativePath("assets")
	cfo.CopyFile("application.png")

	// MAIN FILES
	cfo.From.RelativePath = "cmd/gtk-startup"
	cfo.To.RelativePath = fmt.Sprintf("cmd/%s", t.project.Name)
	cfo.CopyFile("main.go")

	// RUN CONFIGURATION
	cfo.SetRelativePath(".run")
	cfo.From.FileName = "project-name.run.xml"
	cfo.To.FileName = fmt.Sprintf("%s.run.xml", t.project.Name)
	cfo.CopyFile("")
}

func (t *p5Template) goMod() {
	fmt.Printf("Running : go mod init github.com/hultan/%s...\n", t.project.Name)

	command := fmt.Sprintf("cd %s;go mod init github.com/hultan/%s", t.project.Path, t.project.Name)
	cmd := exec.Command("bash", "-c", command)
	// Forces the new process to detach from the GitDiscover process
	// so that it does not die when GitDiscover dies
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to run : go mod init github.com/hultan/%s : %v\n", t.project.Name, err)
	}
	err = cmd.Process.Release()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to release process (goMod) : %v\n", err)
	}

	fmt.Println(string(output))
}

func (t *p5Template) gitInit() {
	fmt.Println("Running : git init...")

	command := fmt.Sprintf("cd %s;git init", t.project.Path)
	cmd := exec.Command("bash", "-c", command)
	// Forces the new process to detach from the GitDiscover process
	// so that it does not die when GitDiscover dies
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to run : git init : %v\n", err)
	}
	err = cmd.Process.Release()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to release process (gitInit): %v\n", err)
	}

	fmt.Println(string(output))
}
