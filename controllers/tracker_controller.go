package controllers

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hackshel/tracker-server/models"
	"github.com/hackshel/tracker-server/pkg/bencode"
	"github.com/hackshel/tracker-server/pkg/errs"
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

// "/api/v1/tracker/announce?
// passkey=2be1adfa1f3911f0aeeec3ae841ca01b
// &info_hash=Jp%1f%a0%2b%95)%7d%a8%b9%09%c1%13%19N%cf%c1_%d1%eb
// &peer_id=-qB5040-z5qdnog74hT4
// &port=25114
// &uploaded=0
// &downloaded=0
// &left=1732048795
// &corrupt=0
// &key=F2746450
// &event=started
// &numwant=200
// &compact=1
// &no_peer_id=1
// &supportcrypto=1
// &redundant=0"

type announceRequest struct {
	Passkey     string
	InfoHash    string
	PeerID      string
	Port        uint16
	Uploaded    uint64
	Downloaded  uint64
	Left        uint64
	Corrupt     uint32
	Key         string
	Event       AnnounceType
	NumWant     uint32
	Compact     bool // Force compact always?
	NoPeerID    bool // Force no_peer_id always?
	CryptoLevel CryptoLevel
	IP          string //net.IP
	TrackerID   string
}

// type InfoHash [20]byte

func responseBencodeError(message string) []byte {
	var buf bytes.Buffer
	response := map[string]string{"failure reason": message}
	if err := bencode.Marshal(&buf, response); err != nil {
		log.Printf("Failed to encode error response: %s", err)
	}
	return buf.Bytes()
}

func sendErrorResponse(c *gin.Context, error_code int) {
	error_msg := errs.GetMsg(error_code)
	log.Printf("func sendErrorResponse :Error in request from: %s (%d : %s)\n", c.Request.RequestURI, error_code, error_msg)
	c.Data(http.StatusOK, gin.MIMEPlain, responseBencodeError(error_msg))
}

func portblacklisted(port uint) bool {
	//direct connect
	if port >= 411 && port <= 413 {
		return true
	}
	//bittorrent
	if port >= 6881 && port <= 6889 {
		return true
	}
	//kazaa
	if port == 1214 {
		return true
	}
	//gnutella
	if port >= 6346 && port <= 6347 {
		return true
	}
	//emule
	if port == 4662 {
		return true
	}
	//winmx
	if port == 6699 {
		return true
	}
	return false
}

func AnnounceParams(c *gin.Context) (*announceRequest, int) {
	infoHashStr := c.Query("info_hash")
	if infoHashStr == "" {
		return nil, errs.MSG_INVALID_INFOHASH
	}
	//fmt.Printf("infohashSTR : -- > %v , type %T, length : %v \n", infoHashStr, infoHashStr, len(infoHashStr))

	var infoHash [20]byte
	if len(infoHashStr) != 20 {
		return nil, errs.MSG_INVALID_INFOHASH
	}
	copy(infoHash[:], infoHashStr)

	hexInfohashStr := fmt.Sprintf("%x", infoHash[:])
	//fmt.Printf("xr value %v, type %T\n", xr, xr)
	//fmt.Printf("info hash %v, type %T, length %v , convert string : %v \n", infoHash, infoHash, len(infoHash), string(infoHash[:]))

	//fmt.Printf("info hash convert to string %v, type %T, length : %v\n", hexStr, hexStr, len(hexStr))

	peerID := c.Query("peer_id")
	if peerID == "" || len(peerID) != 20 {
		return nil, errs.MSG_MISSING_PEERID
	}

	peerid_decode, _ := url.QueryUnescape(peerID)
	hexPeerIdStr := hex.EncodeToString([]byte(peerid_decode))
	//fmt.Printf("peer id %v , %v, type %T, length : %v\n", hexPeerIdStr, peerid_decode, hexPeerIdStr, len(hexPeerIdStr))

	ClientIP := c.ClientIP()
	if ClientIP == "" {
		log.Printf("Failed to parse client ip: %s\n", c.Request.RemoteAddr)
		return nil, errs.MSG_MALFORMED_REQUEST
	}
	// fmt.Printf("client_ip is %v , %T\n", ClientIP, ClientIP)

	portstr := c.Query("port")
	port64, _ := strconv.ParseUint(portstr, 10, 16)
	port := uint16(port64)

	if portblacklisted(uint(port)) {
		return nil, errs.MSG_PORT_BANED
	}

	if port < 1024 {
		return nil, errs.MSG_INVALID_PORT
	}

	cryptoLevel := Unencrypted
	sup := c.Query("supportcrypto")

	if sup != "" {
		cryptoLevel = Supported
	}
	req_s := c.Query("requirecrypto")
	if req_s != "" {
		cryptoLevel = Required
	}

	downloaded_str := c.Query("downloaded")
	if downloaded_str == "" {
		return nil, errs.MSG_MISSING_DL
	}
	downloaded, _ := strconv.ParseUint(downloaded_str, 10, 32)

	uploaded_str := c.Query("uploaded")
	if uploaded_str == "" {
		return nil, errs.MSG_MISSING_UL
	}
	uploaded, _ := strconv.ParseUint(uploaded_str, 10, 32)

	left_str := c.Query("left")
	if left_str == "" {
		return nil, errs.MSG_MISSING_LEFT
	}
	left, _ := strconv.ParseUint(left_str, 10, 32)

	numwant_str := c.Query("numwant")
	if numwant_str == "" {
		return nil, errs.MSG_INVALID_NUMWANT
	}
	numwant64, _ := strconv.ParseUint(numwant_str, 10, 32)
	numwant := uint32(numwant64)

	key := c.Query("key")

	return &announceRequest{
		Compact:     true,
		Downloaded:  downloaded,
		Uploaded:    uploaded,
		Event:       ParseAnnounceType(c.Query("event")),
		IP:          ClientIP,
		InfoHash:    hexInfohashStr,
		Left:        left,
		NumWant:     numwant,
		PeerID:      hexPeerIdStr,
		NoPeerID:    true,
		Port:        port,
		Key:         key,
		CryptoLevel: cryptoLevel,
	}, errs.SUCCESS
}

type Response struct {
	Interval    int64
	MinInterval int64
	Complete    uint32
	Incomplete  uint32
	Peers       interface{}
	// Peers       []string

}

type RespPeer struct {
	PeerID string `json:"optional_field,omitempty"`
	IP     string //net.IP
	Port   uint16
}

type RespPeerNoID struct {
	IP   string //net.IP
	Port uint16
}

var RespPeerCompact []string

func Announce(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	passkey := c.Query("passkey")
	if passkey == "" {
		sendErrorResponse(c, errs.MSG_INVALID_AUTH)
		return
	}

	foundUser, user_err := models.UserGetByPasskey(db, passkey)

	if foundUser == nil || user_err != nil {
		log.Printf("passkey (%v) is not in the db, passed ... \n", passkey)
		sendErrorResponse(c, errs.MSG_INVALID_AUTH)
		return
	}
	log.Printf("user info %v ,type %T\n", foundUser, foundUser)

	req, code := AnnounceParams(c)
	if code != errs.SUCCESS {
		sendErrorResponse(c, code)
		return
	}

	// 判断是否是seeder
	is_seeder := "no"
	if req.Left == 0 {
		is_seeder = "yes"
	}

	// 通过info_hash 获取torrent
	torrent, tor_err := models.GetTorrentByInfoHash(db, req.InfoHash)
	if tor_err != nil {
		fmt.Printf("torrent by hash err : %v\n", tor_err)
	}
	//fmt.Printf("torrent val %v, type %T\n", torrent, torrent)

	// 共同下载用户总量计算
	var peers_count uint32
	if is_seeder == "yes" {
		peers_count = torrent.Leechers
	} else {
		peers_count = torrent.Seeders + torrent.Leechers
	}

	// 设置最小返回时间，for client
	dt := time.Now()
	var announce_wait int64 = 300
	var announnce_interval int64 = 1000
	real_announce_interval := announnce_interval

	//判断种子时间，如果时间太长的，做冷处理，间隔时间为
	if (3600 > announce_wait) && (dt.Unix()-torrent.Added_time.Unix() > (30 * 86400)) {
		real_announce_interval = 3600
	} else {
		real_announce_interval = 2700
	}

	peers_res, peers_err := models.GetPeers(db, torrent.Torrent_id,
		is_seeder, peers_count, req.NumWant, req.InfoHash)

	if peers_err != nil {
		fmt.Printf("peers_err err : %v", peers_err)
		sendErrorResponse(c, errs.ERROR_GET_PEERS_DB)
	}

	var curr_peer *models.Peers //peer 类型

	var resp_peer_noid []RespPeerNoID
	var resp_peer []RespPeer

	if len(peers_res) != 0 { // 查询的peers不为空时候开始组装 ip，port 组
		for _, peer := range peers_res {
			if peer.Passkey == passkey { //查询出来的peer 包含本身，下一个
				curr_peer = &peer
				continue
			}

			if req.Compact { //true, url compact = 1

				s, s_err := packPeer(peer.IP, peer.Port)
				if s_err != nil {
					fmt.Printf("s_err %v\n", s_err)
				} else {
					t := fmt.Sprintf("%x", s)
					RespPeerCompact = append(RespPeerCompact, t)
				}

			} else {
				if req.NoPeerID {
					resp_peer_noid = append(resp_peer_noid, RespPeerNoID{
						IP:   peer.IP,
						Port: peer.Port,
					})
				} else {
					resp_peer = append(resp_peer, RespPeer{
						PeerID: peer.PeerID,
						IP:     peer.IP,
						Port:   peer.Port,
					})
				}

			}
		}
	}

	if curr_peer == nil {
		cr, _ := models.GetPeerByPasskey(db, torrent.Torrent_id, passkey)
		if len(cr) == 0 {
			fmt.Printf("no get peer by passkey ")
		} else {
			curr_peer = &cr[0]
		}
	}

	var update_set []string // 更新使用参数 数组

	if curr_peer == nil {
		var connectable string = "no"
		if CheckClientConnectable(req.IP, req.Port) {
			connectable = "yes"
		}

		rid, add_peer_err := models.AddNewPeer(db, torrent.Torrent_id, req.PeerID,
			req.IP, req.Port, req.Uploaded, req.Downloaded,
			req.Left, is_seeder, connectable, passkey, req.InfoHash)

		if add_peer_err != nil {
			sendErrorResponse(c, errs.MSG_ADD_PEER_ERR)
			return
		}
		log.Printf("add new peer  , peer_id : %v, torrent_id %v , insert lastid: %v", req.PeerID, torrent.Torrent_id, rid)

		if is_seeder == "yes" {
			update_set = append(update_set, "seeders = seeders + 1")
		} else {
			update_set = append(update_set, "leechers = leechers + 1")
		}

	} else {
		if req.Event == COMPLETED {
			update_set = append(update_set, "times_completed = times_completed + 1")

			err := models.UpdatePeer(db, torrent.Torrent_id, passkey,
				req.IP, req.Port, req.Uploaded, req.Downloaded, req.Left, is_seeder)

			if err != nil {
				sendErrorResponse(c, errs.MSG_UPDATE_PEER_ERR)
				return
			}

			if is_seeder != curr_peer.IsSeeder {
				if is_seeder == "yes" {
					update_set = append(update_set, "seeders = seeders + 1, leechers = leechers - 1")
				} else {
					update_set = append(update_set, "seeders = seeders - 1, leechers = leechers + 1")
				}
			}

		} else if req.Event == STOPPED {
			err := models.RemovePeerBy(db, torrent.Torrent_id, req.PeerID, passkey)
			if err != nil {
				sendErrorResponse(c, errs.MSG_REMOVE_PEER_FAILD)
				return
			}

			if is_seeder == "yes" {
				update_set = append(update_set, "seeders = seeders - 1")
			} else {
				update_set = append(update_set, "leechers = leechers - 1")
			}

		} else {
			err := models.UpdatePeer(db, torrent.Torrent_id, passkey,
				req.IP, req.Port, req.Uploaded, req.Downloaded, req.Left, is_seeder)
			if err != nil {
				sendErrorResponse(c, errs.MSG_UPDATE_PEER_ERR)
				return
			}
		}
	}

	if len(update_set) > 0 {
		update_set = append(update_set, fmt.Sprintf("last_action = '%s'", dt.Format("2006-01-02 15:04:05")))

		err := models.UpdateTorrent(db, update_set, torrent.Torrent_id)
		if err != nil {
			sendErrorResponse(c, errs.MSG_UPDATE_TORRENT_ERR)
			return
		}
	}
	//fmt.Printf("req %v, type %T\n", req, req)

	// 根据 compact ， no_peer_id 来匹配返回ip 结构
	var peers_finally interface{}
	if req.Compact {
		peers_finally = RespPeerCompact
	} else {
		if req.NoPeerID {
			peers_finally = resp_peer_noid
		} else {
			peers_finally = resp_peer
		}
	}

	af := Response{
		Interval:    real_announce_interval,
		MinInterval: announce_wait,
		Complete:    torrent.Seeders,
		Incomplete:  torrent.Leechers,
		Peers:       peers_finally,
	}
	var outBytes bytes.Buffer
	if err := bencode.Marshal(&outBytes, af); err != nil {
		sendErrorResponse(c, errs.MSG_GENERIC_ERROR)
		return
	}

	c.Data(http.StatusOK, gin.MIMEPlain, outBytes.Bytes())
}

func CheckClientConnectable(ip string, port uint16) bool {
	address := fmt.Sprintf("%s:%d", ip, port)

	// 设置连接超时时间，比如 3 秒
	timeout := 3 * time.Second
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}

	conn.Close()
	return true
}

func packPeer(ipStr string, port uint16) ([]byte, error) {
	ip := net.ParseIP(ipStr).To4() // IPv4 only
	if ip == nil {
		return nil, fmt.Errorf("invalid IPv4 address: %s", ipStr)
	}

	buf := make([]byte, 6)
	copy(buf[:4], ip)                         // 前4字节写 IP
	binary.BigEndian.PutUint16(buf[4:], port) // 后2字节写端口

	return buf, nil
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
