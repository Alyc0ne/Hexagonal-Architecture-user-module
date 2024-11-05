package domain

type (
	LoginUserReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LoginUserRes struct {
		ID           string `json:"id"`
		Email        string `json:"email"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	CreateUserReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	ForgetPasswordReq struct {
		Email string `json:"email"`
	}

	ForgetPasswordRes struct {
		ResetToken string `json:"reset_token"`
	}

	ResetPasswordReq struct {
		ResetToken string `json:"reset_token"`
		Password   string `json:"password"`
	}

	AccessToken struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
)
