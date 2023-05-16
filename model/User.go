package model

type UserShort struct {
	Username string `gorm:"not null;unique_index" json:"username" binding:"required"`
	Password string `gorm:"not null" json:"password,omitempty" binding:"required"`
	Email    string `gorm:"not null;unique_index" json:"email" binding:"required"`
	Role     string `gorm:"not null" json:"role,omitempty"`

	AccessToken  string `gorm:"-" json:"-"`
	RefreshToken string `gorm:"-" json:"-"`
}

type User struct {
	Common

	UserShort

	Reviews []Review `gorm:"foreignKey:UserID" json:"reviews,omitempty"`
}
