package main

import (
	"errors"
	"os"
	"path"
)

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
	GetAvatarURL(ChatUser) (string, error)
}

type AuthAvatar struct {
}

var UseAuthAvatar AuthAvatar

// get away from  (auth AuthAvatar) =>(AuthAvatar) no explicit auth variable to avoid nil reference failure
func (AuthAvatar) GetAvatarURL(c ChatUser) (string, error) {
	url := c.AvatarURL()
	if len(url) == 0 {
		return "", ErrNoAvatarURL
	}
	return url, nil
}

type GravatarAvatar struct {
}

var UseGravatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(c ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + c.UniqueID(), nil
}

type FileSystemAvatar struct {
}

var UseFileSystemAvatar FileSystemAvatar

func (FileSystemAvatar) GetAvatarURL(c ChatUser) (string, error) {
	if files, err := os.ReadDir("avatars"); err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if match, _ := path.Match(c.UniqueID()+"*", file.Name()); match {
				return "/avatars/" + file.Name(), nil
			}
		}
	}
	return "", ErrNoAvatarURL
}
