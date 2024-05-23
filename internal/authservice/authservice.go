package authservice

import (
	"context"
	"errors"

	"authservice/internal/authorization"
	"authservice/internal/model"

	"github.com/redis/go-redis/v9"
)

type Users interface {
	ListUsers(ctx context.Context) ([]*model.User, error)
	DeleteUser(ctx context.Context, username string) error
}

type Tokens interface {
	GetUserToken(ctx context.Context, username string, password string) (jwtToken string, err error)
	CreateUser(ctx context.Context, username string, password string) error
	InvalidateToken(ctx context.Context, userID string) error
	ValidateToken(ctx context.Context, jwtToken string) (bool, error)
}
type AuthService interface {
	Users
	Tokens
}

type SessionStore interface {
	CreateSession(ctx context.Context, session *model.Session) error
	GetSession(ctx context.Context, uid string) (*model.Session, error)
	DeleteSession(ctx context.Context, uid string) error
}

type UserStore interface {
	GetUser(ctx context.Context, uid string) (*model.User, error)
	CreateOrSetUser(ctx context.Context, user *model.User) error
	ListUsers(ctx context.Context) ([]*model.User, error)
	DeleteUser(ctx context.Context, uid string) error
}

type Store interface {
	SessionStore
	UserStore
}

type Auth struct {
	validator authorization.Validator
	store     Store
}

func NewAuthService(validator authorization.Validator, store Store) *Auth {
	return &Auth{
		validator: validator,
		store:     store,
	}
}

func (a *Auth) GetUserToken(ctx context.Context, username string, password string) (string, error) {
	user, err := a.store.GetUser(ctx, username)
	if err != nil {
		return "", err
	}

	if !validatePassword(user.HashedPassword, password) {
		return "", errors.New("invalid password given")
	}

	jwtToken, err := a.validator.CreateToken(user.Username, nil)

	if err != nil {
		return "", err
	}

	err = a.store.CreateSession(ctx, &model.Session{
		Username: user.Username,
		JWTToken: jwtToken,
	})

	return jwtToken, err
}

func (a *Auth) InvalidateToken(ctx context.Context, username string) error {
	return a.store.DeleteSession(ctx, username)
}

func (a *Auth) ValidateToken(ctx context.Context, jwtToken string) (bool, error) {
	claimsFromToken, err := a.validator.Validate(jwtToken)
	if err != nil {
		return false, err
	}
	session, err := a.store.GetSession(ctx, claimsFromToken.UserID)
	if err != nil {
		return false, err
	}
	return session.JWTToken == jwtToken, nil
}

func (a *Auth) CreateUser(ctx context.Context, username string, password string) error {
	user, err := a.store.GetUser(ctx, username)
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	if user != nil {
		return errors.New("user already exists")
	}

	encryptedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}

	return a.store.CreateOrSetUser(ctx, &model.User{
		Username:       username,
		HashedPassword: encryptedPassword,
		LoggedIn:       false,
	})
}

func (a *Auth) ListUsers(ctx context.Context) ([]*model.User, error) {
	return a.store.ListUsers(ctx)
}

func (a *Auth) DeleteUser(ctx context.Context, username string) error {
	user, err := a.store.GetUser(ctx, username)
	if err != nil {
		return err
	}

	session, err := a.store.GetSession(ctx, user.Username)
	if err != nil {
		return err
	}
	if session != nil {
		a.store.DeleteSession(ctx, user.Username)
	}

	return a.store.DeleteUser(ctx, user.Username)
}
