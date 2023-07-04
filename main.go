package main

import (
	initializers "github.com/seetharamugn/wachat/initializers"
	"github.com/seetharamugn/wachat/routers"
)

func init() {
	initializers.LoadEnvVariables()
}

func main() {

	// Set up the Gin routers
	route := routers.SetupRouter()
	route.Run()
}
