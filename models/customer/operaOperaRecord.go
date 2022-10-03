package customer

import (
	"errors"
	"log"
	"qiudaoyu/models"
	"qiudaoyu/models/menuInfo"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// 获取操作记录表
func GetOperaList(uid int, username string, currentPage int, size int, ordertb Ordertb, name string) (map[string]interface{}, error) {
	var order []Ordertb
	var user menuInfo.User
	var total int64

	//判断参数值，来构建where子句
	conditionParam := "Where 1=1 "
	if len(name) > 0 {
		conditionParam += ` and t_orderlist.customerid like "%` + name + `%"`
	}
	if len(ordertb.Item) > 0 {
		conditionParam += ` and t_orderlist.item = "` + ordertb.Item + `"`
	}
	if len(ordertb.Consultteach) > 0 {
		conditionParam += ` and t_orderlist.consultteach = "` + ordertb.Consultteach + `"`
	}
	if len(ordertb.Consumetype) > 0 {
		conditionParam += ` and t_orderlist.consumetype = "` + ordertb.Consumetype + `"`
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
		sqlString := `select t_orderlist.* from t_orderlist  ` + `  ` + conditionParam + `  ` + `and t_orderlist.consultteach = ? limit ? offset  ? `

		sqlStringTotal := `select count(*) from t_orderlist ` + `  ` + conditionParam + `  ` + ` and t_orderlist.consultteach = ? `

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

// 新增操作记录
func AddOpera(order Ordertb) error {
	//判断顾客的新老是否选择正确,在见诊日期之前就已经存在了
	var countRows int64
	// sqlString := `select t_orderlist.* from t_orderlist ` + `  ` + conditionParam + ` limit ? offset  ? `

	sqlStringTotal := `select count(*) from t_orderlist where customerid = "` + order.Customerid + `"  and visittime <  "` + order.Visittime + `"`

	err := models.Conn.Raw(sqlStringTotal).Count(&countRows).Error
	log.Println("sqlStringTotal:", sqlStringTotal)
	log.Println("今天以前的记录数countRows:")
	if err != nil {
		log.Println(err)
		return err
	}
	if countRows >= 1 && strings.Contains(order.Consumetype, "新客") {
		log.Println("此顾客历史已消费过，非新客")
		return errors.New("此顾客历史已消费过，非新客")
	}
	err = models.Conn.Table("t_orderlist").Create(&order).Error
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}

// 更新操作记录
func UpdateOpera(order Ordertb) error {
	//判断顾客的新老是否选择正确,在见诊日期之前就已经存在了
	var countRows int64
	sqlStringTotal := `select count(*) from t_orderlist where customerid = "` + order.Customerid + `"  and visittime <  "` + order.Visittime + `"`

	err := models.Conn.Raw(sqlStringTotal).Count(&countRows).Error
	log.Println("sqlStringTotal:", sqlStringTotal)
	log.Println("今天以前的记录数countRows:")
	if err != nil {
		log.Println(err)
		return err
	}
	if countRows >= 1 && strings.Contains(order.Consumetype, "新客") {
		log.Println("此顾客历史已消费过，非新客")
		return errors.New("此顾客历史已消费过，非新客")
	}
	// var cust Customertb
	sqlString := ` update  t_orderlist set    visittime = ? , customerid = ?
	 , consumetype = ? , item = ? , treatnum = ? , swmdmoney = ?  , qbmoney = ? , mmmoney = ? , mmboxnum = ? , 
	 donateboxnum = ? , dgfmoney = ? , xfymoney = ? , totalmoney = ?  , returnmoney = ? , kkmoney = ? , owedmoney = ?  
	 , paidmoney = ? , consultteach = ? , dqteach = ?  , operateach = ? , comment = ?  where orderid = ?  `
	// log.Println("order:", order)

	err = models.Conn.Exec(sqlString, order.Visittime, order.Customerid, order.Consumetype, order.Item,
		order.Treatnum, order.Swmdmoney, order.Qbmoney, order.Mmmoney, order.Mmboxnum, order.Donateboxnum, order.Dgfmoney,
		order.Xfymoney, order.Totalmoney, order.Returnmoney, order.Kkmoney, order.Owedmoney, order.Paidmoney, order.Consultteach,
		order.Dqteach, order.Operateach, order.Comment, order.Orderid).Error
	// 使用 Struct 进行 Select（会 select 零值的字段）
	// err := models.Conn.Table("t_customer").Select("Name", "Gender", "Phone", "Shop", "Consultteach", "Visittime").Where("customerid = ?", customer.Customerid).Updates(Customertb{Name: customer.Name,
	// 	Gender: customer.Gender, Phone: customer.Phone, Shop: customer.Shop, Consultteach: customer.Consultteach, Visittime: customer.Visittime}).Error
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}

// 删除操作记录
func DeleteOpera(orderId string) error {
	sqlString := `DELETE from t_orderlist where orderid = ?`
	err := models.Conn.Exec(sqlString, orderId).Error
	// err := models.Conn.Table("t_customer").Where("phone= ? ", phone).Delete(&Customertb{}).Error
	log.Println(orderId)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}

// 导出操作记录
// 获取操作记录信息表,返回文件地址
func ExportOperaList(uid int, username string) (string, error) {
	var order []Ordertb
	var user menuInfo.User
	//存储信息
	smMap := make(map[string]interface{})

	err := models.Conn.Table("t_admin").Where("id = ? ", uid).Find(&user).Error
	if err != nil {
		log.Println(err)
		return "", err
	}
	rid := user.Roleid

	if rid == 1 || rid == 2 {
		//, date_format(t_orderlist.visittime,"%Y-%m-%d") as visittimeFormat,
		sqlString := `select t_orderlist.*  
		from t_orderlist `

		// err := models.Conn.Raw(sqlString).Count(&total).Limit(size).Offset(offsetVal).Scan(&customer).Error
		err := models.Conn.Raw(sqlString).Scan(&order).Error
		if err != nil {
			log.Println(err)
			return "查询失败:" + err.Error(), err
		}
		// models.Conn.Raw(sqlString, userID).Scan(&user)
		// db.Raw(sqlString, userID).Create(&user)
		if len(order) == 0 {
			log.Println("数据库无数据")
			smMap["message"] = "数据库无数据"
			return "", errors.New("数据库无数据")
		} else {
			var filePath string = "./assets/exportFile/" + username + strconv.FormatInt(time.Now().Unix(), 10) + ".xlsx"
			//数据库数据整理到表格中
			f := excelize.NewFile()
			// 设置单元格的值
			f.SetCellValue("Sheet1", "A1", "订单编号")
			f.SetCellValue("Sheet1", "B1", "见诊日期")
			f.SetCellValue("Sheet1", "C1", "会员id")
			f.SetCellValue("Sheet1", "D1", "消费类型")
			f.SetCellValue("Sheet1", "E1", "项目")
			f.SetCellValue("Sheet1", "F1", "疗程次数")
			f.SetCellValue("Sheet1", "G1", "四维美雕金额")
			f.SetCellValue("Sheet1", "H1", "祛斑金额")
			f.SetCellValue("Sheet1", "I1", "面膜金额")
			f.SetCellValue("Sheet1", "J1", "面膜盒数")
			f.SetCellValue("Sheet1", "K1", "赠送盒数")
			f.SetCellValue("Sheet1", "L1", "冻干粉金额")
			f.SetCellValue("Sheet1", "M1", "修复液金额")
			f.SetCellValue("Sheet1", "N1", "合计成交金额")
			f.SetCellValue("Sheet1", "O1", "回款金额")
			f.SetCellValue("Sheet1", "P1", "卡扣金额")
			f.SetCellValue("Sheet1", "Q1", "欠款金额")
			f.SetCellValue("Sheet1", "R1", "实付金额")
			f.SetCellValue("Sheet1", "S1", "咨询师")
			f.SetCellValue("Sheet1", "T1", "导前师")
			f.SetCellValue("Sheet1", "U1", "操作师")
			f.SetCellValue("Sheet1", "V1", "已操作次数")
			f.SetCellValue("Sheet1", "W1", "未操作次数")
			f.SetCellValue("Sheet1", "X1", "备注")

			//数据库取值
			// sqlStringTotal := `select * from t_customer`
			// err := models.Conn.Raw(sqlStringTotal).Scan(&order).Error
			// if err != nil {
			// 	log.Println(err)
			// 	return "", err
			// }
			var i int = 2
			for _, value := range order {
				// log.Println(key, "---", value)
				f.SetCellValue("Sheet1", "A"+strconv.Itoa(i), value.Orderid)
				f.SetCellValue("Sheet1", "B"+strconv.Itoa(i), strings.Split(value.Visittime, "T")[0])
				f.SetCellValue("Sheet1", "C"+strconv.Itoa(i), value.Customerid)
				f.SetCellValue("Sheet1", "D"+strconv.Itoa(i), value.Consumetype)
				f.SetCellValue("Sheet1", "E"+strconv.Itoa(i), value.Item)
				f.SetCellValue("Sheet1", "F"+strconv.Itoa(i), value.Treatnum)
				f.SetCellValue("Sheet1", "G"+strconv.Itoa(i), value.Swmdmoney)
				f.SetCellValue("Sheet1", "H"+strconv.Itoa(i), value.Qbmoney)
				f.SetCellValue("Sheet1", "I"+strconv.Itoa(i), value.Mmmoney)
				f.SetCellValue("Sheet1", "J"+strconv.Itoa(i), value.Mmboxnum)
				f.SetCellValue("Sheet1", "K"+strconv.Itoa(i), value.Donateboxnum)
				f.SetCellValue("Sheet1", "L"+strconv.Itoa(i), value.Dgfmoney)
				f.SetCellValue("Sheet1", "M"+strconv.Itoa(i), value.Xfymoney)
				f.SetCellValue("Sheet1", "N"+strconv.Itoa(i), value.Totalmoney)
				f.SetCellValue("Sheet1", "O"+strconv.Itoa(i), value.Returnmoney)
				f.SetCellValue("Sheet1", "P"+strconv.Itoa(i), value.Kkmoney)
				f.SetCellValue("Sheet1", "Q"+strconv.Itoa(i), value.Owedmoney)
				f.SetCellValue("Sheet1", "R"+strconv.Itoa(i), value.Paidmoney)
				f.SetCellValue("Sheet1", "S"+strconv.Itoa(i), value.Consultteach)
				f.SetCellValue("Sheet1", "T"+strconv.Itoa(i), value.Dqteach)
				f.SetCellValue("Sheet1", "U"+strconv.Itoa(i), value.Operateach)
				f.SetCellValue("Sheet1", "V"+strconv.Itoa(i), value.Operanum)
				f.SetCellValue("Sheet1", "W"+strconv.Itoa(i), value.Unoperanum)
				f.SetCellValue("Sheet1", "X"+strconv.Itoa(i), value.Comment)
				i = i + 1

			}
			err = f.SaveAs(filePath)
			if err != nil {
				log.Println(err)
				return "订单写入excel失败", err
			}
			// excelReader, _ := os.Open(filePath)
			// fileInfo, _ := excelReader.Stat()
			return filePath, nil
		}
	} else {
		//limit offset
		sqlString := `select t_orderlist.* from t_orderlist where t_orderlist.consultteach = ? `

		err := models.Conn.Raw(sqlString, username).Scan(&order).Error
		if err != nil {
			log.Println(err)
			return "下载会员信息失败", err
		}
		// models.Conn.Raw(sqlString, userID).Scan(&user)
		// db.Raw(sqlString, userID).Create(&user)
		if len(order) == 0 {
			return "无数据,无法下载", errors.New("无数据,无法下载")
		} else {
			var filePath string = "./assets/exportFile/" + username + strconv.FormatInt(time.Now().Unix(), 10) + ".xlsx"
			// var filePath string = "/exportFile/" + strconv.FormatInt(time.Now().Unix(), 10) + ".xlsx"
			// var filePath string = "d:/test.xlsx"
			//数据库数据整理到表格中
			f := excelize.NewFile()
			// 设置单元格的值
			f.SetCellValue("Sheet1", "A1", "订单编号")
			f.SetCellValue("Sheet1", "B1", "见诊日期")
			f.SetCellValue("Sheet1", "C1", "会员id")
			f.SetCellValue("Sheet1", "D1", "消费类型")
			f.SetCellValue("Sheet1", "E1", "项目")
			f.SetCellValue("Sheet1", "F1", "疗程次数")
			f.SetCellValue("Sheet1", "G1", "四维美雕金额")
			f.SetCellValue("Sheet1", "H1", "祛斑金额")
			f.SetCellValue("Sheet1", "I1", "面膜金额")
			f.SetCellValue("Sheet1", "J1", "面膜盒数")
			f.SetCellValue("Sheet1", "K1", "赠送盒数")
			f.SetCellValue("Sheet1", "L1", "冻干粉金额")
			f.SetCellValue("Sheet1", "M1", "修复液金额")
			f.SetCellValue("Sheet1", "N1", "合计成交金额")
			f.SetCellValue("Sheet1", "O1", "回款金额")
			f.SetCellValue("Sheet1", "P1", "卡扣金额")
			f.SetCellValue("Sheet1", "Q1", "欠款金额")
			f.SetCellValue("Sheet1", "R1", "实付金额")
			f.SetCellValue("Sheet1", "S1", "咨询师")
			f.SetCellValue("Sheet1", "T1", "导前师")
			f.SetCellValue("Sheet1", "U1", "操作师")
			f.SetCellValue("Sheet1", "V1", "已操作次数")
			f.SetCellValue("Sheet1", "W1", "未操作次数")
			f.SetCellValue("Sheet1", "X1", "备注")
			var i int = 2

			for _, value := range order {
				log.Println("value:", value)
				// log.Println(key, "---", value)
				f.SetCellValue("Sheet1", "A"+strconv.Itoa(i), value.Orderid)
				f.SetCellValue("Sheet1", "B"+strconv.Itoa(i), strings.Split(value.Visittime, "T")[0])
				f.SetCellValue("Sheet1", "C"+strconv.Itoa(i), value.Customerid)
				f.SetCellValue("Sheet1", "D"+strconv.Itoa(i), value.Consumetype)
				f.SetCellValue("Sheet1", "E"+strconv.Itoa(i), value.Item)
				f.SetCellValue("Sheet1", "F"+strconv.Itoa(i), value.Treatnum)
				f.SetCellValue("Sheet1", "G"+strconv.Itoa(i), value.Swmdmoney)
				f.SetCellValue("Sheet1", "H"+strconv.Itoa(i), value.Qbmoney)
				f.SetCellValue("Sheet1", "I"+strconv.Itoa(i), value.Mmmoney)
				f.SetCellValue("Sheet1", "J"+strconv.Itoa(i), value.Mmboxnum)
				f.SetCellValue("Sheet1", "K"+strconv.Itoa(i), value.Donateboxnum)
				f.SetCellValue("Sheet1", "L"+strconv.Itoa(i), value.Dgfmoney)
				f.SetCellValue("Sheet1", "M"+strconv.Itoa(i), value.Xfymoney)
				f.SetCellValue("Sheet1", "N"+strconv.Itoa(i), value.Totalmoney)
				f.SetCellValue("Sheet1", "O"+strconv.Itoa(i), value.Returnmoney)
				f.SetCellValue("Sheet1", "P"+strconv.Itoa(i), value.Kkmoney)
				f.SetCellValue("Sheet1", "Q"+strconv.Itoa(i), value.Owedmoney)
				f.SetCellValue("Sheet1", "R"+strconv.Itoa(i), value.Paidmoney)
				f.SetCellValue("Sheet1", "S"+strconv.Itoa(i), value.Consultteach)
				f.SetCellValue("Sheet1", "T"+strconv.Itoa(i), value.Dqteach)
				f.SetCellValue("Sheet1", "U"+strconv.Itoa(i), value.Operateach)
				f.SetCellValue("Sheet1", "V"+strconv.Itoa(i), value.Operanum)
				f.SetCellValue("Sheet1", "W"+strconv.Itoa(i), value.Unoperanum)
				f.SetCellValue("Sheet1", "X"+strconv.Itoa(i), value.Comment)
				i = i + 1

			}
			log.Println("i:", i)
			err = f.SaveAs(filePath)
			if err != nil {
				log.Println(err)
				return "订单写入excel失败", err
			}
			// excelReader, _ := os.Open(filePath)
			// fileInfo, _ := excelReader.Stat()
			return filePath, nil
		}
	}

}
