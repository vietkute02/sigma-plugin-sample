package main

import (
	"fmt"
	"io/ioutil"
	"plugin"
	pdk "sigma-plugin/pkg"

	"gopkg.in/yaml.v2"
)

type ListPluginManifest struct {
	Plugins []PluginManifest `yaml:"plugins"`
}

type PluginManifest struct {
	Name    string                 `yaml:"name"`
	Type    string                 `yaml:"type"`
	Path    string                 `yaml:"path"`
	Summary string                 `yaml:"summary"`
	Config  map[string]interface{} `yaml:"config"`
}

func LoadPluginFromManifest(manifestPath string) error {
	// Read the manifest file into memory.
	listPluginData, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to read manifest file: %v", err)
	}

	// Parse the manifest data into a PluginManifest struct.
	var listPlugin ListPluginManifest
	err = yaml.Unmarshal(listPluginData, &listPlugin)
	if err != nil {
		return fmt.Errorf("failed to parse manifest file: %v", err)
	}

	for _, pluginManifest := range listPlugin.Plugins {
		// Open the plugin file and load its symbols.
		p, err := plugin.Open(pluginManifest.Path)
		if err != nil {
			return fmt.Errorf("failed to load plugin file: %v", err)
		}

		// Look up the `MyFunction` symbol in the plugin.
		initSymbol, err := p.Lookup("Init")
		if err != nil {
			return fmt.Errorf("failed to lookup MyFunction symbol: %v", err)
		}

		// Convert the symbol to a function with the expected signature.
		pluginInit, ok := initSymbol.(func(map[string]interface{}))
		if !ok {
			return fmt.Errorf("failed to convert MyFunction symbol to expected function signature")
		}

		pluginInit(pluginManifest.Config)

		// Look up the `MyFunction` symbol in the plugin.
		myFunctionSymbol, err := p.Lookup("MainScript")
		if err != nil {
			return fmt.Errorf("failed to lookup MyFunction symbol: %v", err)
		}

		// Convert the symbol to a function with the expected signature.
		myFunction, ok := myFunctionSymbol.(func(*pdk.SigmaPDK))
		if !ok {
			return fmt.Errorf("failed to convert MyFunction symbol to expected function signature")
		}

		// Create a new PDK instance.
		sigmaPDK := pdk.New()

		// Call the plugin's MyFunction with the PDK instance and config data.
		myFunction(sigmaPDK)

		// Plugin loaded successfully!
	}
	return nil

}

func main() {
	// Load the plugin from the manifest file.
	err := LoadPluginFromManifest("plugin.yaml")
	if err != nil {
		fmt.Printf("failed to load plugin: %v\n", err)
	}
}
