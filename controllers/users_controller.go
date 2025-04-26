package controllers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/hackshel/tracker-server/models"
	"github.com/hackshel/tracker-server/pkg/errs"
	"github.com/hackshel/tracker-server/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ErrorResp struct {
	Status   int    `json:"error_code"`
	ErrorMsg string `json:"error_msg"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	//fmt.Printf("req : %v\n", req)
	err := c.ShouldBindJSON(&req)
	/*
		fmt.Println("用户名:", req.Username)
		fmt.Println("密码:", req.Password)
	*/
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error_code": errs.ERROR_REQUEST,
			"error_msg":  errs.GetMsg(errs.ERROR_REQUEST),
		})
		return
	}

	sqldb := c.MustGet("db").(*gorm.DB)

	user, err := models.GetUserByUsername(sqldb, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_code": errs.ERROR_USER,
			"error_msg":  errs.GetMsg(errs.ERROR_USER),
		})
		return
	}

	//fmt.Printf("user %v \n", user)
	fmt.Printf("pass: %v\n\n", user.Role)

	dst := make([]byte, base64.StdEncoding.DecodedLen(len(user.Passwd)))
	n, b64_err := base64.StdEncoding.Decode(dst, []byte(user.Passwd))
	if b64_err != nil {
		log.Fatalf("decode error: %v", b64_err)
		return
	}
	dst = dst[:n]

	//fmt.Printf("+++: %v\n", dst)
	pass_err := bcrypt.CompareHashAndPassword([]byte(dst), []byte(req.Password))
	if pass_err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_code": errs.ERROR_USER_PASSWD,
			"error_msg":  errs.GetMsg(errs.ERROR_USER_PASSWD),
		})
		return
	}

	user_id := strconv.FormatInt(user.User_id, 10)

	token, token_err := utils.GenerateToken(req.Username, user.Role, user_id, user.Passkey)
	//fmt.Printf("errors: %v", token_err)
	if token_err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_code": errs.ERROR_GEN_TOKEN,
			"error_msg":  errs.GetMsg(errs.ERROR_GEN_TOKEN),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error_code": errs.SUCCESS,
		"error_msg":  errs.GetMsg(errs.SUCCESS),
		"token":      token,
	})
}

func Profile(c *gin.Context) {
	username := c.GetString("username")
	expire_at, _ := c.Get("expire_at")
	userRole, _ := c.Get("userRole")

	c.JSON(http.StatusOK, gin.H{
		"error_code": errs.SUCCESS,
		"error_msg":  "Authenticated",
		"username":   username,
		"userRole":   userRole,
		"expirt_at":  expire_at,
	})
}

func Logout(c *gin.Context) {

}

type ListUsrRequest struct {
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset" `
	Search   string `json:"search"`
	ListWith string `json:"listwith"`
}

func ListUser(c *gin.Context) {
	userRole := c.GetString("userRole")
	//user_id := c.GetString("userID")

	var req ListUsrRequest

	if userRole != "admin" {
		c.JSON(http.StatusOK, gin.H{
			"error_code": errs.ERROR_USER_ROLE,
			"error_msg":  errs.GetMsg(errs.ERROR_USER_ROLE),
		})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {

		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make(map[string]string)
			for _, fe := range ve {
				out[fe.Field()] = "字段验证错误，要求：" + fe.Tag()
			}
		}
		fmt.Printf("binding error: %v ,  out %v \n", err, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": errs.ERROR_INVALID_DATA,
			"error_msg":  errs.GetMsg(errs.ERROR_INVALID_DATA),
		})
		return
	}
	fmt.Printf("print req ===== %v, === type %T ===", req, req)
	fmt.Printf("req.Offset %v\n", req.Offset)
	fmt.Printf("req.Limit %v\n", req.Limit)
	fmt.Printf("req.search %v\n", req.Search)
	fmt.Printf("req.listwith %v\n", req.ListWith)

	// 解析分页参数，默认 page=1, page_size=10
	// pageStr := c.DefaultQuery("offset", "1")
	// pageSizeStr := c.DefaultQuery("limit", "10")

	// listWith := c.DefaultQuery("listwith", "")
	// fmt.Printf("listWith args %v, type %T\n", listWith, listWith)
	// page, err := strconv.Atoi(pageStr)
	// if err != nil || page <= 0 {
	// 	page = 1
	// }
	// pageSize, err := strconv.Atoi(pageSizeStr)
	// if err != nil || pageSize <= 0 {
	// 	pageSize = 10
	// }
	var page int
	var pageSize int
	if req.Offset <= 0 {
		page = 1
	} else {
		page = int(req.Offset)
	}
	if req.Limit <= 0 {
		pageSize = 10
	} else {
		pageSize = int(req.Limit)
	}
	listWith := req.ListWith
	db := c.MustGet("db").(*gorm.DB)
	//users, total, err := models.GetUserList(db, page, pageSize)
	// if err != nil {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"error_code": errs.ERROR_USER_LIST_FAILD,
	// 		"error_msg":  errs.GetMsg(errs.ERROR_USER_LIST_FAILD),
	// 	})
	// 	return
	// }

	var (
		rows  interface{}
		total int64
		err   error
	)
	if listWith == "hospital" {
		users, t, e := models.GetUserListWithHospital(db, page, pageSize)

		rows = users
		total = t
		err = e
	} else {
		users, t, e := models.GetUserList(db, page, pageSize)
		rows = users
		total = t
		err = e
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error_code": errs.ERROR_USER_LIST_FAILD,
			"error_msg":  errs.GetMsg(errs.ERROR_USER_LIST_FAILD),
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
		"rows":       rows,
		// "rows":       users,

	})
}

type AddUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	//HospitalID int64  `json:"hospital_id" binding:"required"`
}

func AddUser(c *gin.Context) {

	var req AddUserRequest

	userRole := c.GetString("userRole")
	//user_id := c.GetString("userID")

	if userRole != "admin" {
		c.JSON(http.StatusOK, gin.H{
			"error_code": errs.ERROR_USER_ROLE,
			"error_msg":  errs.GetMsg(errs.ERROR_USER_ROLE),
		})
		return
	}

	// 绑定 JSON 数据
	// c.ShouldBindJSON(&req)
	// fmt.Printf("binding json data %v, source data %v, type %T\n", req, req, req)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": errs.ERROR_INVALID_DATA,
			"error_msg":  errs.GetMsg(errs.ERROR_INVALID_DATA),
		})
		return
	}

	passwdhash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_code": errs.ERROR_PASSWD_HASH_FAILD,
			"error_msg":  errs.GetMsg(errs.ERROR_PASSWD_HASH_FAILD),
		})
	}
	salt := passwdhash[:29]
	fmt.Printf("slat %v, salt length %v, pass %v, pass length %v\n", salt, len(salt), passwdhash, len(passwdhash))
	saltStr := base64.StdEncoding.EncodeToString(salt)
	passwdStr := base64.StdEncoding.EncodeToString(passwdhash)

	user := models.Users{
		User_name: req.Username,
		// Code:         req.HospitalID,
		Code:         0,
		Role:         "doctor",
		Public_Key:   "",
		Access_Level: 1,
		Salt:         saltStr,
		Passwd:       passwdStr,
		Last_Login:   time.Now(),
		Passkey:      generateUUIDNoDash(),
	}

	fmt.Printf("copied user data for : %v, type %T\n", user, user)
	db := c.MustGet("db").(*gorm.DB)

	lastid, add_err := models.AddNewUser(db, user)
	if add_err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_code": errs.ERROR_PASSWD_HASH_FAILD,
			"error_msg":  errs.GetMsg(errs.ERROR_PASSWD_HASH_FAILD),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"error_code": errs.SUCCESS,
		"error_msg":  errs.GetMsg(errs.SUCCESS),
		"user_id":    lastid,
	})

}

func generateUUIDNoDash() string {
	u := uuid.New()
	return strings.ReplaceAll(u.String(), "-", "")
}

type ModiyfyUserHosRequest struct {
	UserId       int64 `json:"User_id" binding:"required"`
	HospitalCode int64 `json:"hospital_code" binding:"required"`
}

// 修改用户和医院关系
func ModifyUser(c *gin.Context) {
	userRole := c.GetString("userRole")
	//user_id := c.GetString("userID")

	if userRole != "admin" {
		c.JSON(http.StatusOK, gin.H{
			"error_code": errs.ERROR_USER_ROLE,
			"error_msg":  errs.GetMsg(errs.ERROR_USER_ROLE),
		})
		return
	}

	var req ModiyfyUserHosRequest
	// 绑定 JSON 数据
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error_code": errs.ERROR_INVALID_DATA,
			"error_msg":  errs.GetMsg(errs.ERROR_INVALID_DATA),
		})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	err := models.UpdateUserHospitalCode(db, req.UserId, req.HospitalCode)
	//err := models.UpdateUserFields(db, req.UserId, req.HospitalCode)
	if err != nil {
		fmt.Printf("update error %v", err)
		c.JSON(http.StatusOK, gin.H{
			"error_code": errs.ERROR_USER_UPDATE_CDDE,
			"error_msg":  errs.GetMsg(errs.ERROR_USER_UPDATE_CDDE),
		})
		return
	}

	fmt.Printf("delete user OK\n")
	c.JSON(http.StatusOK, gin.H{
		"error_code": errs.SUCCESS,
		"error_msg":  errs.GetMsg(errs.SUCCESS),
	})
}

type DlUserRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
}

func DeleteUser(c *gin.Context) {
	userRole := c.GetString("userRole")
	//user_id := c.GetString("userID")

	if userRole != "admin" {
		c.JSON(http.StatusOK, gin.H{
			"error_code": errs.ERROR_USER_ROLE,
			"error_msg":  errs.GetMsg(errs.ERROR_USER_ROLE),
		})
		return
	}

	var req DlUserRequest
	// 绑定 JSON 数据
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error_code": errs.ERROR_INVALID_DATA,
			"error_msg":  errs.GetMsg(errs.ERROR_INVALID_DATA),
		})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	err := models.DeleteUserByID(db, req.UserID)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error_code": errs.ERROR_USER_DEL,
			"error_msg":  errs.GetMsg(errs.ERROR_USER_DEL),
		})
		return
	}

	fmt.Printf("delete user OK\n")
	c.JSON(http.StatusOK, gin.H{
		"error_code": errs.SUCCESS,
		"error_msg":  errs.GetMsg(errs.SUCCESS),
	})

}
