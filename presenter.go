package main

import (
	"fmt"

	"github.com/ykpythemind/funho/model"
)

type Presenter struct {
	LoginUser model.User
	Users     []model.User
	Chats     []model.Chat
	Path      map[string]Path
}

type Path string

var pathForApp = map[string]Path{
	"root_path":        Path("/"),
	"new_chat_path":    Path("/chat"),
	"chat_path":        Path("/chat"),
	"login_form_path":  Path("/login"),
	"new_session_path": Path("/login"),
	"logout_path":      Path("/logout"),
}

func (p Path) toString() string {
	return string(p)
}

func NewPresenter(handler *Handler, loginUser *model.User, chats []model.Chat) Presenter {
	var users []model.User
	handler.DB.Find(&users)

	fmt.Printf("users[1].ID = %+v\n", users[1].ID)
	fmt.Printf("users[1].ID = %+v\n", users[1].Password)

	p := Presenter{
		LoginUser: *loginUser,
		Chats:     chats,
		Users:     users,
		Path:      pathForApp,
	}
	return p
}

// return sample data
func newChats(user *model.User) []model.Chat {

	s := []string{"ya", "hey", "ho!"}

	chats := make([]model.Chat, len(s))

	for i, body := range s {
		var u *model.User
		if i == 1 {
			u = user
		} else {
			u = &model.User{Name: "hoge"}
		}
		chats[i] = model.Chat{
			Body: body,
			User: *u,
		}
	}

	return chats
}
