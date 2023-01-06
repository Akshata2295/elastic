package models

import (

	//es7 "github.com/elastic/go-elasticsearch/v7"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
)


func Client(c *gin.Context)  {
	// Create a new Elasticsearch client
	client, err := elastic.NewClient()
	if err != nil {
		// Handle error
	}
	fmt.Println(client)
}