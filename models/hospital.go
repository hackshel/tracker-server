package models

import (
	"time"

	"gorm.io/gorm"
)

type Hospital struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Code         int64     `json:"code"`
	HspName      string    `json:"hsp_name"`
	Grade        int32     `json:"grade"`
	ProvinceId   int32     `json:"provinceId"`
	ProvinceName string    `json:"provinceName"`
	CityId       int32     `json:"cityId"`
	CityName     string    `json:"cityName"`
	CountyId     int32     `json:"countyId"`
	CountyName   string    `json:"countyName"`
	AddTime      time.Time `json:"add_time"`
	LastModify   time.Time `json:"last_modify"`
}

func GetHospitalList(db *gorm.DB, page, pageSize int) ([]Hospital, int64, error) {
	var total int64
	var hospitals []Hospital

	offset := (page - 1) * pageSize
	// 获取总记录数
	if err := db.Model(&Hospital{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	// 获取分页数据
	if err := db.Limit(pageSize).Offset(offset).Order("id ASC").Find(&hospitals).Error; err != nil {
		return nil, 0, err
	}
	return hospitals, total, nil
}

func AddHospital(db *gorm.DB, hospital Hospital) (int64, error) {
	rs := db.Create(&hospital)
	if rs.Error != nil {
		return -1, rs.Error
	}
	return hospital.ID, nil
}

func DeleteHospitalByID(db *gorm.DB, id int64) error {
	var hospital Hospital
	r := db.Where("id = ?", id).Delete(&hospital)

	return r.Error
}
