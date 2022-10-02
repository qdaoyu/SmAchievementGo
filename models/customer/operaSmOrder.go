package customer

import (
	"log"
	"qiudaoyu/models"
	"qiudaoyu/models/menuInfo"
	"strconv"
)

// 获取尚美会员信息表
func GetSmOrderList(uid int, username string, currentPage int, size int, ordertb Ordertb) (map[string]interface{}, error) {
	var order []Ordertb
	var user menuInfo.User
	var total int64

	//判断参数值，来构建where子句
	conditionParam := "Where 1=1 "
	if len(ordertb.Item) > 0 {
		conditionParam += ` and t_orderlist.item = "` + ordertb.Item + `"`
	}
	if len(ordertb.Consultteach) > 0 {
		conditionParam += ` and t_orderlist.consultteach = "` + ordertb.Consultteach + `"`
	}
	if len(ordertb.Consumetype) > 0 {
		conditionParam += ` and t_orderlist.name = "` + ordertb.Consumetype + `"`
	}
	if ordertb.Unoperanum != -1 {
		conditionParam += ` and t_orderlist.unoperanum >= ` + strconv.Itoa(ordertb.Unoperanum)
	}
	log.Println("OrderConditionParam:", conditionParam)

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
		// sqlString := `select t_customer.*,t_order..list.treatnum,t_orderlist.operanum,t_orderlist.unoperanum
		// from t_customer left join t_orderlist on t_customer.customerid = t_orderlist.customerid `
		// sqlString := `select t_customer.*,
		// t_orderlist.treatnum,t_orderlist.operanum ,t_orderlist.unoperanum
		// from t_customer left join t_orderlist on t_customer.customerid = t_orderlist.customerid limit ? offset  ? `
		sqlString := `select t_orderlist.* from t_orderlist ` + `  ` + conditionParam + ` limit ? offset  ? `

		sqlStringTotal := `select count(*) from t_orderlist` + `  ` + conditionParam

		models.Conn.Raw(sqlStringTotal).Count(&total)
		// err := models.Conn.Raw(sqlString).Count(&total).Limit(size).Offset(offsetVal).Scan(&customer).Error
		err := models.Conn.Raw(sqlString, size, offsetVal).Scan(&order).Error
		log.Println("order:", order)
		// err := models.Conn.Table("t_syachieve").Count(&total).Limit(size).Offset(offsetVal).Find(&syAchieve).Error
		log.Println("条数：", total)
		if err != nil {
			log.Println(err)
			smMap["message"] = "查询失败" + err.Error()
			return smMap, err
		}
		// models.Conn.Raw(sqlString, userID).Scan(&user)
		// db.Raw(sqlString, userID).Create(&user)
		if len(order) == 0 {
			log.Println("数据库无数据")
			smMap["message"] = "数据库无数据"
			return smMap, err
		} else {
			smMap["data"] = order
			smMap["message"] = "查询成功"
			smMap["total"] = total
			return smMap, nil
		}
	} else {
		//limit offset
		sqlString := `select t_orderlist.* from t_orderlist  ` + `  ` + conditionParam + `and t_orderlist.consultteach = ? limit ? offset  ? `

		sqlStringTotal := `select count(*) from t_customer ` + `  ` + conditionParam + ` and t_orderlist.consultteach = ? `

		models.Conn.Raw(sqlStringTotal, username).Count(&total)
		err := models.Conn.Raw(sqlString, username, size, offsetVal).Scan(&order).Error
		log.Println("条数：", total)
		if err != nil {
			log.Println(err)
			smMap["message"] = "查询失败" + err.Error()
			return smMap, err
		}
		// models.Conn.Raw(sqlString, userID).Scan(&user)
		// db.Raw(sqlString, userID).Create(&user)
		if len(order) == 0 {
			log.Println("无此人数据")
			smMap["message"] = "查无数据"
			return smMap, err
		} else {
			smMap["data"] = order
			smMap["message"] = "查询成功"
			smMap["total"] = total
			return smMap, nil
		}
	}

}
