package models

type Website struct {
	ID        uint   `gorm:"primaryKey"`
	ImageURL  string `gorm:"column:image_url;not null"`
	IsActive  bool   `gorm:"column:is_active;default:true"`
	LastCheck string `gorm:"column:last_check"`
}
