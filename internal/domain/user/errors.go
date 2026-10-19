package user

import "mg-back/internal/domain/shared"

var (
	ErrUserNotFound       = shared.NewDomainError("USER_NOT_FOUND", "usuario no encontrado")
	ErrEmailAlreadyExists = shared.NewDomainError("EMAIL_ALREADY_EXISTS", "el correo ya está registrado")
	ErrUsernameTaken      = shared.NewDomainError("USERNAME_TAKEN", "el nombre de usuario ya está en uso")
	ErrInvalidCredentials = shared.NewDomainError("INVALID_CREDENTIALS", "credenciales inválidas")

	ErrRefreshTokenNotFound = shared.NewDomainError("REFRESH_TOKEN_NOT_FOUND", "token de actualización no encontrado")
	ErrRefreshTokenExpired  = shared.NewDomainError("REFRESH_TOKEN_EXPIRED", "el token de actualización expiró")
	ErrRefreshTokenReused   = shared.NewDomainError("REFRESH_TOKEN_REUSED", "se detectó reutilización de un token revocado; se cerraron todas las sesiones")

	// Un solo error genérico para código incorrecto / vencido / sin solicitud / demasiados intentos,
	// así nadie puede deducir cuál fue la causa (ni si el correo existe).
	ErrInvalidResetCode  = shared.NewDomainError("INVALID_RESET_CODE", "el código es inválido o expiró; solicita uno nuevo")
	ErrInvalidResetToken = shared.NewDomainError("INVALID_RESET_TOKEN", "la sesión de restablecimiento expiró; vuelve a empezar")
)
