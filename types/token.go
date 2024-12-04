package types

import "github.com/golang-jwt/jwt/v5"

type BearerTokenClaims struct {
	jwt.RegisteredClaims
	UID      uint64 `json:"uid"`
	Username string `json:"username"`
}
