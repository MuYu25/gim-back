package main

import (
	"project/model"
	"project/routes"
)

func main() {
	model.InitDb()
	routes.InitRouter()
}
