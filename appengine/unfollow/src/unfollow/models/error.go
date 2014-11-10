package models

import (
    "errors"
)

var (
    ErrUsernameAlreadyTaken = errors.New("models: username already taken")
)
