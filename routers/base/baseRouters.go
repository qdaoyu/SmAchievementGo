package base

import (
	"github.com/gin-gonic/gin"
)

func LoadBase(e *gin.Engine) {

	e.POST("/upload/suyan/achiveTable'", UploadSyAchieveTb)

	e.POST("/login", Login)

	e.POST("/calDupName", CalDupName)

	e.GET("/home", HomeMenuHandler)

	e.GET("/admin/info", GetUserInfoHandler)

	e.GET("/admin/basic/info", GetUserBasicInfoHandler)

	// e.POST("/achieve/sy/info/add", AddSyAchieveInfoHandler)

	// e.GET("/achieve/sy/info/get", GetSyAchieveInfoHandler)
	//尚美会员列表
	e.GET("/customer/list", GetSmCustomerListHandler)
	// e.POST("/auth", AuthHandler)

	// e.NoRoute(Page404)

}
