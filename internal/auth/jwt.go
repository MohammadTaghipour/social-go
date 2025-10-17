package auth

import "github.com/golang-jwt/jwt/v5"

type JwtAuthenticator struct {
	secret string
	aud    string
	iss    string
}

func NewJwtAuthenticator(secret, audience, issuer string) *JwtAuthenticator {
	return &JwtAuthenticator{
		secret: secret,
		aud:    audience,
		iss:    issuer,
	}
}

func (a *JwtAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *JwtAuthenticator) ValidateToken(tokenString string) (*jwt.Token, error) {
	return nil, nil
}
