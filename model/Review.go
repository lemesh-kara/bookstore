package model

type ReviewShort struct {
	UserID uint64 `json:"user_id" binding:"required"`
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty" binding:"-"`

	BookID uint64 `gorm:"not null" json:"book_id" binding:"required"`
	Book   Book   `gorm:"foreignKey:BookID" json:"book,omitempty" binding:"-"`

	ReviewMark float64 `gorm:"not null;check:review_mark >= 0 AND review_mark <= 5" json:"review_mark" binding:"required"`
	Text       string  `gorm:"not null" json:"text" binding:"required"`
}

type Review struct {
	Common

	ReviewShort
}

// UserID  uint64  `json:"user_id"`
// User    User    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:UserID" json:"-"`
// UserID  uint64 `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
