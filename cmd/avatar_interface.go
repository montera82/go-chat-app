package main

import "errors"

var ErrNoAvatarURL = errors.New("unable to find Avatar url")

type Avatar interface {
	GetAvatarURL(c *client) (string, error)
}
