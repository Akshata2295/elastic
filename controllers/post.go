package controllers

import (
	"bytes"
	"context"
	models "elastic/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"

	"github.com/gin-gonic/gin"
)

func GetESClient() *elasticsearch.Client {
	/* Fetching elastic Search Client */
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		panic(err)
	}
	return client

}


func CreateUser(c *gin.Context) {
	client := GetESClient()
	// Create a new user.
	var users models.User
	if err := c.BindJSON(&users); err != nil {
		panic(err)
	}

	// Convert the user to JSON.
	userJSON, err := json.Marshal(&users)
	if err != nil {
		panic(err)
	}

	
	index := c.Param("index")
	// Create a new document in the "users" index.
	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: strconv.Itoa(users.ID),
		Body:       strings.NewReader(string(userJSON)),
		Refresh:    "false",
	}

	res, err := req.Do(c, client)
	if err != nil {
		// handle error
	}
	defer res.Body.Close()

	// Print the response.
	fmt.Println(users)
	c.JSON(http.StatusCreated, gin.H{
		// "id": users.ID,
		"name": users.Name,
		"age":  users.Age,
	})

}

func UpdateUser(c *gin.Context) {
	client := GetESClient()

	var users models.User
	if err := c.BindJSON(&users); err != nil {
		panic(err)
	}
	body := map[string]interface{}{
		"doc": map[string]interface{}{
			"name": users.Name,
			"age": users.Age,
		},
	}

	index := c.Param("index")
	id := c.Param("id")
	
	jsonBody, _ := json.Marshal(body)
	req := esapi.UpdateRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewReader(jsonBody),
	}
	res, _ := req.Do(context.Background(), client)
	defer res.Body.Close()
	fmt.Println(res.String())
	c.JSON(http.StatusCreated, gin.H{
		"id": id,
		"name": users.Name,
		"age":  users.Age,
	})
}

func DeleteUser(c *gin.Context) {
	client := GetESClient()
	id := c.Param("id")
	index := c.Param("index")

	// Set up the update request
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: id,
	}

	// Perform the update

	// Perform the update
	res, err := req.Do(context.Background(), client)
	if err != nil {
		// handle error
	}
	defer res.Body.Close()

	fmt.Println(res)
	c.JSON(http.StatusOK, "User deleted Successfully")

}



func GetUser(c *gin.Context) {
    var doc map[string]interface{}
    client := GetESClient()

    id := c.Param("id")
    index := c.Param("index")
    
    // Set up the update request
    getReq := esapi.GetRequest{
        Index:      index,
        DocumentID: id,
    }

    // Perform the update
    res, err := getReq.Do(context.Background(), client)
    if err != nil {
        // handle error
        panic(err)
    }

    fmt.Println(res)

    err = json.NewDecoder(res.Body).Decode(&doc)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing response body"})
        return
    }

    // Return the document in the response
    c.JSON(http.StatusOK, doc)
}




func GetAllUser(c *gin.Context) {
	client := GetESClient()

	index := c.Param("index")


	// Set up the update request
	req := esapi.SearchRequest{
		Index: []string{index},
	}

	// Perform the update
	res, err := req.Do(context.Background(), client)
	if err != nil {
		// handle error
		panic(err)
	}

	fmt.Println(res)
	var results struct {
		Hits struct {
			Hits []json.RawMessage `json:"hits"`
		} `json:"hits"`
	}

	err = json.NewDecoder(res.Body).Decode(&results)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing response body"})
		return
	}

	// Return the search results in the response
	c.JSON(http.StatusOK, results.Hits.Hits)
}




func CreateUserBatch(c *gin.Context) {
    client := GetESClient()

    // Create a new user
    var bulkRequest bytes.Buffer
    var users []models.User

    if err := c.BindJSON(&users); err != nil {
        panic(err)
    }

    for _, user := range users {
        // Convert the user to JSON
        userJSON, err := json.Marshal(&user)
        if err != nil {
            panic(err)
        }

        // Add the user to the Bulk request
        bulkRequest.Write(userJSON)
        bulkRequest.Write([]byte("\n"))

        req := esapi.IndexRequest{
            Index:      "bulk",
            DocumentID: strconv.Itoa(user.ID),
            Body:       strings.NewReader(string(userJSON)),
            Refresh:    "true",
        }

        // Send the Index request
        res, err := req.Do(context.Background(), client)
        if err != nil {
            panic(err)
        }

        defer res.Body.Close()

        // Print the response
        fmt.Println(res)
    }
    c.JSON(http.StatusCreated, gin.H{
        "msg": "user added successfully"})

}



func SearchUser(c *gin.Context)  {
	client := GetESClient()

	
	//var response map[string]interface{}
	var buf bytes.Buffer

	
	//limit, _ := strconv.Atoi("limit")

	/*
		Query sort in Elasticsearch.
		It will produce query like this:
		{
			"sort": {
				"_geo_distance": {
					"location": {
						"lat": splitLatLon[0],
						"lon": splitLatLon[1]
					},
					"order": __order__,
					"unit": __unit__
				}
			},
		}
	*/
	sort := map[string]interface{}{
	}

	name := c.Query("name")
	//age := c.Query("age")
	index:= c.Param("index")

	 page,_ := strconv.Atoi(c.Query("page"))
    size,_ := strconv.Atoi(c.Query("size"))

	// We encode from map string-interface into json format.
	if err := json.NewEncoder(&buf).Encode(sort); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// Process the query

	search := esapi.SearchRequest{
		Index: []string{index},
		// Query:   age,
		// Size: *size,
		Body:  strings.NewReader(`{"query":{"query_string":{"query": "` +name+ `"}},"from": ` + strconv.Itoa((page-1)*size) + `,
        "size": ` + strconv.Itoa(size) + `}`),

	}
	

	// Perform the update
	res, err := search.Do(context.Background(), client)
	if err != nil {
		// handle error
		panic(err)
	}

	fmt.Println(res)
	var results struct {
		Hits struct {
			Hits []json.RawMessage `json:"hits"`
		} `json:"hits"`
	}

	err = json.NewDecoder(res.Body).Decode(&results)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing response body"})
		return
	}

	// Return the search results in the response
	c.JSON(http.StatusOK, results.Hits.Hits)

}