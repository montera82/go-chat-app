package main

import (
	"fmt"
)

type GravatarAvatar struct{}

var useGravatarAvatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userId, ok := c.userData["userId"]; ok {
		if userIdStr, ok := userId.(string); ok {
			return fmt.Sprintf("//www.gravatar.com/avatar/%s", userIdStr), nil
		}
	}

	return "", ErrNoAvatarURL
}
