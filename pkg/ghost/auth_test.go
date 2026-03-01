package ghost

import (
	"encoding/hex"
	"testing"
)

func TestGenerateJWT(t *testing.T) {
	// test key: known id and hex-encoded secret
	id := "abcdef0123456789abcdef01"
	secret := hex.EncodeToString([]byte("test-secret-key-value!!!"))
	key := id + ":" + secret

	token, err := GenerateJWT(key)
	if err != nil {
		t.Fatalf("GenerateJWT failed: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	// verify the token round-trips
	claims, err := VerifyJWT(token, secret)
	if err != nil {
		t.Fatalf("VerifyJWT failed: %v", err)
	}

	aud, ok := claims["aud"].(string)
	if !ok || aud != "/admin/" {
		t.Errorf("expected aud=/admin/, got %v", claims["aud"])
	}

	if _, ok := claims["iat"]; !ok {
		t.Error("expected iat claim")
	}
	if _, ok := claims["exp"]; !ok {
		t.Error("expected exp claim")
	}
}

func TestGenerateJWT_InvalidFormat(t *testing.T) {
	_, err := GenerateJWT("not-a-valid-key")
	if err == nil {
		t.Fatal("expected error for invalid key format")
	}
}

func TestGenerateJWT_InvalidHex(t *testing.T) {
	_, err := GenerateJWT("abcdef0123456789abcdef01:not-hex!")
	if err == nil {
		t.Fatal("expected error for non-hex secret")
	}
}

func TestValidateKeyFormat(t *testing.T) {
	tests := []struct {
		key     string
		wantErr bool
	}{
		{"abc123:aabbcc", false},
		{"abc123:", true},
		{":aabbcc", true},
		{"no-colon", true},
		{"abc123:not-hex!", true},
	}

	for _, tt := range tests {
		err := ValidateKeyFormat(tt.key)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateKeyFormat(%q) error=%v, wantErr=%v", tt.key, err, tt.wantErr)
		}
	}
}
