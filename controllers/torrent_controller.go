package controllers

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hackshel/tracker-server/models"
	"github.com/hackshel/tracker-server/pkg/bencode"
	"github.com/hackshel/tracker-server/pkg/errs"
	"github.com/hackshel/tracker-server/pkg/setting"
	"gorm.io/gorm"
)

type TorrentInfoRequest struct {
	TorrentID int64  `json:"torrent_id" binding:"required"`
	InfoHash  string `json:"info_hash" binding:"required"`
}
type TorrentDownloadRequest struct {
	TorrentID int64  `json:"torrent_id"`
	InfoHash  string `json:"info_hash"`
}

type FileInfo struct {
	Length int64    `bencode:"length"`
	Path   []string `bencode:"path"` // 对于多层目录，path 是一个数组
}

type InfoDict struct {
	Name        string     `bencode:"name"`
	PieceLength int        `bencode:"piece length"`
	Pieces      string     `bencode:"pieces"`
	Length      int64      `bencode:"length,omitempty"` // 单文件才有
	Files       []FileInfo `bencode:"files,omitempty"`  // 多文件才有
	Private     int        `bencode:"private,omitempty"`
}

type TorrentFile struct {
	Announce     string     `bencode:"announce"`
	AnnounceList [][]string `bencode:"announce-list,omitempty"`
	CreationDate int64      `bencode:"creation date,omitempty"` // Unix 时间戳
	CreatedBy    string     `bencode:"created by,omitempty"`
	Encoding     string     `bencode:"encoding,omitempty"` // 一般是 "UTF-8"
	Comment      string     `bencode:"comment,omitempty"`
	Info         InfoDict   `bencode:"info"`
}

// /api/v1/torrent/list?limit=20&offset=0&navStatus=&start_date=2025-03-25&end_date=2025-04-23&search=&_=1745402592007
func TorrentList(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	// 解析分页参数，默认 page=1, page_size=10
	pageStr := c.DefaultQuery("offset", "1")
	pageSizeStr := c.DefaultQuery("limit", "10")

	userRole := c.GetString("userRole")

	user_id := c.GetString("userID")
	userid_int, _ := strconv.ParseInt(user_id, 10, 64)

	if userRole == "admin" {
		userid_int = -1
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	torrents, total, err := models.GetTorrentList(db, page, pageSize, userid_int)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get torrent list",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error_code": errs.SUCCESS,
		"error_msg":  errs.GetMsg(errs.SUCCESS),
		"page":       page,
		"page_size":  pageSize,
		"total":      total,
		"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		"rows":       torrents,
	})
}

func TorrentsCount(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	user_id := c.GetString("userID")
	fmt.Printf("userid ==> %v , type %T\n", user_id, user_id)

	userid_int, _ := strconv.ParseInt(user_id, 10, 64)

	torrent_count, err := models.GetTorentsCountByUserID(db, userid_int)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": errs.ERROR_COUNT,
			"error_msg":  errs.GetMsg(errs.ERROR_COUNT),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"error_code": errs.SUCCESS,
		"error_msg":  errs.GetMsg(errs.SUCCESS),
		"cnt":        torrent_count,
	})
}

func TorrenInfo(c *gin.Context) {
	var req TorrentInfoRequest

	// 绑定 JSON 数据
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": errs.ERROR_TORRENT_INFO,
			"error_msg":  errs.GetMsg(errs.ERROR_TORRENT_INFO),
		})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	torrent, err := models.GetTorrentByIDAndHash(db, req.TorrentID, req.InfoHash)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error_code": errs.ERROR_TORRENT_NOT_FOUND,
			"error_msg":  errs.GetMsg(errs.ERROR_TORRENT_NOT_FOUND),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error_code": errs.SUCCESS,
		"error_msg":  errs.GetMsg(errs.SUCCESS),
		"data":       torrent,
	})
}

func TorrenDownload(c *gin.Context) {

	// var req TorrentDownloadRequest
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error_code": errs.ERROR_DL_TORRENT_ID_OR_INFOHASH,
	// 		"error_msg":  errs.GetMsg(errs.ERROR_DL_TORRENT_ID_OR_INFOHASH),
	// 	})
	// 	return
	// }

	TorrentIDStr := c.Query("torrent_id")
	InfoHash := c.Query("info_hash")

	TorrentID, _ := strconv.ParseInt(TorrentIDStr, 10, 64)
	db := c.MustGet("db").(*gorm.DB)

	// add logs 20250427
	user_id := c.GetString("userID")
	fmt.Printf("userid ==> %v , type %T\n", user_id, user_id)
	userid_int, _ := strconv.ParseInt(user_id, 10, 64)

	log := models.AuditLogs{
		UserID:         userid_int,
		OperationType:  "DOWNLOAD",
		TargetTorrent:  TorrentID,
		Detail:         "DOWNLOAD",
		Timestamp:      time.Now(),
		BlockchainHash: "",
	}
	log_err := models.AddLogs(db, log)
	if log_err != nil {
		fmt.Printf("Logout system!\n")
	}
	// end logs

	// torrent, err := models.GetTorrentByIDAndHash(db, req.TorrentID, req.InfoHash)
	// fmt.Printf("torrent info %v, type %T \n", torrent, torrent)
	torrent, err := models.GetTorrentByIDAndHash(db, TorrentID, InfoHash)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error_code": errs.ERROR_TORRENT_NOT_FOUND,
			"error_msg":  errs.GetMsg(errs.ERROR_TORRENT_NOT_FOUND),
		})
		return
	}

	// 获取当前登录用户的 passkey
	passkey := c.GetString("PassKey")
	//fmt.Printf("passky: %v", passkey)

	// 加载 torrent 文件
	torrentPath := filepath.Join(setting.TorrentSavePath, torrent.Filename)
	fp, err := os.Open(torrentPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_code": errs.ERROR_TORRENT_READ_FAIL,
			"error_msg":  errs.GetMsg(errs.ERROR_TORRENT_READ_FAIL),
		})
		return
	}
	defer fp.Close()

	// 3. 解析 .torrent 文件
	var tf TorrentFile

	if err := bencode.Unmarshal(fp, &tf); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": errs.ERROR_TORRENT_DECODE_FAIL,
			"error_msg":  errs.GetMsg(errs.ERROR_TORRENT_DECODE_FAIL),
			"detail":     err.Error(),
		})
		return
	}
	var builder strings.Builder
	builder.WriteString(setting.TrackerURL)
	builder.WriteString("/api/v1/tracker/announce?passkey=")
	builder.WriteString(passkey)
	tracker_url_with_passkey := builder.String()

	//fmt.Printf("new tracker %v, type %T", tracker_url_with_passkey, tracker_url_with_passkey)

	trackerList := [][]string{
		{tracker_url_with_passkey},
	}
	tf.Announce = tracker_url_with_passkey
	tf.AnnounceList = trackerList

	// 编码为 bencode 并返回
	var buf bytes.Buffer
	if err := bencode.Marshal(&buf, tf); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_code": errs.ERROR_ENCODE_FAILD,
			"error_msg":  errs.GetMsg(errs.ERROR_ENCODE_FAILD),
		})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", torrent.Filename))
	c.Header("Content-Type", "application/x-bittorrent")
	c.Data(http.StatusOK, "application/x-bittorrent", buf.Bytes())

}

func TorrentUpload(c *gin.Context) {
	// 1. 获取用户信息
	userID := c.GetString("userID")
	//username := c.GetString("username")

	userid_int, _ := strconv.ParseInt(userID, 10, 64)
	//fmt.Printf("user info => id: %v type: %T  name: %v\n", userID, userID, username)
	file, err := c.FormFile("torrent")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": errs.ERROR_UPLOAD_FAILD,
			"error_msg":  errs.GetMsg(errs.ERROR_UPLOAD_FAILD),
		})
		return
	}

	// torrent_size := file.Size
	torrent_filename := file.Filename

	//fmt.Printf("file object filename: %v, filename type : %T\n", torrent_filename, torrent_filename)
	//fmt.Printf("file object filesize: %v, filesize type : %T\n", torrent_size, torrent_size)
	fp, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_code": errs.ERROR_OPEN_FILE_FAILD,
			"error":      errs.GetMsg(errs.ERROR_OPEN_FILE_FAILD),
		})
		return
	}
	defer fp.Close()

	// 3. 解析 .torrent 文件
	var tf TorrentFile

	if err := bencode.Unmarshal(fp, &tf); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": errs.ERROR_TORRENT_DECODE_FAIL,
			"error_msg":  errs.GetMsg(errs.ERROR_TORRENT_DECODE_FAIL),
			"detail":     err.Error(),
		})
		return
	}

	// fmt.Printf("torrent File announce %v\n", tf.Announce)
	// fmt.Printf("torrent File announce-list %v\n", tf.AnnounceList)
	// fmt.Printf("torrent File  CreatedBy %v\n", tf.CreatedBy)
	// fmt.Printf("torrent File Comment %v\n", tf.Comment)
	// fmt.Printf("torrent File Encoding %v\n", tf.Encoding)
	// fmt.Printf("torrent File info.name %v\n", tf.Info.Name)
	// fmt.Printf("torrent File info.Pieces %T\n", tf.Info.Pieces)
	// fmt.Printf("torrent File info.PieceLength %v\n", tf.Info.PieceLength)
	// fmt.Printf("torrent File info.Length %v\n", tf.Info.Length)
	// fmt.Printf("torrent File info.Private %v\n", tf.Info.Private)

	r, cnt := SplitPieces(tf.Info.Pieces)

	log.Printf("picecs type %T\n", r)
	// fmt.Printf("picecs count : %v", cnt)

	// fmt.Printf("rr %T", r)
	/*
		for i, p := range r {
			fmt.Printf("Piece #%d: %s\n", i+1, hex.EncodeToString(p))
		}*/
	// 4. 保存文件到磁盘
	//filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	savePath := filepath.Join(setting.TorrentSavePath, torrent_filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_code": errs.ERROR_SAVE_FILE_FAILD,
			"error_msg":  errs.GetMsg(errs.ERROR_SAVE_FILE_FAILD),
		})
		return
	}

	file_sha1 := GetTorrentSHA1(savePath)
	// fmt.Printf("file sha1 : %v\n", file_sha1)

	// 5. 生成 info_hash（可使用解析库或 sha1 计算 info 字段）
	// 这里简化处理：
	infoHash, hash_err := GenerateInfoHash(tf.Info)
	if hash_err != nil {
		log.Printf("error : %v\n", hash_err)
	}
	// fmt.Printf("info hash sha1 %v \n", infoHash)

	ssize, f_cnt := GetPackageSize(tf.Info.Files)
	// fmt.Printf("package size : %v, file count: %v", ssize, f_cnt)

	// last_action := time.Now().Format("2006-01-02 15:04:05")
	// fmt.Printf("last action :%v\n", last_action)

	// 6. 保存到数据库
	db := c.MustGet("db").(*gorm.DB)

	t := models.Torrents{
		Info_hash:       infoHash,
		Name:            torrent_filename,
		Filename:        torrent_filename,
		Dcp_uuid:        tf.Comment,
		Dcp_size:        ssize,
		Piece_length:    tf.Info.PieceLength,
		Pices_count:     cnt,
		Added_time:      time.Now(),
		Dcp_type:        "single",
		Numfiles:        f_cnt,
		Tracker_url:     setting.TrackerURL + "/api/v1/tracker/announce",
		F_sha1:          file_sha1,
		Seeders:         0,
		Leechers:        0,
		Times_completed: 0,
		Last_action:     time.Now(),
		User_id:         userid_int,
	}

	if err := db.Create(&t).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_code": errs.ERROR_INSERT_TORRENT_FAILD,
			"error_msg":  errs.GetMsg(errs.ERROR_INSERT_TORRENT_FAILD),
		})
		return
	}
	// add logs 20250427

	log := models.AuditLogs{
		UserID:         userid_int,
		OperationType:  "UPLOAD",
		TargetTorrent:  t.Torrent_id,
		Detail:         "UPLOAD",
		Timestamp:      time.Now(),
		BlockchainHash: "",
	}
	log_err := models.AddLogs(db, log)
	if log_err != nil {
		fmt.Printf("Logout system!\n")
	}

	c.JSON(http.StatusOK, gin.H{
		"error_code": errs.ERROR_UPLOAD_SECCESS,
		"error_msg":  errs.GetMsg(errs.ERROR_UPLOAD_SECCESS),
		"torrent_id": t.Torrent_id,
	})

}

type TorrentDeleteRequest struct {
	TorrentId int64 `json:"torrent_id"`
}

func TorrentDelete(c *gin.Context) {

	var req TorrentDeleteRequest
	userID := c.GetString("userID")
	//username := c.GetString("username")
	userid_int, _ := strconv.ParseInt(userID, 10, 64)

	// 绑定 JSON 数据
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": errs.ERROR_TORRENT_INFO,
			"error_msg":  errs.GetMsg(errs.ERROR_TORRENT_INFO),
		})
		return
	}
	db := c.MustGet("db").(*gorm.DB)

	torrent, err := models.GetTorrentByID(db, req.TorrentId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error_code": errs.ERROR_TORRENT_NOT_FOUND,
			"error_msg":  errs.GetMsg(errs.ERROR_TORRENT_NOT_FOUND),
		})
		return
	}
	//删除磁盘文件
	savePath := filepath.Join(setting.TorrentSavePath, torrent.Filename)

	rm_err := os.Remove(savePath)
	if rm_err != nil {
		fmt.Println("删除失败:", rm_err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": errs.ERROR_TORRENT_RM_FAILD,
			"error_msg":  errs.GetMsg(errs.ERROR_TORRENT_RM_FAILD),
		})

		return
	}

	// add logs 20250427

	log := models.AuditLogs{
		UserID:         userid_int,
		OperationType:  "DELETE",
		TargetTorrent:  req.TorrentId,
		Detail:         "DELETE",
		Timestamp:      time.Now(),
		BlockchainHash: "",
	}
	log_err := models.AddLogs(db, log)
	if log_err != nil {
		fmt.Printf("Logout system!\n")
	}
	del_db_err := models.DeleteTorrentByTorrentID(db, req.TorrentId)

	if del_db_err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": errs.ERROR_TORRENT_DEL_FAILD,
			"error_msg":  errs.GetMsg(errs.ERROR_TORRENT_DEL_FAILD),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error_code": errs.SUCCESS,
		"error_msg":  errs.GetMsg(errs.SUCCESS),
	})
}

func GetTorrentSHA1(filepath string) string {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	sha1 := fmt.Sprintf("%x", h.Sum(nil))
	return sha1
}

func GetPackageSize(f_info []FileInfo) (int64, int) {

	var sum int64 = 0
	var f_count int = 0
	// 多文件模式
	if len(f_info) > 0 {
		for _, file := range f_info {
			sum += file.Length
			f_count += 1
		}
	}
	return sum, f_count
}

func SplitPieces(pieces string) ([][]byte, int) {
	count := 0
	pieceLen := 20
	pieceBytes := []byte(pieces)
	var result [][]byte

	for i := 0; i < len(pieceBytes); i += pieceLen {
		if i+pieceLen > len(pieceBytes) {
			break
		}
		result = append(result, pieceBytes[i:i+pieceLen])
		count += 1
	}
	return result, count
}

func GenerateInfoHash(info InfoDict) (string, error) {
	// 用 bencode 编码 info 字典
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, info)
	if err != nil {
		return "", fmt.Errorf("bencode marshal failed: %v", err)
	}

	// 计算 SHA1 哈希
	hash := sha1.Sum(buf.Bytes())

	// 转为 16 进制字符串（也可以返回 []byte）
	return fmt.Sprintf("%x", hash), nil
}
