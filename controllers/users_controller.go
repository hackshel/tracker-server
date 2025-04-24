package controllers

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
