// main.go
package main

import (
	"github.com/ayo-69/stage-2/controllers"
	"github.com/ayo-69/stage-2/models"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// PostgreSQL DSN from .env
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required in .env")
	}

	// Connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Optional: Disable foreign key constraints (if needed)
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}

	// Auto-migrate models
	if err := db.AutoMigrate(&models.Country{}, &models.Status{}); err != nil {
		log.Fatal("Migration failed:", err)
	}

	// Optional: Seed Status table
	var status models.Status
	if err := db.First(&status, 1).Error; err != nil {
		db.Create(&models.Status{ID: 1})
	}

	// Setup Gin router
	r := gin.Default()

	// Routes
	r.POST("/countries/refresh", controllers.RefreshCountries(db))
	r.GET("/countries", controllers.GetCountries(db))
	r.GET("/countries/:name", controllers.GetCountry(db))
	r.DELETE("/countries/:name", controllers.DeleteCountry(db))
	r.GET("/status", controllers.GetStatus(db))
	r.GET("/countries/image", controllers.GetSummaryImage)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Server failed:", err)
	}
}
