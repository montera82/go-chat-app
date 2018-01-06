package main

import "testing"

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar

	//var client = new(client)
	//
	//client.userData = map[string]interface{}{
	//	"email": "MyEmailAddress@example.com",
	//}

	var client = &client{
		userData: map[string]interface{}{
			"userId": "d41d8cd98f00b204e9800998ecf8427e",
		},
	}

	url, err := gravatarAvatar.GetAvatarURL(client)

	if err != nil {
		t.Errorf("GetAvatarURL should not return error but got this %s", err)
	}

	if url != "//www.gravatar.com/avatar/d41d8cd98f00b204e9800998ecf8427e" {
		t.Errorf("Wrong url returned %s", url)
	}
}
