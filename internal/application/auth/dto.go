package auth

type RegisterUserInput struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=30"`
	Password string `json:"password" validate:"required,min=8"`
}

type RegisterUserOutput struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type LoginInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ---- Recuperación de contraseña ----

type ForgotPasswordInput struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyResetCodeInput struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6,number"`
}

type VerifyResetCodeOutput struct {
	ResetToken string `json:"reset_token"`
}

type ResetPasswordInput struct {
	ResetToken      string `json:"reset_token" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}
