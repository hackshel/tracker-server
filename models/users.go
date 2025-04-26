package models

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Users struct {
	User_id      int64     `gorm:"primaryKey;autoIncrement" json:"user_id"`
	User_name    string    `json:"user_name"`
	Code         int64     `json:"code"`
	Role         string    `json:"role"`
	Public_Key   string    `json:"public_key"`
	Access_Level int       `json:"access_level"`
	Salt         string    `json:"salt"`
	Passwd       string    `json:"passwd"`
	Last_Login   time.Time `json:"last_login"`
	Passkey      string    `json:"passkey"`
}

type UserWithHospital struct {
	UserID       int64     `json:"user_id"`
	UserName     string    `json:"user_name"`
	Code         int64     `json:"code"`
	Role         string    `json:"role"`
	PublicKey    string    `json:"public_key"`
	AccessLevel  int       `json:"access_level"`
	Salt         string    `json:"salt"`
	Passwd       string    `json:"passwd"`
	LastLogin    time.Time `json:"last_login"`
	Passkey      string    `json:"passkey"`
	HospitalName string    `json:"hospital_name"` // 新增字段：医院名称
	HospitalCode string    `json:"hospital_code"` // 新增字段：医院编码
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

func GetUserList(db *gorm.DB, page int, pageSize int) ([]Users, int64, error) {
	var users []Users
	var total int64

	offset := (page - 1) * pageSize
	// 获取总记录数
	if err := db.Model(&Users{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	// 获取分页数据
	if err := db.Limit(pageSize).Offset(offset).Order("user_id ASC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil

}

func GetUserListWithHospital(db *gorm.DB, page int, pageSize int) ([]UserWithHospital, int64, error) {
	var users []UserWithHospital
	var total int64

	offset := (page - 1) * pageSize

	// 查询总数（还是只算user表的总数）
	if err := db.Model(&Users{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 联表查询 users 和 hospitals
	err := db.Table("tk_users").
		Select("tk_users.*, tk_hospital.hsp_name as hospital_name, tk_hospital.code as hospital_code").
		Joins("left join tk_hospital on tk_users.code = tk_hospital.code").
		Order("tk_users.user_id ASC").
		Limit(pageSize).
		Offset(offset).
		Scan(&users).Error

	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func DeleteUserByID(db *gorm.DB, user_id int64) error {

	var user Users
	r := db.Where("user_id = ?", user_id).Delete(&user)

	return r.Error

}

func AddNewUser(db *gorm.DB, user Users) (int64, error) {
	rs := db.Create(&user)
	if rs.Error != nil {
		return -1, rs.Error
	}
	return user.User_id, nil
}

func UpdateUserHospitalCode(db *gorm.DB, user_id int64, code int64) error {
	rs := db.Model(&Users{}).Where("user_id = ?", user_id).Update("code", code)
	return rs.Error
}

func UpdateUserFields(db *gorm.DB, user_id int64, fields map[string]interface{}) error {
	rs := db.Model(&Users{}).Where("user_id = ?", user_id).Updates(fields)
	return rs.Error
}
