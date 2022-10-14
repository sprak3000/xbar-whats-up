// Package configuration handles reading and writing the configuration file for the plugin
package configuration

import "io/ioutil"

// Reader provides the requirements for anyone implementing reading configuration files from disk
type Reader interface {
	ReadFile(filename string) ([]byte, error)
}

// FileReader implements the Reader interface for disk based configuration files
type FileReader struct {
}

// ReadFile allows us to read configuration off of disk
func (fr FileReader) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}
