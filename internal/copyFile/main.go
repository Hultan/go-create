package copyFile

import (
	"fmt"
	"io"
	"os"
	"path"
)

type CopyFilePath struct {
	BasePath     string
	RelativePath string
	FileName     string
}

type CopyFileOperation struct {
	From *CopyFilePath
	To   *CopyFilePath
}

func (c *CopyFileOperation) SetFileName(fileName string) {
	c.From.FileName = fileName
	c.To.FileName = fileName
}

func (c *CopyFileOperation) SetRelativePath(relativePath string) {
	c.From.RelativePath = relativePath
	c.To.RelativePath = relativePath
}

func (c *CopyFileOperation) SetBasePath(basePath string) {
	c.From.BasePath = basePath
	c.To.BasePath = basePath
}

func (c *CopyFileOperation) FromFilePath() string {
	return path.Join(c.From.BasePath, c.From.RelativePath, c.From.FileName)
}

func (c *CopyFileOperation) ToFilePath() string {
	return path.Join(c.To.BasePath, c.To.RelativePath, c.To.FileName)
}

func (c *CopyFileOperation) getSuccessMessage() string {
	return fmt.Sprintf("Copied '%s' to '%s'...\n", c.FromFilePath(), c.ToFilePath())
}

func (c *CopyFileOperation) getFailureMessage(err error) string {
	return fmt.Sprintf("Failed to copy '%s' to '%s' :\n\t%v\n", c.FromFilePath(), c.ToFilePath(), err)
}

func (c *CopyFileOperation) CopyFile() {
	err := c.copyFileInternal(c.FromFilePath(), c.ToFilePath())
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, c.getFailureMessage(err))
	}
	_, _ = fmt.Fprintf(os.Stderr, c.getSuccessMessage())
}

func (c *CopyFileOperation) copyFileInternal(fromPath, toPath string) error {
	from, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	to, err := os.Create(toPath)
	if err != nil {
		return err
	}

	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}
	return nil
}
