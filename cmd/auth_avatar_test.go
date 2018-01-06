package main

import "testing"

func TestGetAvatarURL(t *testing.T) {

	var authAvatar AuthAvatar
	c := new(client)

	_, err := authAvatar.GetAvatarURL(c)

	if err != ErrNoAvatarURL {
		t.Error("error returned by AuthAvater when no" +
			"url is present in client should be ErrNoAvatarURL ")
	}

	c.userData = map[string]interface{}{"avatar_url": "foo.png"}

	url, err := authAvatar.GetAvatarURL(c)

	if url != "foo.png" {
		t.Errorf("expected foo.png, but got '%s'", url)
	}

	if err != nil {
		t.Errorf("GetAvatar should work as expected when url is present"+
			"but got error %s", err)
	}
}
