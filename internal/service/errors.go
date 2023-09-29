package service

import "errors"

var (
	ErrCannotSignToken  = errors.New("cannot sign token")
	ErrCannotParseToken = errors.New("cannot parse token")

	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrCannotCreateUser  = errors.New("cannot create user")
	ErrCannotGetUser     = errors.New("cannot get user")

	ErrCannotUploadFile = errors.New("cannot upload file")
)
