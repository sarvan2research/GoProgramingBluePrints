package main

import "errors"

// ErrNoAvatar is the error that is returned when the
// Avatar instance is unable to provide an avatar URL.

var ErrNoAvatarURL = errors.New("chat:Unable to get an avatar URL.")

// Avatar represents types capable of representing
// user profile pictures.
type Avatar interface {
	// GetAvatarURL gets the avatar URL for the specified client,
	// or returns an error if something goes wrong.
	// ErrNoAvatarURL is returned if the object is unable to get
	// a URL for the specified client.
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct {
}

// get away from  (auth AuthAvatar) =>(AuthAvatar) no explicit auth variable to avoid nil reference failure
func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlString, ok := url.(string); ok {
			return urlString, nil
		}
	}
	return "", ErrNoAvatarURL
}
