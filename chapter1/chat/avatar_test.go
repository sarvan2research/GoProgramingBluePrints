package main

import (
	gomniauthtest "github.com/stretchr/gomniauth/test"
	"os"
	"path/filepath"
	"testing"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	testUser := &gomniauthtest.TestUser{}
	testUser.On("AvatarURL").Return("", ErrNoAvatarURL)
	testChatUser := &chatUser{
		User: testUser,
	}
	url, err := authAvatar.GetAvatarURL(testChatUser)
	if err != ErrNoAvatarURL {
		t.Error("AuthAvatar.GetAvatar should return ErrNoAvatarURL when no value present")
	}
	// set a value
	testUrl := "http://urlto-gravatar"
	testUser = &gomniauthtest.TestUser{}
	testUser.On("AvatarURL").Return(testUrl, nil)
	testChatUser = &chatUser{
		User: testUser,
	}
	url, err = authAvatar.GetAvatarURL(testChatUser)
	if err != nil {
		t.Error("AuthAvatar.GetAvatar should return ErrNoAvatarURL when no value present")
	}
	if url != testUrl {
		t.Error("AuthAvatar.GetAvatar should return correct URL")
	}
}

/*
Gravatar hash creation principle
Creating Hash from gravatar
Trim leading and trailing whitespace from an email address
Force all characters to lower-case
md5 hash the final string
https://en.gravatar.com/site/implement/hash/
*/
func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	user := &chatUser{uniqueID: "abc"}
	url, err := gravatarAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURL should not return an error")
	}
	if url != "//www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346" {
		t.Errorf("GravatarAvatar.GetAvatarURL wrongly returned %s", url)
	}
}

func TestFileSystemAvatar(t *testing.T) {
	fileName := filepath.Join("avatars", "123.jpg")
	os.WriteFile(fileName, []byte{}, 0777)
	defer os.Remove(fileName)
	var fileSystemAvatar FileSystemAvatar
	user := &chatUser{uniqueID: "abc"}
	url, err := fileSystemAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("FileSystemAvatar.GetAvatarURL should not return an error")
	}
	if url != "/avatars/abc.jpg" {
		t.Errorf("FileSystemAvatar.GetAvatarURL wrongly returned %s", url)
	}
}
