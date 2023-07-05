package main

import (
	"fmt"
	docker_secrets "github.com/docker/go-plugins-helpers/secrets"
)

func main() {
	fmt.Println("Initializing OCI Secrets plugin")

	driver := OCIDriver{}
	err := driver.Init()
	if (err != nil) {
		fmt.Println("Failed to init driver")
		return
	}

	handler := docker_secrets.NewHandler(driver)
	handler.ServeUnix("plugin", 0)
}