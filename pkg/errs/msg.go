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
}

func GetMsg(code int) string {

	msg, ok := MsgFlags[code]

	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
