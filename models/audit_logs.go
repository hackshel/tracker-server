package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

/*
CREATE TABLE `tk_audit_logs` (
  `log_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT ' ID',
  `user_id` bigint unsigned NOT NULL COMMENT ' ID',
  `operation_type` enum('UPLOAD','DOWNLOAD','DELETE','MODIFY','SHARE','LOGIN') NOT NULL,
  `target_torrent` bigint unsigned DEFAULT NULL COMMENT ' ID',
  `detail` varchar(500) NOT NULL,
  `timestamp` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'UTC ',
  `blockchain_hash` varchar(66) DEFAULT NULL,
  PRIMARY KEY (`log_id`),
  KEY `target_torrent` (`target_torrent`),
  CONSTRAINT `tk_audit_logs_ibfk_1` FOREIGN KEY (`target_torrent`) REFERENCES `tk_torrents` (`torrent_id`) ON DELETE SET NULL
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

为保存完整日志删除外键
 alter table tk_audit_logs drop FOREIGN KEY  tk_audit_logs_ibfk_1;
*/

type AuditLogs struct {
	LogId          int64     `gorm:"primaryKey;autoIncrement" json:"log_id"`
	UserID         int64     `json:"user_id"`
	OperationType  string    `json:"operation_type"`
	TargetTorrent  int64     `json:"target_torrent"`
	Detail         string    `json:"detail"`
	Timestamp      time.Time `json:"timestamp"`
	BlockchainHash string    `json:"blockchain_hash"`
}

func (AuditLogs) TableName() string {
	return "tk_audit_logs"
}

func AddLogs(db *gorm.DB, log AuditLogs) error {

	fmt.Printf("add log data  : %v, type %T\n", log, log)
	fmt.Printf("log target_torrent %v\n", log.TargetTorrent)
	err := db.Create(&log).Error
	fmt.Printf("err %v \n", err)
	return err
}
