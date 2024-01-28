package cmd

import (
	"io"
	"log"
	"os"
	"sync"

	"github.com/redskal/steve/pkg/azrecon"
	"gopkg.in/yaml.v3"
)

var (
	sm     sync.Map
	logger = log.New(os.Stderr, "Steve: ", 0)
)

type empty struct{}

// hasStdin checks for piped input from character devices
// or FIFO. Original:
// https://github.com/projectdiscovery/fileutil/blob/380e33ef95825c6b781f289d8cd9c0d48d6c67f5/file.go#L141
func hasStdin() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	mode := stat.Mode()
	return (mode&os.ModeCharDevice) == 0 || (mode&os.ModeNamedPipe) != 0
}

// unique checks if the given string was unique for the given map.
// Modified version of https://github.com/hakluke/hakrawler/blob/master/hakrawler.go#L285
func unique(s string) bool {
	_, present := sm.Load(s)
	if present {
		return false
	}
	sm.Store(s, true)
	return true
}

// readConfigFromFile reads the specified YAML file and unmarshals
// into global config 'cfg'
func readConfigFromFile(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	fileBytes, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(fileBytes, &cfg)
	if err != nil {
		return err
	}
	return nil
}

// writeDefaultConfig creates a configuration file using
// the default resolvers from the azrecon package.
func writeDefaultConfig(fileName string) error {
	c := config{}
	c.Resolvers = append(c.Resolvers, azrecon.Resolvers...)
	fc, err := yaml.Marshal(&c)
	if err != nil {
		return err
	}
	if err := os.WriteFile(fileName, fc, 0644); err != nil {
		return err
	}
	return nil
}
