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
	// AddSyAchieveInfoHandler
	// e.GET("/achieve/sy/info/get", GetSyAchieveInfoHandler)
	//获取尚美会员列表
	e.GET("/customer/list", GetSmCustomerListHandler)
	// e.POST("/auth", AuthHandler)
	//获取门店信息
	e.GET("/basic/shop", GetShopHandler)
	//获取咨询师信息
	e.GET("/basic/consultteach", GetConsultteachHandler)
	//获取项目信息

	e.GET("/basic/item", GetItemHandler)
	//获取消费类型信息
	e.GET("/basic/consumetype", GetConsumetypeHandler)
	//添加会员信息
	e.POST("/customer/add", AddSmCustomerHandler)
	//删除会员信息
	e.DELETE("/customer/delete/:id", DeleteSmCustomerHandler)
	//更新会员信息
	e.PUT("/customer/update", UpdateSmCustomerHandler)
	//下载会员信息
	e.GET("/customer/export", ExportSmCustomerHandler)

}
