package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokensecret string, expiresIn time.Duration) (string, error) {
        claims := jwt.RegisteredClaims{
                Issuer: "chirpy",
                IssuedAt: jwt.NewNumericDate(time.Now()),
                ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
                Subject: userID.String(),
        }
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
        signedJWT, err := token.SignedString([]byte(tokensecret)) 
        if err!= nil {
                return "", err
        }
        return signedJWT, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error){
        var claims jwt.RegisteredClaims
        token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token)(interface{}, error){
                return []byte(tokenSecret), nil
        }, jwt.WithLeeway(5*time.Second))
        if err != nil {
                return uuid.Nil, err
        }else if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok{
                fmt.Printf("%v", claims.Subject)
                return uuid.MustParse(claims.Subject), nil
        } else {
                return uuid.Nil, err
        }
}

func GetBearerToken(headers http.Header) (string, error){
        TokenString, _ := strings.CutPrefix(headers.Get("Authorization"), "Bearer ")
        if TokenString == "" {
                return "", fmt.Errorf("No token string received")
        }
        return TokenString, nil
}
