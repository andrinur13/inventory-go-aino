package middleware

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"twc-ota-api/db/entities"

	"github.com/dgrijalva/jwt-go"
)

//JwtDuration : duration for generated JWT token
var JwtDuration = time.Hour * 24

//CreateJwtToken : generate JWT token
func CreateJwtToken(data *entities.Users) (string, error) {

	/* Create the token */
	token := jwt.New(jwt.SigningMethodHS256)

	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)

	/* Set token claims */
	claims["authorized"] = true
	claims["user"] = data.Name
	// claims["user"] = "useragent"
	claims["user_id"] = data.ID
	claims["email"] = data.Email
	claims["agent_id"] = data.Typeid
	// exp, _ := time.ParseDuration(os.Getenv("JWT_EXPIRED_DURATION"))
	exp := JwtDuration
	claims["exp"] = time.Now().Add(exp).Unix()

	/* Sign the token with our secret */
	// tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	tokenString, err := token.SignedString([]byte("TnTSys+3m!8"))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

//Authorize : auth JWT Token
func Authorize(tokenString string) (string, error) {
	result := strings.Split(tokenString, " ")
	if len(result) == 1 || strings.ToLower(result[0]) != "bearer" {
		return "", errors.New("Invalid token")
	}

	tokenString = result[1]
	_, err := jwt.Parse(result[1], func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// return []byte(os.Getenv("JWT_SECRET")), nil
		return []byte("TnTSys+3m!8"), nil
	})

	return result[1], err
}

//Decode : decode JWT Token
func Decode(tokenString string) *entities.Users {

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("TnTSys+3m!8"), nil
	})

	if err != nil {
		panic(err)
	}

	user := claims["user"].(string)
	email := claims["email"].(string)
	agentID := claims["agent_id"].(float64)
	userID := claims["user_id"].(float64)

	resp := entities.Users{
		Name: user,
		Email:  email,
		Typeid: int(agentID),
		ID:     int(userID),
	}

	return &resp
}
