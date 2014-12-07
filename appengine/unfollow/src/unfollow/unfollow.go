// This is the main package for unfollow.
// It will import the sub-packages automatically.
// For example, in app.go we have:
//
//     package main
//
//     import _ "unfollow"
//
//     func init() {
//     }
//
package unfollow

import (
    _ "unfollow/network"
    _ "unfollow/page"
    _ "unfollow/task"
    _ "unfollow/urls"
    _ "unfollow/user"
    _ "unfollow/web"
)

func init() {
}
