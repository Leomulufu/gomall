package service

import (
	"context"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	//"github.com/cloudwego/biz-demo/gomall/app/user/biz/dal/mysql"
	//"github.com/cloudwego/biz-demo/gomall/app/user/biz/model"
	user "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/user"
	"github.com/cloudwego/kitex/pkg/klog"
)

// Secret key used to sign the JWT token
var jwtSecret = []byte("your_secret_key")

type AuthService struct {
	ctx context.Context
}

func NewAuthService(ctx context.Context) *AuthService {
	return &AuthService{ctx: ctx}
}

// Claims represents the claims structure in the JWT token
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken generates a new JWT token for the given user ID
func (s *AuthService) GenerateToken(userID int) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token expires in 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(jwtSecret)
	if err != nil {
		klog.Errorf("Failed to generate token: %v", err)
		return "", err
	}

	return ss, nil
}

// ValidateToken validates the given JWT token and returns the user ID if valid
func (s *AuthService) ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		klog.Errorf("Failed to validate token: %v", err)
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// LoginWithToken combines login and token generation
func (s *AuthService) LoginWithToken(req *user.LoginReq) (resp *user.LoginResp, token string, err error) {
	loginService := NewLoginService(s.ctx)
	loginResp, err := loginService.Run(req)
	if err != nil {
		return nil, "", err
	}

	token, err = s.GenerateToken(int(loginResp.UserId))
	if err != nil {
		return nil, "", err
	}

	return loginResp, token, nil
}
