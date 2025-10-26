package controllers

import (
	"github.com/ayo-69/stage-2/services"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RefreshCountries(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := services.RefreshCountries(db); err != nil {
			if os.IsTimeout(err) || err.Error() == "EOF" {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"error":   "External data source unavailable",
					"details": "Could not fetch data from " + err.Error(),
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Refresh completed successfully"})
	}
}

func GetCountries(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		region := c.Query("region")
		currency := c.Query("currency")
		sort := c.Query("sort")

		countries, err := services.GetCountries(db, region, currency, sort)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		c.JSON(http.StatusOK, countries)
	}
}

func GetCountry(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		country, err := services.GetCountryByName(db, name)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Country not found"})
			return
		}
		c.JSON(http.StatusOK, country)
	}
}

func DeleteCountry(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		if err := services.DeleteCountryByName(db, name); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Country not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Country deleted successfully"})
	}
}

func GetStatus(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		status, err := services.GetStatus(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		c.JSON(http.StatusOK, status)
	}
}

func GetSummaryImage(c *gin.Context) {
	path := "cache/summary.png"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Summary image not found"})
		return
	}
	c.File(path)
}
