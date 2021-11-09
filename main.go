package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/sensu/sensu-go/types"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	Example string
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-test-check",
			Short:    "test sensu check",
			Keyspace: "sensu.io/plugins/sensu-test-check/config",
		},
	}

	options = []*sensu.PluginConfigOption{
		&sensu.PluginConfigOption{
			Path:      "example",
			Env:       "CHECK_EXAMPLE",
			Argument:  "example",
			Shorthand: "e",
			Default:   "",
			Usage:     "An example string configuration option",
			Value:     &plugin.Example,
		},
	}
)

var urls = map[string]string{
	"MainPage":  "http://geocitizen.link:8080/citizen",
	"LoginPage": "http://geocitizen.link:8080/citizen/#/auth",
}

func main() {
	useStdin := false
	fi, err := os.Stdin.Stat()
	if err != nil {
		fmt.Printf("Error check stdin: %v\n", err)
		panic(err)
	}
	//Check the Mode bitmask for Named Pipe to indicate stdin is connected
	if fi.Mode()&os.ModeNamedPipe != 0 {
		log.Println("using stdin")
		useStdin = true
	}

	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, useStdin)
	check.Execute()
}

func checkArgs(event *types.Event) (int, error) {
	if len(plugin.Example) == 0 {
		return sensu.CheckStateWarning, fmt.Errorf("--example or CHECK_EXAMPLE environment variable is required")
	}
	return sensu.CheckStateOK, nil
}

func executeCheck(event *types.Event) (int, error) {
	for service, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("%s: check ERROR: %s\n", service, err)
			return sensu.CheckStateCritical, nil
		}
		resp.Body.Close()
		if resp.StatusCode != 200 {
			log.Printf("%s: status check ERROR: %d != 200\n", service, resp.StatusCode)
			return sensu.CheckStateCritical, nil
		}
		log.Printf("%s: status OK", service)
	}
	return sensu.CheckStateOK, nil
}
