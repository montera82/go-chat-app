package main

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	url, ok := c.userData["avatar_url"]

	if !ok {
		return "", ErrNoAvatarURL
	}

	return url.(string), nil
}
