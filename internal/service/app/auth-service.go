package app

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"github.com/imarrche/jwt-auth-example/internal/config"
	"github.com/imarrche/jwt-auth-example/internal/model"
	"github.com/imarrche/jwt-auth-example/internal/store"
)

// authService implements authorization business logic.
type authService struct {
	store store.Store
}

// newAuthServer creates and returns a new authService instance.
func newAuthService(s store.Store) *authService {
	return &authService{store: s}
}

// Sign up signes up a user.
func (s *authService) SignUp(u model.User) (model.User, error) {
	if err := u.Validate(); err != nil {
		return model.User{}, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}
	u.Password = ""
	u.PasswordHash = string(hashedPassword)

	u, err = s.store.Users().Create(u)
	if err != nil {
		return model.User{}, err
	}

	return u, err
}

// generateJWT generates access/refresh JSON Web Token for user.
func (s *authService) generateJWT(userID int, tokenType string) (string, error) {
	secret := []byte(config.Get().JWT.Secret)

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"type": tokenType, "user_id": userID, "exp": time.Now().Add(time.Hour * 24).Unix(),
	}).SignedString(secret)
	if err != nil {
		return "", err
	}

	return token, nil
}

// SignIn returns access and refresh JSON Web Tokens for a user if credentials are valid.
func (s *authService) SignIn(email string, password string) (string, string, error) {
	u, err := s.store.Users().GetByEmail(email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	accessJWT, err := s.generateJWT(u.ID, "access")
	if err != nil {
		return "", "", err
	}
	refreshJWT, err := s.generateJWT(u.ID, "refresh")
	if err != nil {
		return "", "", err
	}

	return accessJWT, refreshJWT, err
}

// ValidateJWT validates access JSON Web Token and returns error if it's invalid.
func (s *authService) ValidateJWT(token, tokenType string) (int, error) {
	secret := []byte(config.Get().JWT.Secret)

	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected token signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if ok && t.Valid {
		if tType, ok := claims["type"].(string); ok {
			if tType != tokenType {
				return 0, errors.New("invalid JWT type")
			}
		} else {
			return 0, errors.New("couldn't parse JWT's type")
		}

		expTime, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["exp"]), 10, 64)
		if err != nil {
			return 0, errors.New("couldn't parse JWT's expiration time")
		}
		if time.Unix(expTime, 0).Before(time.Now()) {
			return 0, errors.New("JWT expired")
		}

		userID, err := strconv.ParseInt(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
		if err != nil {
			return 0, errors.New("couldn't parse JWT's user ID")
		}
		return int(userID), nil
	}

	return 0, errors.New("JWT is invalid")
}

// RefreshAccessJWT returns new access JSON Web Token if valid refresh token was provided.
func (s *authService) RefreshAccessJWT(refreshToken string) (string, error) {
	userID, err := s.ValidateJWT(refreshToken, "refresh")
	if err != nil {
		return "", err
	}

	return s.generateJWT(userID, "access")
}
