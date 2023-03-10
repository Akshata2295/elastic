package routes

import (
	"fmt"
	controllers "elastic/controllers"
	"github.com/gin-gonic/gin"
)


func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Logger())


	r.Use(func(c *gin.Context) {
		//allow all
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers",  "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			platform := c.GetHeader("User-agent")
			fmt.Println("platofor", platform)
		}
		c.Next()
	})
	r.POST("/users/:index", controllers.CreateUser)
	r.PATCH("/users/:index/:id", controllers.UpdateUser)
	r.DELETE("/users/:index/:id", controllers.DeleteUser)
	r.GET("/users/:index/:id", controllers.GetUser)
	r.GET("/users/:index", controllers.GetAllUser)
	r.POST("/users", controllers.CreateUserBatch)
	r.GET("/users/:index/search",controllers.SearchUser)
	
	return r
}