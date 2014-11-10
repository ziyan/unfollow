// This package contains all models of Unfollow
package models

import (
    "encoding/gob"
)

func init() {
    gob.Register(Session{})
    gob.Register(User{})
    gob.Register(Username{})
}
