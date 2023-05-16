package model

type CartShort struct {
	UserID uint64 `json:"user_id"`
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty" binding:"-"`

	BookID uint64 `gorm:"not null" json:"book_id"`
	Book   Book   `gorm:"foreignKey:BookID" json:"book,omitempty" binding:"-"`
}

type Cart struct {
	Common

	CartShort
}
