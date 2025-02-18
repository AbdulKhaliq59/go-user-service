package config

import (
	"log"

	"gorm.io/gorm"
)

func listAllTables(db *gorm.DB) {
	var tables []string
	query := `SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'`
	err := db.Raw(query).Scan(&tables).Error
	if err != nil {
		log.Printf("Error fetching tables: %v", err)
	} else {
		log.Println("Available tables in the database:", tables)
	}
}
