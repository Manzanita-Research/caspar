package ghost

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT creates a Ghost Admin API JWT from an id:secret key pair.
// The token uses HS256, has a 5-minute expiry, and targets the /admin/ audience.
func GenerateJWT(adminAPIKey string) (string, error) {
	parts := strings.SplitN(adminAPIKey, ":", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid admin API key format — expected id:secret")
	}

	id, secret := parts[0], parts[1]

	// Ghost expects the secret decoded from hex, used as HMAC-SHA256 key
	keyBytes, err := hex.DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("decoding API secret: %w", err)
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(5 * time.Minute).Unix(),
		"aud": "/admin/",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = id

	// Ghost uses the raw HMAC key, not a PEM-encoded key
	signed, err := token.SignedString(keyBytes)
	if err != nil {
		return "", fmt.Errorf("signing JWT: %w", err)
	}

	return signed, nil
}

// ValidateKeyFormat checks that an API key looks like id:secret without
// actually making a request.
func ValidateKeyFormat(key string) error {
	parts := strings.SplitN(key, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("invalid format — expected id:secret (find this in Ghost Admin → Settings → Integrations)")
	}

	// secret should be valid hex
	if _, err := hex.DecodeString(parts[1]); err != nil {
		return fmt.Errorf("secret is not valid hex — check your API key")
	}

	return nil
}

// VerifyJWT verifies a Ghost Admin API JWT. Used for testing.
func VerifyJWT(tokenString string, secret string) (jwt.MapClaims, error) {
	keyBytes, err := hex.DecodeString(secret)
	if err != nil {
		return nil, fmt.Errorf("decoding secret: %w", err)
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return keyBytes, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// hmacSHA256 is a helper for testing — not used in JWT flow directly.
func hmacSHA256(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}
