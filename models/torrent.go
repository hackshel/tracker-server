package models

import (
	"time"

	"gorm.io/gorm"
)

type Torrents struct {
	Torrent_id      int64     `gorm:"primaryKey;autoIncrement" json:"torrent_id"`
	Info_hash       string    `json:"info_hash"`
	Name            string    `json:"name"`
	Filename        string    `json:"filename"`
	Dcp_uuid        string    `json:"dcp_uuid"`
	Dcp_size        int64     `json:"dcp_size"`
	Piece_length    int       `json:"piece_length"`
	Pices_count     int       `json:"pices_count"`
	Added_time      time.Time `json:"added_time"`
	Dcp_type        string    `json:"dcp_type"`
	Numfiles        int       `json:"numfiles"`
	Tracker_url     string    `json:"tracker_url"`
	F_sha1          string    `json:"f_sha1"`
	Seeders         int       `json:"seeders"`
	Leechers        int       `json:"leechers"`
	Times_completed int       `json:"times_completed"`
	Last_action     time.Time `json:"last_action"`
	User_id         int64     `json:"user_id"`
}

func (t *Torrents) TableName() string {
	return "tk_torrents"
}

func GetTorrentList(db *gorm.DB, page, pageSize int, user_id int64) ([]Torrents, int64, error) {
	var torrents []Torrents
	var total int64

	offset := (page - 1) * pageSize

	if user_id == -1 {

		// 获取总记录数
		if err := db.Model(&Torrents{}).Count(&total).Error; err != nil {
			return nil, 0, err
		}
		// 获取分页数据
		if err := db.Limit(pageSize).Offset(offset).Order("added_time DESC").Find(&torrents).Error; err != nil {
			return nil, 0, err
		}

	} else {
		// 获取总记录数
		if err := db.Model(&Torrents{}).Where("user_id = ?", user_id).Count(&total).Error; err != nil {
			return nil, 0, err
		}
		// 获取分页数据
		if err := db.Where("user_id = ?", user_id).Limit(pageSize).Offset(offset).Order("added_time DESC").Find(&torrents).Error; err != nil {
			return nil, 0, err
		}
	}
	return torrents, total, nil
}

func GetTorrentByIDAndHash(db *gorm.DB, torrentID int64, infoHash string) (*Torrents, error) {
	var torrent Torrents
	err := db.Where("torrent_id = ? AND info_hash = ?", torrentID, infoHash).First(&torrent).Error
	if err != nil {
		return nil, err
	}
	return &torrent, nil
}

func GetTorentsCountByUserID(db *gorm.DB, user_id int64) (int64, error) {
	var cnt int64
	// err := db.Debug().Where("user_id = ? ", user_id).Count(&cnt).Error
	err := db.Model(&Torrents{}).Where("user_id = ? ", user_id).Count(&cnt).Error
	if err != nil {
		return -1, err
	}

	return cnt, nil
}
