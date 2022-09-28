package customer

import (
	"log"
	"qiudaoyu/models"
)

// 会员表
//
//	type Customer struct {
//		Username string `form:"username" json:"username" binding:"required"`
//		Password string `form:"password" json:"password" binding:"required"`
//	}
//
// 新增会员信息
func AddSmCustomer(customer Customertb) error {
	var err error = models.Conn.Table("t_customer").Create(&customer).Error
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}
