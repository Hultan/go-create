package template

import (
	"fmt"
	"os"
)

type Type int

const (
	Normal Type = iota
	GTK
	P5
)

type Project struct {
	Name        string
	Path        string
	Description string
	Template    Type
}

func (p *Project) Create() {
	switch p.Template {
	case Normal:
		normal := normalTemplate{p}
		normal.create()
	case GTK:
		gtk := gtkTemplate{p}
		gtk.create()
	case P5:
		p5 := p5Template{p}
		p5.create()
	}
}

func createFolder(path string) {
	err := os.Mkdir(path, 0755)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to create folder (%s) : %v\n", path, err)
	}
}
