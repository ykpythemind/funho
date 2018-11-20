package main

import (
	"github.com/ykpythemind/funho/model"
)

type Presenter struct {
	LoginUser model.User
	Users     []model.User
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

func NewPresenter(handler *Handler, loginUser *model.User) Presenter {
	var users []model.User
	handler.DB.Find(&users)

	p := Presenter{
		LoginUser: *loginUser,
		Users:     users,
		Path:      pathForApp,
	}
	return p
}
