package main

import (
	"log"
	Routers "elastic/routes"
)



func main() {
	
	router := Routers.SetupRouter()

	if err := router.Run(":8900"); err != nil {
		log.Fatal(err)
	}
}