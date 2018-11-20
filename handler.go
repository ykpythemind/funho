package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/ykpythemind/funho/model"
)

type Handler struct {
	DB *gorm.DB
}

func (h *Handler) GetChatIndex(c echo.Context) error {
	cc := c.(*CustomContext)

	user := cc.LoginUser
	if user == nil {
		return errors.New("no")
	}

	p := NewPresenter(h, user, newChats(user))

	return c.Render(http.StatusOK, "chat", &p)
}

func (h *Handler) GetChats(c echo.Context) error {
	id := c.Param("room_id")
	var chats []model.Chat
	h.DB.Preload("User").Order("created_at desc").Where("room_user_id = ?", id).Find(&chats)

	return c.JSON(http.StatusOK, chats)
}

func (h *Handler) PostChat(c echo.Context) error {
	cc := c.(*CustomContext)

	loginUser := cc.LoginUser
	if loginUser == nil {
		return errors.New("no")
	}

	roomId := c.Param("room_id")

	newChat := new(model.Chat)
	if err := c.Bind(newChat); err != nil {
		return err
	}
	newChat.UserID = loginUser.ID
	newChat.RoomUserID, _ = strconv.Atoi(roomId)

	if newChat.Body == "" {
		return c.String(http.StatusUnprocessableEntity, "body is blank")
	}

	h.DB.Create(newChat)

	return c.String(http.StatusOK, "")
}

func (h *Handler) GetSessionForm(c echo.Context) error {
	cc := c.(*CustomContext)

	user := cc.LoginUser
	if user == nil {
		return c.Render(http.StatusOK, "login", "you must login")
	}
	return c.Redirect(http.StatusMovedPermanently, pathForApp["chat_path"].toString())
}

func (h *Handler) CreateSession(c echo.Context) error {
	cc := c.(*CustomContext)
	username := c.FormValue("username")
	password := c.FormValue("password")

	return cc.LoginAndStoreSession(username, password)
}

func (h *Handler) DestroySession(c echo.Context) error {
	cc := c.(*CustomContext)
	cc.DestroySession()
	return cc.Redirect(http.StatusMovedPermanently, pathForApp["login_form_path"].toString())
}
