package controllers

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hackshel/tracker-server/pkg/bencode"
	"gorm.io/gorm"
)

type AnnounceType string

const (
	STARTED   AnnounceType = "started"
	STOPPED   AnnounceType = "stopped"
	COMPLETED AnnounceType = "completed"
	PAUSED    AnnounceType = "paused"
	ANNOUNCE  AnnounceType = ""
)

// func AnnounceType 返回 AnnounceType string
func ParseAnnounceType(t string) AnnounceType {
	switch t {
	case "started":
		return STARTED
	case "stopped":
		return STOPPED
	case "completed":
		return COMPLETED
	case "paused":
		return PAUSED
	default:
		return ANNOUNCE
	}
}

type CryptoLevel uint

const (
	Unencrypted CryptoLevel = 0
	Supported   CryptoLevel = 1
	Required    CryptoLevel = 2
)

type announceRequest struct {
	Compact     bool // Force compact always?
	Downloaded  uint32
	Left        uint32
	Uploaded    uint32
	Corrupt     uint32
	Event       AnnounceType
	IP          net.IP
	IPv6        bool
	InfoHash    string
	NumWant     uint
	Passkey     string
	PeerID      string
	Port        uint16
	TrackerID   string
	Key         string
	CryptoLevel CryptoLevel
}
type InfoHash [20]byte

func Announce(c *gin.Context) {
	// start := time.Now()
	c.Header("Content-Type", "text/plain")
	c.Header("charset", "utf8")
	c.Header("Cache-Control", "no-cache")

	// // 1. 获取请求参数
	// port := c.Query("port")
	// downloaded := c.Query("downloaded")
	// uploaded := c.Query("uploaded")
	// left := c.Query("left")
	// compact := c.Query("compact")
	// numwant := c.Query("numwant")
	// infoHash := c.Query("info_hash")
	// peerID := c.Query("peer_id")
	// event := c.Query("event")
	// ip := c.ClientIP()

	// // 参数校验
	// if infoHash == "" || peerID == "" || port == "" {
	// 	bencode.Marshal(c.Writer, map[string]string{
	// 		"failure reason": "Missing required parameters",
	// 	})
	// 	return
	// }

}

type TorrentStat struct {
	InfoHash   string `gorm:"column:info_hash"`
	Completed  int    `gorm:"column:times_completed"`
	Downloaded int    `gorm:"column:downloaded"`
	Seeders    int    `gorm:"column:seeders"`  // complete
	Leechers   int    `gorm:"column:leechers"` // incomplete
}

func Scrape(c *gin.Context) {

	// 获取所有 info_hash 参数
	infoHashes := c.Request.URL.Query()["info_hash"]
	if len(infoHashes) == 0 {
		c.String(http.StatusBadRequest, "Missing info_hash")
		return
	}

	// 初始化返回结构
	response := map[string]interface{}{
		"files": make(map[string]map[string]int),
	}

	db := c.MustGet("db").(*gorm.DB)

	for _, hash := range infoHashes {

		hashStr := string(hash)

		var stat TorrentStat
		if err := db.Where("info_hash = ?", hashStr).First(&stat).Error; err != nil {
			// 不存在的 info_hash 也需要返回空结构体
			(response["files"].(map[string]map[string]int))[hashStr] = map[string]int{
				"complete":   0,
				"incomplete": 0,
				"downloaded": 0,
			}
			continue
		}

		// 将结果填入 response
		(response["files"].(map[string]map[string]int))[hashStr] = map[string]int{
			"complete":   stat.Seeders,
			"incomplete": stat.Leechers,
			"downloaded": stat.Downloaded,
		}
	}

	c.Header("Content-Type", "text/plain")
	if err := bencode.Marshal(c.Writer, response); err != nil {
		c.String(http.StatusInternalServerError, "Failed to encode scrape response")
	}

}
