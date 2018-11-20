package main

import (
	"errors"
	"html/template"
	"io"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	"github.com/ykpythemind/funho/model"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	temp := template.Must(template.Must(t.templates.Lookup("layout").Clone()).AddParseTree("content", t.templates.Lookup(name).Tree))
	return temp.ExecuteTemplate(w, "layout", data)
}

type CustomContext struct {
	echo.Context
	DB        *gorm.DB
	LoginUser *model.User
}

// GetUserBySession gets user from session
func (c *CustomContext) GetUserFromSession() error {
	sess, _ := session.Get("session", c)
	n := sess.Values["login_user"]

	name, ok := n.(string)

	if !ok {
		return errors.New("no")
	}

	if name == "" {
		return errors.New("name is blank")
	}

	user := &model.User{}
	if err := c.DB.Where(model.User{Name: name}).First(&user).Error; gorm.IsRecordNotFoundError(err) {
		return err
	}

	c.LoginUser = user

	return nil
}

func (c *CustomContext) DestroySession() error {
	sess, _ := session.Get("session", c)

	sess.Values["login_user"] = nil
	sess.Save(c.Request(), c.Response())
	return nil
}

func (c *CustomContext) LoginAndStoreSession(username, password string) error {
	user := &model.User{}
	// FIXME
	// 手抜きログイン
	if err := c.DB.Where(model.User{Name: username, Password: password}).First(&user).Error; gorm.IsRecordNotFoundError(err) {
		return c.Redirect(http.StatusMovedPermanently, pathForApp["login_form_path"].toString())
	}

	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}

	// FIXME
	sess.Values["login_user"] = user.Name
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusMovedPermanently, pathForApp["chat_path"].toString())
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(*CustomContext)

		if err := cc.GetUserFromSession(); err != nil {
			return c.Redirect(http.StatusMovedPermanently, pathForApp["login_form_path"].toString())
		}

		return next(c)
	}
}

func main() {
	db, err := gorm.Open("sqlite3", "development.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	e := NewEcho(db)
	e.Logger.Fatal(e.Start("localhost:1323"))
}

func NewEcho(db *gorm.DB) *echo.Echo {
	e := echo.New()

	// This middleware should be registered before any other middleware.
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{c, db, nil}
			return h(cc)
		}
	})

	// Middleware
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	e.GET("/", func(c echo.Context) error {
		cc := c.(*CustomContext)
		cc.GetUserFromSession()
		user := cc.LoginUser
		if user != nil {
			return c.Redirect(http.StatusMovedPermanently, pathForApp["chat_path"].toString())
		}
		return c.String(http.StatusOK, "Hello, World!")
	})

	h := Handler{DB: db}

	// Session
	// e.GET("/session", h.GetSessionForm, AuthMiddleware)
	e.GET("/login", h.GetSessionForm)
	e.POST("/login", h.CreateSession)
	e.POST("/logout", h.DestroySession)

	// Chat API
	m := e.Group("chat", AuthMiddleware)
	m.GET("", h.GetChatIndex)
	m.GET("/:room_id", h.GetChats)
	m.POST("/:room_id", h.PostChat)

	// Static files
	e.Static("/static", "assets")

	e.Renderer = setupTemplates()

	return e
}

func setupTemplates() *Template {
	return &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
}