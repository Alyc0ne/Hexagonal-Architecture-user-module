package domain

type (
	User struct {
		ID       string `json:"id" gorm:"id"`
		Email    string `json:"email" gorm:"email"`
		Password string `json:"password" gorm:"password"`
	}

	ForgetPassword struct {
		ID         string `json:"id" gorm:"id"`
		Email      string `json:"email" gorm:"email"`
		ResetToken string `json:"reset_token" gorm:"reset_token"`
	}
)
