package main

import (
	"encoding/json"
	"flag"
	"github.com/tracyde/aquadyno/config"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

const (
	CONFIGMANAGER = "127.0.0.1:3191"
)

var conf = flag.String("conf", "config.json", "the config file to manage")

type Configuration struct {
	File string
}

func (c *Configuration) GetConfig(_ *struct{}, reply *config.Config) error {
	config, err := readConfig(c.File)
	if err != nil {
		return err
	}

	reply = config
	return nil
}

func writeConfig(config *config.Config, f string) error {
	// marshall config adding indents
	b, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// write config to disk
	err = ioutil.WriteFile(f, b, 0644)
	if err != nil {
		return err
	}

	return nil
}

func readConfig(f string) (*config.Config, error) {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	var config config.Config
	// unmarshal the read bytes
	err = json.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func checkConfigFile(f string) (*config.Config, error) {
	var c *config.Config

	// check if config file exists
	if _, err := os.Stat(f); err == nil {
		// file exists lets read it
		c, err = readConfig(f)
		if err != nil {
			return nil, err
		}
	} else {
		// file does not exist lets initialize it with default config
		c = config.NewDefaultConfig()
		if err := writeConfig(c, f); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func main() {
	// Parse our commandline flags
	flag.Parse()

	// When we first startup we are going to check the config file for correctness
	// we only want to fail if the first check is incorrect, otherwise we will just
	// log the error.
	_, err := checkConfigFile(*conf)
	if err != nil {
		log.Panicf("Error with config file %v: %v\n", *conf, err)
	}

	configuration := &Configuration{*conf}
	rpc.Register(configuration)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", CONFIGMANAGER)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	http.Serve(l, nil)
}
