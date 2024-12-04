package generators

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/mehakhanaa/complex-micro-blog/consts"
	"github.com/mehakhanaa/complex-micro-blog/types"
)

func GenerateToken(uid uint64, username string) (string, *types.BearerTokenClaims, error) {

	claims := &types.BearerTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(consts.TOKEN_EXPIRE_DURATION * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    consts.TOKEN_ISSUER,
			Subject:   "BearerToken",
			ID:        uuid.New().String(),
		},
		UID:      uid,
		Username: username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(consts.TOKEN_SECRET))

	return tokenString, claims, err
}
