package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/ykpythemind/funho/model"
)

func main() {
	db, err := gorm.Open("sqlite3", "development.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	createUsersRecord(db)
}

func createUsersRecord(db *gorm.DB) {

	users := []string{
		"test1",
		"test2",
		"test3",
	}

	for _, userName := range users {
		user := model.User{Name: userName, Password: "password"}
		db.Create(&user)
	}

}
