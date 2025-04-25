package models

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

const (
	YesSeeder = "yes"
	NoSeeder  = "no"
	YesConn   = "yes"
	NoConn    = "no"
)

type Peers struct {
	Id             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TorrentID      int64     `json:"torrent_id"`
	PeerID         string    `json:"peer_id"`
	IP             string    `json:"ip"`
	Port           uint16    `json:"port"`
	Uploaded       uint64    `json:"uploaded"`
	Downloaded     uint64    `json:"downloaded"`
	LeftSize       uint64    `json:"left_size"`
	IsSeeder       string    `json:"is_seeder"`
	Started        time.Time `json:"started"`
	LastAction     time.Time `json:"last_action"`
	PrevAction     string    `json:"prev_action"`
	Connectable    string    `json:"connectable"`
	Agent          string    `json:"agent"`
	FinishedAt     time.Time `json:"finished_at"`
	Downloadoffset uint64    `json:"downloadoffset"`
	Uploadoffset   uint64    `json:"uploadoffset"`
	DLSpeed        string    `json:"dl_speed"`
	Passkey        string    `json:"passkey"`
	InfoHash       string    `json:"InfoHash"`
}

func GetPeers(db *gorm.DB, torrent_id int64, is_seeder string, peers_count uint32, numwant uint32, info_hash string) ([]Peers, error) {

	var peers []Peers
	var peer_limit string

	var only_leech string = ""

	if peers_count > numwant {
		peer_limit = fmt.Sprintf(" ORDER BY RAND() LIMIT %d", numwant)
	}

	// 判断请求者是否是种子者
	if is_seeder == "yes" {
		only_leech = " and is_seeder = 'no'"
	}

	dt := time.Now().Unix()

	fields := fmt.Sprintf("id, is_seeder, peer_id, passkey, ip, port, uploaded, downloaded, (%d - UNIX_TIMESTAMP(last_action)) AS announce_time, UNIX_TIMESTAMP(prev_action) AS prevts", dt)

	strSQL := fmt.Sprintf("select %s from tk_peers where torrent_id = %d and info_hash = '%s' and connectable = 'yes' %s %s", fields, torrent_id, info_hash, only_leech, peer_limit)

	log.Printf("GetPeers sql is : %v\n", strSQL)
	res := db.Raw(strSQL).Scan(&peers)

	if res.Error != nil {
		log.Println("sql exec error : ", res.Error)
		return nil, res.Error
	}
	return peers, nil

}

func GetPeerByPasskey(db *gorm.DB, torrent_id int64, passkey string) ([]Peers, error) {
	t := time.Now().Unix()
	var peers []Peers

	fields := fmt.Sprintf("id, is_seeder, peer_id, ip, port, uploaded, downloaded, (%d - UNIX_TIMESTAMP(last_action)) AS announce_time, UNIX_TIMESTAMP(prev_action) AS prev_action", t)

	strSQL := fmt.Sprintf(" SELECT %s FROM tk_peers WHERE torrent_id = %d AND passkey = '%s'", fields, torrent_id, passkey)

	fmt.Printf("debug SQL %v,  %T\n", strSQL, strSQL)
	res := db.Raw(strSQL).Scan(&peers)

	if res.Error != nil {
		log.Println("sql exec error : ", res.Error)
		return nil, res.Error
	}
	return peers, nil

}

func AddNewPeer(db *gorm.DB, torrent int64, peer_id string, client_ip string, port uint16, uploaded uint64, downloaded uint64, left uint64, is_seeder string, connectable string, passkey string, info_hash string) (int64, error) {

	dt := time.Now()
	prev := dt.Format("2006-01-02 15:04:05")
	var finished_at_time time.Time

	fmt.Printf("debug === > is_seeder : %v , type %T\n", is_seeder, is_seeder)
	if is_seeder == "yes" {
		finished_at_time = dt
	} else {
		finished_at_time = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	peer := &Peers{
		TorrentID:      torrent,
		PeerID:         peer_id,
		IP:             client_ip,
		Port:           port,
		Uploaded:       uploaded,
		Downloaded:     downloaded,
		LeftSize:       left,
		IsSeeder:       is_seeder,
		Started:        dt,
		LastAction:     dt,
		PrevAction:     prev,
		Connectable:    connectable,
		Agent:          "",
		FinishedAt:     finished_at_time,
		Downloadoffset: downloaded,
		Uploadoffset:   uploaded,
		DLSpeed:        "",
		Passkey:        passkey,
		InfoHash:       info_hash,
	}
	result := db.Create(&peer)

	if result.Error != nil {
		return 0, result.Error
	} else {
		return peer.Id, nil
	}
}

func UpdatePeer(db *gorm.DB, torernt_id int64, passkey string,
	client_ip string, port uint16, uploaded uint64, downloaded uint64,
	left uint64, is_seeder string) error {

	dt := time.Now()

	// updateSQL = fmt.Sprintf("UPDATE tk_peers  SET"+
	// 	" ip = '%s' "+
	// 	" port = %d "+
	// 	" uploaded = %d "+
	// 	" downloaded = %d"+
	// 	" left_size = %d "+
	// 	" prev_action =  '%s' "+
	// 	" last_action = '%s' "+
	// 	" is_seeder = '%s'  "+
	// 	" agent = ''  "+
	// 	" finished_at = '%s'  "+
	// 	" WHERE torrent_id = %d  and passkey  = '%s'  ",
	// 	client_ip, port, uploaded, downloaded, left, dt, dt, is_seeder, dt, torernt_id, passkey,
	// )

	rs := db.Model(&Peers{}).
		Where("torrent_id = ? AND passkey = ?", torernt_id, passkey).
		Updates(map[string]interface{}{
			"ip":          client_ip,
			"port":        port,
			"uploaded":    uploaded,
			"downloaded":  downloaded,
			"left_size":   left,
			"prev_action": dt,
			"last_action": dt,
			"is_seeder":   is_seeder,
			"agent":       "",
			"finished_at": dt,
		})

	if rs.Error != nil {
		return rs.Error
	}

	return nil

}

func RemovePeerBy(db *gorm.DB, torrent_id int64, peer_id string, passkey string) error {

	err := db.Where("torrent_id = ? and passkey = ?", torrent_id, passkey).Delete(&Peers{}).Error

	return err

}
