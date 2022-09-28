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

// 更新会员信息
func UpdateSmCustomer(customer Customertb) error {
	// var cust Customertb
	sqlString := ` update  t_customer set    name = ? , gender = ?
	 , shop = ? , consultteach = ? , visittime = ? where phone = ? `
	// var err error = models.Conn.Table("t_customer").Where(`customerid = ? and  name = ? and gender = ?
	// 		and phone = ? and shop = ? and consultteach = ? and visittime = ? `, customer.Customerid, customer.Name,
	// 	customer.Gender, customer.Phone, customer.Shop, customer.Consultteach, customer.Visittime).Updates(&Customertb{}).Error
	err := models.Conn.Exec(sqlString, customer.Name, customer.Gender, customer.Shop, customer.Consultteach, customer.Visittime, customer.Phone).Error
	// 使用 Struct 进行 Select（会 select 零值的字段）
	// err := models.Conn.Table("t_customer").Select("Name", "Gender", "Phone", "Shop", "Consultteach", "Visittime").Where("customerid = ?", customer.Customerid).Updates(Customertb{Name: customer.Name,
	// 	Gender: customer.Gender, Phone: customer.Phone, Shop: customer.Shop, Consultteach: customer.Consultteach, Visittime: customer.Visittime}).Error
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}

// 删除会员信息
func DeleteSmCustomer(phone string) error {
	sqlString := `DELETE from t_customer where phone = ?`
	err := models.Conn.Exec(sqlString, phone).Error
	// err := models.Conn.Table("t_customer").Where("phone= ? ", phone).Delete(&Customertb{}).Error
	log.Println(phone)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}
