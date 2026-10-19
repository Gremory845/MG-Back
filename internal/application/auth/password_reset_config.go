package auth

import "time"

// Parámetros de seguridad de la recuperación de contraseña.
const (
	resetCodeTTL     = 10 * time.Minute // vigencia del código de 6 dígitos
	resetTokenTTL    = 10 * time.Minute // tiempo para escribir la nueva contraseña tras verificar el código
	resetMaxAttempts = 5                // intentos de código por solicitud; luego hay que pedir otro
	resetCooldown    = 60 * time.Second // espera mínima entre dos correos para el mismo usuario
	resetMaxPerHour  = 5                // máximo de correos por usuario en una hora
	resetMailTimeout = 30 * time.Second
)
