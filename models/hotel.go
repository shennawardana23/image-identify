package models

type Hotel struct {
	HotelID     int    `gorm:"column:hotel_id;primaryKey;table:tb_hotels"`
	BrandID     int    `gorm:"column:brand_id"`
	HotelName   string `gorm:"column:hotel_name"`
	WebsiteLink string `gorm:"column:website_link"`
	Brand       Brand  `gorm:"foreignKey:BrandID"`
}

func (Hotel) TableName() string {
	return "tb_hotels"
}

type Brand struct {
	BrandID   int    `gorm:"column:brand_id;primaryKey"`
	BrandName string `gorm:"column:brand_name"`
}

func (Brand) TableName() string {
	return "tb_brands"
}
