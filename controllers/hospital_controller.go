package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hackshel/tracker-server/models"
	"github.com/hackshel/tracker-server/pkg/errs"
	"gorm.io/gorm"
)

func ListHospital(c *gin.Context) {
	userRole := c.GetString("userRole")
	//user_id := c.GetString("userID")
	fmt.Printf("run this ? ")

	if userRole != "admin" {
		c.JSON(http.StatusOK, gin.H{
			"error_code": errs.ERROR_USER_ROLE,
			"error_msg":  errs.GetMsg(errs.ERROR_USER_ROLE),
		})
		return
	}

	// 解析分页参数，默认 page=1, page_size=10
	pageStr := c.DefaultQuery("offset", "1")
	pageSizeStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	db := c.MustGet("db").(*gorm.DB)
	hospitals, total, err := models.GetHospitalList(db, page, pageSize)
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
		"rows":       hospitals,
	})
}

type AddHospitalRequest struct {
	HospitalName  string `json:"hospital_name" binding:"required"`
	HospitalID    int64  `json:"hospital_code" binding:"required"`
	HospitalGrade int32  `json:"hospital_grade" binding:"required"`
}

func AddHospital(c *gin.Context) {
	var req AddHospitalRequest

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

	hospital := models.Hospital{
		Code:         req.HospitalID,
		HspName:      req.HospitalName,
		Grade:        req.HospitalGrade,
		ProvinceId:   0,
		ProvinceName: "",
		CityId:       0,
		CityName:     "",
		CountyId:     0,
		CountyName:   "",
		AddTime:      time.Now(),
		LastModify:   time.Now(),
	}

	db := c.MustGet("db").(*gorm.DB)
	lastid, add_err := models.AddHospital(db, hospital)
	if add_err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_code": errs.ERROR_ADD_HOSPITAIL_FAILD,
			"error_msg":  errs.GetMsg(errs.ERROR_ADD_HOSPITAIL_FAILD),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"error_code": errs.SUCCESS,
		"error_msg":  errs.GetMsg(errs.SUCCESS),
		"user_id":    lastid,
	})
}

type DeleteHospitalRequest struct {
	Id int64
}

func DeleteHospital(c *gin.Context) {
	userRole := c.GetString("userRole")
	//user_id := c.GetString("userID")

	if userRole != "admin" {
		c.JSON(http.StatusOK, gin.H{
			"error_code": errs.ERROR_USER_ROLE,
			"error_msg":  errs.GetMsg(errs.ERROR_USER_ROLE),
		})
		return
	}

	var req DeleteHospitalRequest
	// 绑定 JSON 数据
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error_code": errs.ERROR_INVALID_DATA,
			"error_msg":  errs.GetMsg(errs.ERROR_INVALID_DATA),
		})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	err := models.DeleteHospitalByID(db, req.Id)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error_code": errs.ERROR_HOSPITAIL_DEL,
			"error_msg":  errs.GetMsg(errs.ERROR_HOSPITAIL_DEL),
		})
	}

	fmt.Printf("delete user OK\n")
	c.JSON(http.StatusOK, gin.H{
		"error_code": errs.SUCCESS,
		"error_msg":  errs.GetMsg(errs.SUCCESS),
	})
}
