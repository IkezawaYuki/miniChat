package main

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"path/filepath"
)

var ErrNoAvatarURL = errors.New("chat:アバターのURLを取得できません。")

type Avatar interface {
	GetAvatarURL(ChatUser) (string, error)
}

type TryAvatars []Avatar

func (a TryAvatars)GetAvatarURL(u ChatUser)(string, error){
	for _, avatar := range a{
		if url, err := avatar.GetAvatarURL(u); err == nil{
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}

type AuthAvatar struct {}

var UseAuthAvatar AuthAvatar

func (_ AuthAvatar) GetAvatarURL(u ChatUser)(string, error){
	url := u.AvatarURL()
	if url != ""{
		return url, nil
	}
	return "", ErrNoAvatarURL
}

type GravatarAvatar struct {}
var UserGravatar GravatarAvatar

func (_ GravatarAvatar) GetAvatarURL(u ChatUser)(string, error){
	return "//www/gravatar.com/avatar/"+u.UniqueID(), nil
}

type FileSystemAvatar struct {}

var UserFileSystemAvatar FileSystemAvatar

func (_ FileSystemAvatar) GetAvatarURL(u ChatUser)(string, error){
	if files, err := ioutil.ReadDir("avatars"); err == nil{
		for _, file := range files{
			if file.IsDir(){
				continue
			}
			if match, _ := filepath.Match(u.UniqueID() +"*", file.Name()); match{
				return "/avatars/" + file.Name(), nil
			}
		}
	}
	return "", ErrNoAvatarURL
}
























