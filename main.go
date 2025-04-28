// main.go
package main

import (
	"go-user-service/api/routes"
	"go-user-service/config"
	"go-user-service/services"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "go-user-service/docs" // This line is important for swagger
)

// @title           LONDONWINNER
// @version         1.0
// @description     A backend service written in Go with Gin framework
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// No BasePath annotation here

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env")
	}
	log.Printf("KEYCLOAK_BASE_URL: %s", os.Getenv("KEYCLOAK_BASE_URL"))
	log.Printf("KEYCLOAK_REALM: %s", os.Getenv("KEYCLOAK_REALM"))
	log.Printf("KEYCLOAK_CLIENT_ID: %s", os.Getenv("KEYCLOAK_CLIENT_ID"))
	log.Printf("KEYCLOAK_CLIENT_SECRET: %s", os.Getenv("KEYCLOAK_CLIENT_SECRET")) // Don't log the actual secret

	// Initialize router
	router := gin.Default()

	router.Use(cors.Default())

	// Connect to database
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	keycloakService, err := services.NewKeycloakService()
	if err != nil {
		log.Fatalf("Failed to initialize Keycloak service: %v", err)
	}

	// Get the SQL DB object
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}
	defer sqlDB.Close()

	// API version group

	v1 := router.Group("/api/v1")
	{

		routes.SetupReportRoutes(v1, db)
		routes.SetupTransactionRoutes(v1, db)
		routes.SetupProductRoutes(v1, db)
		routes.SetupRoleRoutes(v1, db)
		routes.SetupGroupRoutes(v1, db)
		routes.SetupUserRoutes(v1, keycloakService)
		routes.SetupBlacklistRoutes(v1, db)
	}

	// Swagger documentation route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
