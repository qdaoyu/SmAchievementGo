package basic

import (
	"log"
	"qiudaoyu/models"
)

type SmShop struct {
	Id       int
	Shopname string
}

type SmConsultteach struct {
	Id   int
	Name string
}

type SmItem struct {
	Id   int
	Item string
}

type SmConsumetype struct {
	Id          int
	Consumetype string
}

type SmCustomerid struct {
	Id         int
	Customerid string
}

// 获取尚美门店信息表
func GetSmShop() (map[string]interface{}, error) {
	var smshop []SmShop
	//存储信息
	resMp := make(map[string]interface{})
	//limit offset
	sqlString := `select t_smshop.* from t_smshop`
	err := models.Conn.Raw(sqlString).Scan(&smshop).Error
	if err != nil {
		log.Println(err)
		resMp["message"] = "查询失败" + err.Error()
		return resMp, err
	}
	// models.Conn.Raw(sqlString, userID).Scan(&user)
	// db.Raw(sqlString, userID).Create(&user)
	if len(smshop) == 0 {
		log.Println("无此人数据")
		resMp["message"] = "查无数据"
		return resMp, err
	} else {
		resMp["data"] = smshop
		resMp["message"] = "查询成功"
		return resMp, nil
	}

}

// 获取尚美咨询师信息表
func GetSmConsultteach() (map[string]interface{}, error) {
	var smconsultteach []SmConsultteach
	//存储信息
	resMp := make(map[string]interface{})
	//limit offset
	sqlString := `select t_consultteach.* from t_consultteach`
	err := models.Conn.Raw(sqlString).Scan(&smconsultteach).Error
	if err != nil {
		log.Println(err)
		resMp["message"] = "查询失败" + err.Error()
		return resMp, err
	}
	// models.Conn.Raw(sqlString, userID).Scan(&user)
	// db.Raw(sqlString, userID).Create(&user)
	if len(smconsultteach) == 0 {
		log.Println("查无数据")
		resMp["message"] = "查无数据"
		return resMp, err
	} else {
		resMp["data"] = smconsultteach
		resMp["message"] = "查询成功"
		return resMp, nil
	}

}

// 获取尚美项目信息表
func GetSmItem() (map[string]interface{}, error) {
	var smitem []SmItem
	//存储信息
	resMp := make(map[string]interface{})
	//limit offset
	sqlString := `select t_item.* from t_item`
	err := models.Conn.Raw(sqlString).Scan(&smitem).Error
	if err != nil {
		log.Println(err)
		resMp["message"] = "查询失败" + err.Error()
		return resMp, err
	}
	// models.Conn.Raw(sqlString, userID).Scan(&user)
	// db.Raw(sqlString, userID).Create(&user)
	if len(smitem) == 0 {
		log.Println("无项目数据")
		resMp["message"] = "查无数据"
		return resMp, err
	} else {
		resMp["data"] = smitem
		resMp["message"] = "查询成功"
		return resMp, nil
	}

}

// 获取尚美消费类型信息表
func GetSmConsumetype() (map[string]interface{}, error) {
	var smconsumetype []SmConsumetype
	//存储信息
	resMp := make(map[string]interface{})
	//limit offset
	sqlString := `select t_consumetype.* from t_consumetype`
	err := models.Conn.Raw(sqlString).Scan(&smconsumetype).Error
	if err != nil {
		log.Println(err)
		resMp["message"] = "查询失败" + err.Error()
		return resMp, err
	}
	// models.Conn.Raw(sqlString, userID).Scan(&user)
	// db.Raw(sqlString, userID).Create(&user)
	if len(smconsumetype) == 0 {
		log.Println("无消费类型数据")
		resMp["message"] = "查无数据"
		return resMp, err
	} else {
		resMp["data"] = smconsumetype
		resMp["message"] = "查询成功"
		return resMp, nil
	}

}

// 获取尚美会员id信息表
func GetSmCustomerid() (map[string]interface{}, error) {
	var smCustomerid []SmCustomerid
	//存储信息
	resMp := make(map[string]interface{})
	//limit offset
	sqlString := `select t_customer.id,t_customer.customerid from t_customer`
	err := models.Conn.Raw(sqlString).Scan(&smCustomerid).Error
	if err != nil {
		log.Println(err)
		resMp["message"] = "查询失败" + err.Error()
		return resMp, err
	}
	// models.Conn.Raw(sqlString, userID).Scan(&user)
	// db.Raw(sqlString, userID).Create(&user)
	if len(smCustomerid) == 0 {
		log.Println("无消费类型数据")
		resMp["message"] = "查无数据"
		return resMp, err
	} else {
		resMp["data"] = smCustomerid
		resMp["message"] = "查询成功"
		return resMp, nil
	}

}
