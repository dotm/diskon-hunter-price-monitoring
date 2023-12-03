package jwttoken

import (
	"fmt"
	"strings"

	"github.com/lestrrat-go/jwx/v2/jwt"
)

const diskonHunterCustomClaimName = "dhcc"
const diskonHunterCustomClaimKeyUserId = "user_id"

type CustomClaimMap = map[string]interface{}

func GetUserId(token jwt.Token) string {
	return GetCustomClaims(token)[diskonHunterCustomClaimKeyUserId].(string)
}
func GetCustomClaims(token jwt.Token) CustomClaimMap {
	return token.PrivateClaims()[diskonHunterCustomClaimName].(CustomClaimMap)
}

func ParseJWTFromString(jwtTokenString string) jwt.Token {
	jwtToken, err := jwt.Parse([]byte(CleanUpJWT(jwtTokenString)), jwt.WithVerify(false))
	if err != nil {
		err := fmt.Errorf("error parsing jwtToken %s", err)
		panic(err)
	}

	return jwtToken
}

func CleanUpJWT(jwtToken string) string {
	if jwtToken == "" {
		return ""
	}
	jwtToken = strings.TrimPrefix(jwtToken, "jwt=")
	jwtToken = strings.Split(jwtToken, " ")[0]
	jwtToken = strings.TrimSuffix(jwtToken, ";")
	return jwtToken
}
