package services

import (
	"encoding/json"
	"github.com/ayo-69/stage-2/models"
	"github.com/ayo-69/stage-2/utils"
	"math/rand"
	"net/http"
	"os"
	"time"

	"gorm.io/gorm"
)

type CountryRaw struct {
	Name       string `json:"name"`
	Capital    string `json:"capital"`
	Region     string `json:"region"`
	Population int    `json:"population"`
	Flag       string `json:"flag"`
	Currencies []struct {
		Code string `json:"code"`
	} `json:"currencies"`
}

func RefreshCountries(db *gorm.DB) error {
	// Fetch countries
	resp, err := http.Get("https://restcountries.com/v2/all?fields=name,capital,region,population,flag,currencies")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var rawCountries []CountryRaw
	if err := json.NewDecoder(resp.Body).Decode(&rawCountries); err != nil {
		return err
	}

	// Fetch exchange rates
	rates, err := GetExchangeRates()
	if err != nil {
		return err
	}

	now := time.Now()

	for _, raw := range rawCountries {
		if raw.Name == "" || raw.Population == 0 {
			continue
		}

		currencyCode := ""
		if len(raw.Currencies) > 0 {
			currencyCode = raw.Currencies[0].Code
		}

		var exchangeRate *float64
		var estimatedGDP *float64

		if currencyCode != "" {
			if rate, ok := rates[currencyCode]; ok && rate > 0 {
				r := &rate
				exchangeRate = r
				multiplier := rand.Float64()*(2000-1000) + 1000
				gdp := float64(raw.Population) * multiplier / rate
				estimatedGDP = &gdp
			}
		} else {
			gdp := 0.0
			estimatedGDP = &gdp
		}

		var country models.Country
		result := db.Where("LOWER(name) = LOWER(?)", raw.Name).First(&country)
		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			continue
		}

		country.Name = raw.Name
		country.Capital = nil
		if raw.Capital != "" {
			country.Capital = &raw.Capital
		}
		country.Region = nil
		if raw.Region != "" {
			country.Region = &raw.Region
		}
		country.Population = raw.Population
		country.CurrencyCode = nil
		if currencyCode != "" {
			country.CurrencyCode = &currencyCode
		}
		country.ExchageRate = exchangeRate
		country.EstimateGDP = estimatedGDP
		country.FlagURL = nil
		if raw.Flag != "" {
			country.FlagURL = &raw.Flag
		}
		country.LastRefreshedAt = now

		if result.Error == gorm.ErrRecordNotFound {
			db.Create(&country)
		} else {
			db.Save(&country)
		}
	}

	// Update status
	var status models.Status
	db.FirstOrCreate(&status, models.Status{ID: 1})
	status.LastRefreshedAt = now
	db.Save(&status)

	// Invalidate cache
	utils.Set("countries", nil, 0)

	// Generate image
	return GenerateSummaryImage(db)
}

func GetCountries(db *gorm.DB, region, currency, sort string) ([]map[string]interface{}, error) {
	cacheKey := "countries_list_" + region + "_" + currency + "_" + sort
	if cached, ok := utils.Get(cacheKey); ok {
		return cached.([]map[string]interface{}), nil
	}

	var countries []models.Country
	query := db.Model(&models.Country{})

	if region != "" {
		query = query.Where("region = ?", region)
	}
	if currency != "" {
		query = query.Where("currency_code = ?", currency)
	}

	if sort == "gdp_desc" {
		query = query.Order("estimated_gdp DESC")
	} else if sort == "gdp_asc" {
		query = query.Order("estimated_gdp ASC")
	}

	if err := query.Find(&countries).Error; err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(countries))
	for i, c := range countries {
		result[i] = map[string]interface{}{
			"id":                c.ID,
			"name":              c.Name,
			"capital":           c.Capital,
			"region":            c.Region,
			"population":        c.Population,
			"currency_code":     c.CurrencyCode,
			"exchange_rate":     c.ExchageRate,
			"estimated_gdp":     c.EstimateGDP,
			"flag_url":          c.FlagURL,
			"last_refreshed_at": c.LastRefreshedAt.Format(time.RFC3339),
		}
	}

	ttl := time.Duration(3600) * time.Second
	if seconds := os.Getenv("CACHE_TTL_SECONDS"); seconds != "" {
		if s, err := time.ParseDuration(seconds + "s"); err == nil {
			ttl = s
		}
	}
	utils.Set(cacheKey, result, ttl)
	return result, nil
}

func GetCountryByName(db *gorm.DB, name string) (*map[string]interface{}, error) {
	var country models.Country
	if err := db.Where("LOWER(name) = LOWER(?)", name).First(&country).Error; err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"id":                country.ID,
		"name":              country.Name,
		"capital":           country.Capital,
		"region":            country.Region,
		"population":        country.Population,
		"currency_code":     country.CurrencyCode,
		"exchange_rate":     country.ExchageRate,
		"estimated_gdp":     country.EstimateGDP,
		"flag_url":          country.FlagURL,
		"last_refreshed_at": country.LastRefreshedAt.Format(time.RFC3339),
	}
	return &result, nil
}

func DeleteCountryByName(db *gorm.DB, name string) error {
	result := db.Where("LOWER(name) = LOWER(?)", name).Delete(&models.Country{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	// Invalidate all list caches
	utils.Set("countries", nil, 0)
	return nil
}

func GetStatus(db *gorm.DB) (map[string]interface{}, error) {
	count := int64(0)
	db.Model(&models.Country{}).Count(&count)

	var status models.Status
	db.First(&status, 1)

	resp := map[string]interface{}{
		"total_countries":   count,
		"last_refreshed_at": nil,
	}
	if !status.LastRefreshedAt.IsZero() {
		resp["last_refreshed_at"] = status.LastRefreshedAt.Format(time.RFC3339)
	}
	return resp, nil
}
