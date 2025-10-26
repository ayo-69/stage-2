package models

import "time"

type Country struct {
	ID              uint      `gorm:"primaryKey" json"id"`
	Name            string    `gorm:"uniqueIndex;not null" json:"name"`
	Capital         *string   `json:"capital,omitempty"`
	Region          *string   `json:"region,omitempty"`
	Population      int       `gorm:"not null" json"population"`
	CurrencyCode    *string   `gorm:"size:10" json:"currency_code,omitempty"`
	ExchageRate     *float64  `json:"exchange_rate,omitemtpy"`
	EstimateGDP     *float64  `json:"estimated_gdp,omitemtpy"`
	FlagURL         *string   `json:"flag_url,omitempty"`
	LastRefreshedAt time.Time `gorm:"autoUpdateTime" json:"last_refreshed_at"`
}

type Status struct {
	ID              uint      `gorm:"primaryKey"`
	LastRefreshedAt time.Time `json:"last_refreshed_at,omitrempty"`
}
