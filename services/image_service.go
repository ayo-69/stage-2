package services

import (
	"fmt"
	"github.com/ayo-69/stage-2/models"
	"os"
	"time"

	"github.com/fogleman/gg"
	"gorm.io/gorm"
)

func GenerateSummaryImage(db *gorm.DB) error {
	var count int64
	db.Model(&models.Country{}).Count(&count)

	var top5 []models.Country
	db.Order("estimated_gdp DESC").Limit(5).Find(&top5)

	var status models.Status
	db.First(&status, 1)

	const width, height = 800, 600
	dc := gg.NewContext(width, height)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	if err := dc.LoadFontFace("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 32); err != nil {
		// fallback
		dc.SetRGB(0, 0, 0)
	}

	dc.SetRGB(0, 0, 0)
	dc.DrawStringAnchored("Country Summary", 20, 40, 0, 1)

	dc.DrawStringAnchored("Total countries: "+string(rune(count)), 20, 100, 0, 1)
	dc.DrawStringAnchored("Top 5 by estimated GDP:", 20, 160, 0, 1)

	y := 220.0
	for i, c := range top5 {
		gdp := "N/A"
		if c.EstimateGDP != nil {
			gdp = fmt.Sprintf("%.2f", *c.EstimateGDP)
		}
		dc.DrawStringAnchored(fmt.Sprintf("%d. %s: %s", i+1, c.Name, gdp), 40, y, 0, 1)
		y += 50
	}

	ts := "N/A"
	if !status.LastRefreshedAt.IsZero() {
		ts = status.LastRefreshedAt.Format(time.RFC3339)
	}
	dc.DrawStringAnchored("Last refresh: "+ts, 20, height-60, 0, 1)

	os.MkdirAll("cache", 0755)
	return dc.SavePNG("cache/summary.png")
}
