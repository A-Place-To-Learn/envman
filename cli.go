package main

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

const (
	defaultEnvStoreName string = ".envstore.yml"
)

var (
	stdinValue              string
	currentEnvStoreFilePath string
)

func isPipedData() bool {
	if stat, err := os.Stdin.Stat(); err != nil {
		return false
	} else if (stat.Mode() & os.ModeCharDevice) == 0 {
		return true
	}
	return false
}

func envStorePathInCurrentDir() (string, error) {
	return filepath.Abs(path.Join("./", defaultEnvStoreName))
}

func before(c *cli.Context) error {
	if level, err := log.ParseLevel(c.String("log-level")); err != nil {
		log.Fatal(err.Error())
	} else {
		log.SetLevel(level)
	}

	// Befor parsing cli, and running command
	// we need to decide wich path will be used by envman
	currentEnvStoreFilePath = c.String(PATH_KEY)
	if currentEnvStoreFilePath == "" {
		if path, err := envStorePathInCurrentDir(); err != nil {
			log.Fatal("Failed to set envman work path in current dir:", err)
		} else {
			currentEnvStoreFilePath = path
		}
	}
	return nil
}

// Run the Envman CLI.
func run() {
	log.SetLevel(log.DebugLevel)

	// Read piped data
	if isPipedData() {
		if bytes, err := ioutil.ReadAll(os.Stdin); err != nil {
			log.Error("Failed to read stdin:", err)
		} else if len(bytes) > 0 {
			stdinValue = string(bytes)
		}
	}

	// Parse cl
	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = "Environment variable manager"
	app.Version = "0.0.5"

	app.Author = ""
	app.Email = ""

	app.Before = before

	app.Flags = flags
	app.Commands = commands

	if err := app.Run(os.Args); err != nil {
		log.Fatal("Envman finished:", err)
	}
}
