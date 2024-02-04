package token

import "github.com/lestrrat-go/jwx/v2/jwt"

func GetEmail(token jwt.Token) string {
	return GetTraits(token)["email"].(string)
}

func GetIdentityID(token jwt.Token) string {
	return GetIdentity(token)["id"].(string)
}

func GetClaims(token jwt.Token) map[string]interface{} {
	return token.PrivateClaims()
}

func GetSession(token jwt.Token) map[string]interface{} {
	return GetClaims(token)["session"].(map[string]interface{})
}

func GetIdentity(token jwt.Token) map[string]interface{} {
	return GetSession(token)["identity"].(map[string]interface{})
}

func GetTraits(token jwt.Token) map[string]interface{} {
	return GetIdentity(token)["traits"].(map[string]interface{})
}
