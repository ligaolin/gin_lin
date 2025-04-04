package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtConfig struct {
	Expir  int64  `json:"expir" toml:"expir" yaml:"expir"` // jwt登录过期时间，分钟，1440一天
	Issuer string `json:"issuer" toml:"issuer" yaml:"issuer"`
}

type Claims struct {
	Id   string `json:"id"`
	Type string `json:"type"`
	jwt.RegisteredClaims
}

type Jwt struct {
}

func NewJwt(cfg JwtConfig) *Jwt {
	return &Jwt{}
}

func (j *Jwt) Set(id *uint32, types string, cfg JwtConfig) (string, error) {
	claims := Claims{
		fmt.Sprintf("%d", *id),
		types,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.Expir) * time.Minute)),
			Issuer:    cfg.Issuer,
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("AllYourBase"))
}

func (j *Jwt) Get(t string, claims *Claims) error {
	token, err := jwt.ParseWithClaims(t, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("身份信无效或已过期")
		}
		return []byte("AllYourBase"), nil
	})
	if err != nil {
		return err
	}

	if _, ok := token.Claims.(*Claims); ok {
		return nil
	} else {
		return errors.New("解析身份信息失败")
	}
}
