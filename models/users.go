package models

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Users struct {
	User_id      int64     `gorm:"primaryKey" json:"user_id"`
	User_name    string    `json:"user_name"`
	Code         int       `json:"code"`
	Role         string    `json:"role"`
	Public_Key   string    `json:"public_key"`
	Access_Level int       `json:"access_level"`
	Salt         string    `json:"salt"`
	Passwd       string    `json:"passwd"`
	Last_Login   time.Time `json:"last_login"`
	Passkey      string    `json:"passkey"`
}

func GetUserByUsername(db *gorm.DB, username string) (*Users, error) {
	var user Users
	//db.Raw("select * from tk_users where user_name = ? ", username).Scan(&Users)
	if err := db.Where("user_name = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByUserID(db *gorm.DB, userID string) (*Users, error) {

	user_id, _ := strconv.ParseInt(userID, 10, 64)
	var user Users
	//db.Raw("select * from tk_users where user_name = ? ", username).Scan(&Users)
	if err := db.Where("user_id = ?", user_id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func UserGetByPasskey(db *gorm.DB, passkey string) (*Users, error) {
	var user Users
	if err := db.Where("passkey = ?", passkey).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
