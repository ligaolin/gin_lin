package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ligaolin/gin_lin/global"
)

type MyCustomClaims struct {
	Id   string `json:"id"`
	Type string `json:"type"`
	jwt.RegisteredClaims
}

func JwtSet(id *uint32, user_type string) (string, error) {
	claims := MyCustomClaims{
		fmt.Sprintf("%d", *id),
		user_type,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(global.Config.Jwt.Expir) * time.Minute)),
			Issuer:    "lin",
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("AllYourBase"))
}

func JwtGet(t string) (MyCustomClaims, error) {
	var c MyCustomClaims
	token, err := jwt.ParseWithClaims(t, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("身份信无效或已过期")
		}
		return []byte("AllYourBase"), nil
	})
	if err != nil {
		return c, err
	}

	if claims, ok := token.Claims.(*MyCustomClaims); ok {
		return *claims, nil
	} else {
		return c, errors.New("解析身份信息失败")
	}
}
