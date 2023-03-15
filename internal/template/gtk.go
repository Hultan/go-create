package template

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"syscall"

	"go-create/internal/copyFile"
)

const gtkBasePath = "/home/per/code/templates/gtk-startup"

type gtkTemplate struct {
	project *Project
}

func (t *gtkTemplate) create() {
	t.createProjectFolders()
	t.copyProjectFiles()
	t.goMod()
	t.gitInit()
}

func (t *gtkTemplate) createProjectFolders() {
	// Create main project directory
	createFolder(t.project.Path)

	// Create project folders
	createFolder(path.Join(t.project.Path, "assets"))
	createFolder(path.Join(t.project.Path, "bin"))
	createFolder(path.Join(t.project.Path, "internal"))
	createFolder(path.Join(t.project.Path, "internal", t.project.Name))
	createFolder(path.Join(t.project.Path, "internal", t.project.Name, "assets"))
	createFolder(path.Join(t.project.Path, ".run"))
}

func (t *gtkTemplate) copyProjectFiles() {
	// BASE FILES
	cfo := &copyFile.CopyFileOperation{
		From:        &copyFile.CopyFilePath{BasePath: gtkBasePath},
		To:          &copyFile.CopyFilePath{BasePath: t.project.Path},
		ProjectName: t.project.Name,
		Description: t.project.Description,
	}
	cfo.CopyFile(".gitignore")
	cfo.CopyFile("makefile")
	cfo.CopyFile("readme.md")

	// ASSETS
	cfo.SetRelativePath("assets")
	cfo.CopyFile("application.png")
	cfo.CopyFile("main.glade")

	// MAIN FILES
	cfo.From.RelativePath = ""
	cfo.To.RelativePath = ""
	cfo.CopyFile("main.go")

	// INTERNAL FILES
	cfo.From.RelativePath = "internal/gtk-startup"
	cfo.To.RelativePath = fmt.Sprintf("internal/%s", t.project.Name)
	cfo.CopyFile("mainForm.go")
	cfo.CopyFile("extraForm.go")
	cfo.CopyFile("dialog.go")
	cfo.CopyFile("aboutDialog.go")
	cfo.CopyFile("gtkBuilder.go")

	// INTERNAL FILES/assets
	cfo.From.RelativePath = "internal/gtk-startup/assets"
	cfo.To.RelativePath = fmt.Sprintf("internal/%s/assets", t.project.Name)
	cfo.CopyFile("application.png")
	cfo.CopyFile("main.glade")

	// RUN CONFIGURATION
	cfo.SetRelativePath(".run")
	cfo.From.FileName = "project-name.run.xml"
	cfo.To.FileName = fmt.Sprintf("%s.run.xml", t.project.Name)
	cfo.CopyFile("")
}

func (t *gtkTemplate) goMod() {
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

func (t *gtkTemplate) gitInit() {
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
