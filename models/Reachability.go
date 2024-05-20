package models

import "time"

type Reachability struct {
	URL                  string    `gorm:"column:URL;primaryKey"`
	IPv4FirstLevelReach  float64   `gorm:"column:IPv4FirstLevelReach"`
	IPv4SecondLevelReach float64   `gorm:"column:IPv4SecondLevelReach"`
	IPv6FirstLevelReach  float64   `gorm:"column:IPv6FirstLevelReach"`
	IPv6SecondLevelReach float64   `gorm:"column:IPv6SecondLevelReach"`
	CreateTime           time.Time `gorm:"column:create_time;primaryKey;default:CURRENT_TIMESTAMP"`
}

// TableName sets the table name for the model.
func (Reachability) TableName() string {
	return "reachability"
}
