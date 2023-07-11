package helpers

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

// JwtWrapper wraps the signing key and the issuer
type JwtWrapper struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

// JwtClaim adds portal userid  as a claim to the token
type JwtClaim struct {
	UserID string
	jwt.StandardClaims
}

func (j *JwtWrapper) GenerateJWTToken(userid string) (signedToken string, err error) {
	claims := &JwtClaim{
		UserID: userid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(j.ExpirationHours)).Unix(),
			Issuer:    j.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return signedToken, err
	}
	return signedToken, err
}

// ValidateToken validates the jwt token
func (j *JwtWrapper) ValidateJWTToken(signedToken string) (claims *JwtClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*JwtClaim)
	if !ok {
		err = errors.New("unable to parse token claim")
		return &JwtClaim{}, err
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token has expired")
		return &JwtClaim{}, err
	}
	return claims, err
}
