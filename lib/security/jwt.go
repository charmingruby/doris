package security

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWT struct {
	issuer    string
	secretKey string
}

type TokenClaim struct {
	Payload Payload `json:"payload"`
	jwt.StandardClaims
}

func NewJWT(issuer, secretKey string) *JWT {
	return &JWT{
		issuer:    issuer,
		secretKey: secretKey,
	}
}

func (s *JWT) Generate(sub string, p Payload) (string, error) {
	// TODO: change this to 1 hour, but for development purposes, we use 1 week
	// tokenDuration := time.Duration(time.Hour * 1) // 1 hour

	tokenDuration := time.Duration(time.Hour * 24 * 7) // 1 week

	claims := &TokenClaim{
		p,
		jwt.StandardClaims{
			Subject:   sub,
			Issuer:    s.issuer,
			ExpiresAt: time.Now().Local().Add(tokenDuration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", nil
	}

	return tokenStr, nil
}

func (j *JWT) Validate(token string) (string, Payload, error) {
	jwtToken, err := jwt.Parse(token, j.isTokenValid)
	if err != nil {
		return "", Payload{}, err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", Payload{}, fmt.Errorf("unable to parse jwt claims")
	}

	payloadClaims, ok := claims["payload"].(map[string]any)
	if !ok {
		return "", Payload{}, fmt.Errorf("payload is missing")
	}

	payload := Payload{
		Tier: payloadClaims["tier"].(string),
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", Payload{}, fmt.Errorf("subject is missing")
	}

	return sub, payload, nil
}

func (j *JWT) isTokenValid(t *jwt.Token) (any, error) {
	if _, isValid := t.Method.(*jwt.SigningMethodHMAC); !isValid {
		return nil, fmt.Errorf("invalid token %v", t)
	}

	return []byte(j.secretKey), nil
}
