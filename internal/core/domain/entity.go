package domain

type (
	CountResult struct {
		Count int `gorm:"column:count"`
	}

	User struct {
		ID       string `json:"id" gorm:"primaryKey;type:varchar(50); not null"`
		Email    string `json:"email" gorm:"type:varchar(100); not null"`
		Password string `json:"password" gorm:"type:text; not null"`
	}

	ForgetPassword struct {
		ID         string `json:"id" gorm:"id"`
		Email      string `json:"email" gorm:"email"`
		ResetToken string `json:"reset_token" gorm:"reset_token"`
	}
)
