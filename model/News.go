package model

type NewsShort struct {
	Title string `gorm:"not null" json:"title"`
	Text string `gorm:"not null" json:"text"`
	PathToPicture string `gorm:"not null" json:"path_to_picture"`
}

type News struct {
	Common

	NewsShort
}

