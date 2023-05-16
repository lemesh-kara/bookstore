package model

type BookShort struct {
	ISBN          string  `gorm:"uniqueIndex" json:"isbn"`
	Name          string  `gorm:"not null" json:"name"`
	Author        string  `gorm:"not null" json:"author"`
	Year          string  `gorm:"not null" json:"year"`
	PublisherData string  `gorm:"not null" json:"publisher_data"`
	Description   string  `gorm:"not null" json:"description"`
	PathToPdf     string  `gorm:"not null" json:"path_to_pdf"`
	PathToCover   string  `gorm:"not null" json:"path_to_cover"`
	ReviewMark    float64 `gorm:"not null;check:review_mark >= 0 AND review_mark <= 5" json:"review_mark"`
	Price         uint64  `gorm:"not null;check:price >= 0" json:"price"`
}

type Book struct {
	Common

	BookShort

	Reviews []Review `gorm:"foreignKey:BookID" json:"reviews,omitempty"`
}

func IsBookField(field string) bool {
	type VoidType struct{}
	void := VoidType{}
	fieldSet := map[string]VoidType{
		"isbn":           void,
		"id":             void,
		"name":           void,
		"author":         void,
		"year":           void,
		"publisher_data": void,
		"description":    void,
		"review_mark":    void,
		"created_at":     void,
		"updated_at":     void,
		"price":          void,
	}
	_, ok := fieldSet[field]
	return ok
}
