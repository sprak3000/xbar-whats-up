package configuration

import (
	"io/fs"
	"io/ioutil"
)

// Writer provides the requirements for anyone implementing writing configuration files to disk
type Writer interface {
	WriteFile(filename string, data []byte, perm fs.FileMode) error
}

// FileWriter implements the Writer interface for disk based configuration files
type FileWriter struct {
}

// WriteFile allows us to write configuration to disk
func (w FileWriter) WriteFile(filename string, data []byte, perm fs.FileMode) error {
	return ioutil.WriteFile(filename, data, perm)
}
