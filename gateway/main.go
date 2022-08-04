package main

import (
	"fmt"

	middleware "gateway/middleware"

	routes "gateway/routes"

	config "gateway/config"

	"github.com/gin-gonic/gin"
)

func main() {
	config := config.LoadConfig()
	fmt.Printf("Config: %v\n", config)
	server := gin.Default()

	// health check
	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Routes for User Service
	userServiceGroup := server.Group("/user")
	routes.UserServiceRoutes(userServiceGroup)

	// Routes for Item Service
	itemServiceGroup := server.Group("/item")
	itemServiceGroup.Use(middleware.Authenticate())
	routes.ItemServiceRoutes(itemServiceGroup)

	err := server.Run(":5000")
	if err != nil {
		fmt.Println("Gateway server failed to start, err: %v\n", err)
	}
}

// func loadConfig() *Config {
// 	yfile, err := ioutil.ReadFile("users.yaml")

// 	if err != nil {

// 		log.Fatal(err)
// 	}

// 	data := make(map[string]User)

// 	err2 := yaml.Unmarshal(yfile, &data)

// 	if err2 != nil {

// 		log.Fatal(err2)
// 	}

// 	for k, v := range data {set

// 		fmt.Printf("%s: %s\n", k, v)
// 	}
// }
