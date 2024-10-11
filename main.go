package main

import (
	"accounts/api"
	"accounts/db"
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	db := db.NewDB(os.Args[1], os.Args[2], os.Args[3], os.Args[4])
	defer db.Close(context.Background())

	router := gin.Default()
	accountAPI := *NewAccountAPI(&db)

	api.RegisterHandlers(router, accountAPI)

	fmt.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		fmt.Printf("Failed to start server: %v", err)
	}
}
