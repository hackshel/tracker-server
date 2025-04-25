package errs

var MsgFlags = map[int]string{

	SUCCESS:        "ok",
	ERROR:          "fail",
	INVALID_PARAMS: "请求参数错误",

	ERROR_SERVER: "Server Error, 结构体转换错误",

	ERROR_EXIST_TAG:         "已存在该标签名称",
	ERROR_NOT_EXIST_TAG:     "该标签不存在",
	ERROR_NOT_EXIST_ARTICLE: "该文章不存在",

	ERROR_NOT_EXIST_CPL: "cpl_uuid 不存在",
	ERROR_NOT_EXIST_SN:  "server_sn 不存在",
	ERROR_NOT_EXIST_NVB: "start_time 不存在",
	ERROR_NOT_EXIST_NVA: "end_time 不存在",

	ERROR_REQUEST:                  "Request 错误",
	ERROR_AUTH_CHECK_TOKEN_FAIL:    "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT: "Token已超时",
	ERROR_AUTH_TOKEN:               "Token无效",
	ERROR_AUTH:                     "Token错误",
	ERROR_AUTH_HEADER:              "auth Header 丢失或无效", // "error": "Authorization header missing or invalid"

	ERROR_USER:           "用户不存在",
	ERROR_USER_PASSWD:    "用户密码错误",
	ERROR_GEN_TOKEN:      "创建JWT TOKEN 失败",
	ERROR_USER_NOT_FOUND: "没有找到用户",

	ERROR_TORRENT_INFO:              "提交的torrent_id 和 hash_info 错误",
	ERROR_TORRENT_NOT_FOUND:         "Torrent 没有找到",
	ERROR_DL_TORRENT_ID_OR_INFOHASH: "下载时提交的torrent_id 和 hash_info 错误",
	ERROR_TORRENT_READ_FAIL:         "种子读取失败",
	ERROR_TORRENT_DECODE_FAIL:       "种子文件解析错误",
	ERROR_UPLOAD_FAILD:              "上传文件失败",
	ERROR_OPEN_FILE_FAILD:           "打开文件失败",
	ERROR_SAVE_FILE_FAILD:           "保存文件失败",
	ERROR_INSERT_TORRENT_FAILD:      "写入数据库失败",
	ERROR_UPLOAD_SECCESS:            "上传成功",
	ERROR_ENCODE_FAILD:              "种子文件编码Bencode失败",
	ERROR_COUNT:                     "计算种子数量失败",
	ERROR_GET_PEERS_DB:              "从数据库获取PEERS 报错",

	MSG_INVALID_REQ_TYPE:   "客户端请求类型错误",
	MSG_MISSING_INFOHASH:   "no info_hash in request", //请求中没有info_hash
	MSG_MISSING_PEERID:     "no peer_id in request",   //请求中没有peer_id
	MSG_MISSING_PORT:       "no port in request",      //请求中没有端口号
	MSG_INVALID_PORT:       "invild port ",
	MSG_MISSING_DL:         "missing downloaded value",
	MSG_MISSING_UL:         "missing uploaded value",
	MSG_MISSING_LEFT:       "missing left value",
	MSG_MISSING_KEY:        "missing key value",
	MSG_INVALID_INFOHASH:   "invild info_hash", //的info_hash
	MSG_INVALID_PEERID:     "无效的Peer_id",
	MSG_INVALID_NUMWANT:    "无效的NumWant",
	MSG_BAD_CLIENT:         "非法客户端",
	MSG_INFOHASH_NOT_FOUND: "未知的info_hash",
	MSG_INVALID_AUTH:       "client request not passkey", // 客户端请求没有passkey ，鉴权失败
	//MSG_CLIENT_REQUEST_TOO_FAST: "客户端请求过快",
	MSG_GENERIC_ERROR:         "未知的错误，自动返回错误",
	MSG_MALFORMED_REQUEST:     "有问题的请求",
	MSG_QUERY_PARSE_FAIL:      "处理参数失败",
	MSG_PASSKEY_ERR:           "Passkey 长度小于40，错误Passkey",
	MSG_PORT_BANED:            "端口被Ban",
	MSG_OUT_MIN_ANNOUNCE_TIME: "超时 最小announce time",
	MSG_INSERT_PEER_ERR:       "插入peer 到数据库错误",
	MSG_ADD_PEER_ERR:          "insert client to tk_peers error.",
	MSG_UPDATE_PEER_ERR:       "PL Err 1",
	MSG_REMOVE_PEER_FAILD:     "D Err 1",
	MSG_UPDATE_TORRENT_ERR:    "update torrent table error",
}

func GetMsg(code int) string {

	msg, ok := MsgFlags[code]

	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
