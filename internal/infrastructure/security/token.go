package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("token inválido o expirado")

type AccessClaims struct {
	UserID string `json:"sub"`
	jwt.RegisteredClaims
}

// GenerateAccessToken crea un JWT firmado con HS256, de corta duración.
// Esto es lo que el cliente manda en el header Authorization en cada
// petición a un endpoint protegido.
func GenerateAccessToken(userID string, secret []byte, ttl time.Duration) (string, error) {
	claims := AccessClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// ParseAccessToken valida la firma y expiración de un JWT y devuelve sus
// claims. Lo vamos a usar pronto en el middleware de autenticación.
func ParseAccessToken(tokenString string, secret []byte) (*AccessClaims, error) {
	claims := &AccessClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GenerateRefreshToken produce un token opaco aleatorio (lo que recibe
// el cliente) y su hash SHA-256 en hex (lo único que se guarda en BD).
// Nunca guardamos el token en texto plano — si la base de datos se
// filtra, los refresh tokens filtrados no sirven para nada sin el
// valor original.
func GenerateRefreshToken() (plain string, hash string, err error) {
	raw := make([]byte, 32)
	if _, err = rand.Read(raw); err != nil {
		return "", "", err
	}

	plain = base64.RawURLEncoding.EncodeToString(raw)
	hash = HashToken(plain)
	return plain, hash, nil
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
