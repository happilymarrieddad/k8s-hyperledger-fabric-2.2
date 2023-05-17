package jwt

import (
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"log"
	"time"

	jwtpkg "github.com/dgrijalva/jwt-go"
)

const (
	privKeyPath  = "../../keys/private-key.pem"
	pubKeyPath   = "../../keys/public-key.pem"
	HOURS_IN_DAY = 24
	DAYS_IN_WEEK = 7
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func init() {
	signBytes, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		panic(err)
	}
	signKey, err = jwtpkg.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		panic(err)
	}
	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		panic(err)
	}
	verifyKey, err = jwtpkg.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		panic(err)
	}
}

func GetToken(user, org string, identityBts []byte) string {
	token := jwtpkg.New(jwtpkg.SigningMethodRS512)
	claims := make(jwtpkg.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * HOURS_IN_DAY * DAYS_IN_WEEK).Unix()
	claims["iat"] = time.Now().Unix()
	claims["user"] = user
	claims["org"] = org
	if identityBts == nil {
		identityBts = []byte(`{}`)
	}
	claims["identityBts"] = string(identityBts)
	token.Claims = claims

	tokenString, _ := token.SignedString(signKey)

	return tokenString
}

func IsTokenValid(val string) (string, string, []byte, error) {
	token, err := jwtpkg.Parse(val, func(token *jwtpkg.Token) (interface{}, error) {
		return verifyKey, nil
	})

	switch e := err.(type) {
	case nil:
		if !token.Valid {
			return "", "", nil, errors.New("token is invalid")
		}

		claims, ok := token.Claims.(jwtpkg.MapClaims)
		if !ok {
			return "", "", nil, errors.New("token is invalid")
		}

		user := claims["user"].(string)
		org := claims["org"].(string)
		identityBts := claims["identityBts"].(string)

		return user, org, []byte(identityBts), nil
	case *jwtpkg.ValidationError:
		switch e.Errors {
		case jwtpkg.ValidationErrorExpired:
			return "", "", nil, errors.New("token expired, get a new one")
		default:
			log.Println(e)
			return "", "", nil, errors.New("error while parsing token")
		}
	default:
		return "", "", nil, errors.New("unable to parse token")
	}
}
