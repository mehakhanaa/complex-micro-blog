package parsers

import (
	"github.com/golang-jwt/jwt/v5"

	"github.com/mehakhanaa/complex-micro-blog/consts"
	"github.com/mehakhanaa/complex-micro-blog/types"
)

func ParseToken(token string) (*types.BearerTokenClaims, error) {

	claims := new(types.BearerTokenClaims)
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(consts.TOKEN_SECRET), nil
	})

	return claims, err
}
