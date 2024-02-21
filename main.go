/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"rocket/cmd"
	"rocket/service"
)

func main() {
	//init K8s clientset
	service.K8s.Init()
	cmd.Execute()
}
