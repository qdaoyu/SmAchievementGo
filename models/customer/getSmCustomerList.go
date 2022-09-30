package customer

import (
	"log"
	"qiudaoyu/models"
	"qiudaoyu/models/menuInfo"
)

type Customer struct {
	Customerid string
	Name       string
	Gender     string
	// Consumetype  string
	Visittime    string
	Phone        string
	Shop         string
	Consultteach string
	// Item         string
	Treatnum   int
	Operanum   int
	Unoperanum int
}

// 获取尚美会员信息表
func GetSmCustomerList(uid int, username string, currentPage int, size int, cust Customertb) (map[string]interface{}, error) {
	var customer []Customer
	var user menuInfo.User
	var total int64

	//判断参数值，来构建where子句
	conditionParam := "Where 1=1 "
	if len(cust.Shop) > 0 {
		conditionParam += ` and t_customer.shop = "` + cust.Shop + `"`
	}
	if len(cust.Consultteach) > 0 {
		conditionParam += ` and t_customer.consultteach = "` + cust.Consultteach + `"`
	}
	if len(cust.Name) > 0 {
		conditionParam += ` and t_customer.name = "` + cust.Name + `"`
	}
	if len(cust.Phone) > 0 {
		conditionParam += ` and t_customer.phone = "` + cust.Phone + `"`
	}
	log.Println("conditionParam:", conditionParam)

	//存储信息
	smMap := make(map[string]interface{})

	err := models.Conn.Table("t_admin").Where("id = ? ", uid).Find(&user).Error
	if err != nil {
		log.Println(err)
		return nil, err
	}
	//分页的固定写法
	offsetVal := (currentPage - 1) * size

	rid := user.Roleid
	// log.Println("user:", user)
	// log.Println("roleid:", rid)
	if rid == 1 || rid == 2 {
		// sqlString := `select t_customer.*,t_orderlist.treatnum,t_orderlist.operanum,t_orderlist.unoperanum
		// from t_customer left join t_orderlist on t_customer.customerid = t_orderlist.customerid `
		// sqlString := `select t_customer.*,
		// t_orderlist.treatnum,t_orderlist.operanum ,t_orderlist.unoperanum
		// from t_customer left join t_orderlist on t_customer.customerid = t_orderlist.customerid limit ? offset  ? `
		sqlString := `select t_customer.customerid ,t_customer.name ,t_customer.gender ,
		t_customer.phone ,t_customer.shop,t_customer.consultteach , date_format(t_customer.visittime,"%Y-%m-%d") as visittime,
		t_orderlist.treatnum ,t_orderlist.operanum ,t_orderlist.unoperanum  
		from t_customer left join t_orderlist on t_customer.customerid = t_orderlist.customerid` + `  ` + conditionParam + ` limit ? offset  ? `

		sqlStringTotal := `select count(*) from t_customer` + `  ` + conditionParam

		models.Conn.Raw(sqlStringTotal).Count(&total)
		// err := models.Conn.Raw(sqlString).Count(&total).Limit(size).Offset(offsetVal).Scan(&customer).Error
		err := models.Conn.Raw(sqlString, size, offsetVal).Scan(&customer).Error
		log.Println(customer)
		// err := models.Conn.Table("t_syachieve").Count(&total).Limit(size).Offset(offsetVal).Find(&syAchieve).Error
		log.Println("条数：", total)
		if err != nil {
			log.Println(err)
			smMap["message"] = "查询失败" + err.Error()
			return smMap, err
		}
		// models.Conn.Raw(sqlString, userID).Scan(&user)
		// db.Raw(sqlString, userID).Create(&user)
		if len(customer) == 0 {
			log.Println("数据库无数据")
			smMap["message"] = "数据库无数据"
			return smMap, err
		} else {
			smMap["data"] = customer
			smMap["message"] = "查询成功"
			smMap["total"] = total
			return smMap, nil
		}
	} else {
		//limit offset
		sqlString := `select t_customer.customerid ,t_customer.name ,t_customer.gender ,
		t_customer.phone ,t_customer.shop,t_customer.consultteach , date_format(t_customer.visittime,"%Y-%m-%d") as visittime ,
		t_orderlist.treatnum,t_orderlist.operanum ,t_orderlist.unoperanum 
		from t_customer left join t_orderlist on t_customer.customerid = t_orderlist.customerid ` + `  ` + conditionParam + `and t_customer.consultteach = ? limit ? offset  ? `

		sqlStringTotal := `select count(*) from t_customer ` + `  ` + conditionParam + ` and t_customer.consultteach = ? `

		models.Conn.Raw(sqlStringTotal, username).Count(&total)
		err := models.Conn.Raw(sqlString, username, size, offsetVal).Scan(&customer).Error
		log.Println("条数：", total)
		if err != nil {
			log.Println(err)
			smMap["message"] = "查询失败" + err.Error()
			return smMap, err
		}
		// models.Conn.Raw(sqlString, userID).Scan(&user)
		// db.Raw(sqlString, userID).Create(&user)
		if len(customer) == 0 {
			log.Println("无此人数据")
			smMap["message"] = "查无数据"
			return smMap, err
		} else {
			smMap["data"] = customer
			smMap["message"] = "查询成功"
			smMap["total"] = total
			return smMap, nil
		}
	}

}
