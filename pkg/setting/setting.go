package setting

import (
	"fmt"
	"log"
	"time"

	"github.com/go-ini/ini"
)

var (
	Cfg *ini.File

	RunMode string

	HTTPPort        int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	TorrentSavePath string
	TrackerURL      string

	PageSize  int
	JwtSecret string
	AppName   string

	RemoteURL     string
	RemoteBaseURL string
	RemoteSIG     string

	DSN          string
	TABLE_PREFIX string
)

func init() {

	var err error

	Cfg, err = ini.Load("conf/app.ini")

	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}

	LoadBase()
	LoadServer()
	LoadApp()
	LoadRemote()
	LoadDB()
}

func LoadBase() {
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")

}

func LoadServer() {

	sec, err := Cfg.GetSection("server")

	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}

	HTTPPort = sec.Key("HTTP_PORT").MustInt(8080)

	TorrentSavePath = sec.Key("TORRENT_SAVE_PATH").MustString("/opt/torrents/")

	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second

	WriteTimeout = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second

	TrackerURL = sec.Key("TRACKER_URL").MustString("default")
}

func LoadApp() {

	sec, err := Cfg.GetSection("app")

	if err != nil {
		log.Fatalf("Fail to get section 'app': %v", err)
	}

	JwtSecret = sec.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)")

	PageSize = sec.Key("PAGE_SIZE").MustInt(10)

	AppName = sec.Key("APP_NAME").MustString("default")
}

func LoadRemote() {

	r, err := Cfg.GetSection("remote")

	if err != nil {
		log.Fatalf("Fail to get section 'remote': %v", err)
	}

	RemoteURL = r.Key("URL").MustString("default")
	RemoteBaseURL = r.Key("BASE_URL").MustString("default")
	RemoteSIG = r.Key("SIG").MustString("default")

}

func LoadDB() {

	r, err := Cfg.GetSection("database")

	if err != nil {
		log.Fatalf("Fail to get section 'database' : %v", err)
	}

	//db_type := r.Key("TYPE").MustString("default")
	db_user := r.Key("USER").MustString("default")
	db_pass := r.Key("PASSWORD").MustString("default")
	db_host := r.Key("HOST").MustString("default")
	db_port := r.Key("PORT").MustString("default")
	db_name := r.Key("NAME").MustString("default")

	// 构建DSN(Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		db_user,
		db_pass,
		db_host,
		db_port,
		db_name,
	)
	//fmt.Printf("dsn config : %v\n", dsn)
	TABLE_PREFIX = r.Key("TABLE_PREFIX").MustString("default")
	DSN = dsn

}
