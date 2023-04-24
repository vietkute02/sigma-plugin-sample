package main

import (
	pdk "sigma-plugin/pkg"
)

func Init(config map[string]interface{}) {

}

func MainScript(sigma *pdk.SigmaPDK) {
	sigma.Log("Hello from MyFunction")
	// Use the custom PDK to implement your plugin functionality

}
