package main

import "testing"

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	client := new(client)
	url, err := authAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("AuthAvatar.GetAvatar should return ErrNoAvatarURL when no value present")
	}
	// set a value
	testUrl := "http://urlto-gravatar"
	client.userData = map[string]interface{}{"avatar_url": testUrl}
	url, err = authAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("AuthAvatar.GetAvatar should return ErrNoAvatarURL when no value present")
	}
	if url != testUrl {
		t.Error("AuthAvatar.GetAvatar should return correct URL")
	}
}
