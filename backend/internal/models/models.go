package models

import (
	"time"
)

type District struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	Name      string  `gorm:"unique;not null" json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type CrimeReport struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	DistrictID   uint      `json:"district_id"`
	District     District  `gorm:"foreignKey:DistrictID" json:"district"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	IncidentDate time.Time `json:"incident_date"`
	SourceURL    string    `gorm:"unique" json:"source_url"`
	Category     string    `json:"category"`
	CreatedAt    time.Time `json:"created_at"`
}
