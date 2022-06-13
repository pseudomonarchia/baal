package main

import (
	"baal/cmd"
	_ "baal/docs"
)

// @title        Baal API
// @version      1.0
// @description  Baal API Doc

// @host      localhost:7001
// @BasePath  /api/v1

// @securityDefinitions.apikey  BearerToken
// @in                          header
// @name                        Authorization
func main() {
	cmd.Execute()
}
