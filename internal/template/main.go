package template

import (
	"fmt"
	"os"
)

type Type int

const (
	GTK Type = iota
)

type Project struct {
	Name        string
	Path        string
	Description string
	Template    Type
}

func (p *Project) Create() {
	if p.Template == GTK {
		gtk := gtkTemplate{p}
		gtk.create()
	}
}

func createFolder(path string) {
	err := os.Mkdir(path, 0755)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to create folder (%s) : %v", path, err)
	}
}
