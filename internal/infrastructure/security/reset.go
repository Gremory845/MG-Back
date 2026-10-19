package security

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

// GenerateResetCode devuelve un código de 6 dígitos (000000-999999) con
// distribución uniforme, usando el generador criptográfico del sistema.
func GenerateResetCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// HashResetCode calcula el HMAC-SHA256 del código atado al usuario.
// Un código de 6 dígitos solo tiene un millón de combinaciones: un hash simple se
// rompería en milisegundos si la base de datos se filtrara. Con HMAC y un secreto
// que vive solo en el servidor, quien robe la base no puede recuperar los códigos.
func HashResetCode(secret []byte, userID, code string) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte("password-reset:" + userID + ":" + code))
	return hex.EncodeToString(mac.Sum(nil))
}

// ResetCodeMatches compara en tiempo constante.
func ResetCodeMatches(secret []byte, userID, code, storedHash string) bool {
	return hmac.Equal([]byte(HashResetCode(secret, userID, code)), []byte(storedHash))
}

// GenerateResetToken crea el token opaco de un solo uso que se entrega al cliente
// tras verificar el código (se guarda solo su hash SHA-256).
func GenerateResetToken() (plain string, hash string, err error) {
	return GenerateRefreshToken()
}
