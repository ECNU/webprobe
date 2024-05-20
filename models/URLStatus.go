package models

import "time"

type URLStatus struct {
	FatherURL  string    `gorm:"column:father_url"`
	Depth      int       `gorm:"column:depth"`
	URL        string    `gorm:"column:URL;primaryKey"`
	IPVersion  string    `gorm:"column:ip_version;primaryKey"`
	Up         bool      `gorm:"column:up"`
	StatusCode int       `gorm:"column:status_code"`
	Latency    int       `gorm:"column:latency"`
	CertExpire time.Time `gorm:"column:cert_expire"`
	CreateTime time.Time `gorm:"column:create_time;primaryKey;default:CURRENT_TIMESTAMP"`
}

// TableName sets the table name for the model.
func (URLStatus) TableName() string {
	return "url_status"
}
