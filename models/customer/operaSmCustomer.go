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
	 , shop = ? , consultteach = ? , visittime = ? , phone = ? where customerid = ?  `
	log.Println("customer:", customer)
	// var err error = models.Conn.Table("t_customer").Where(`customerid = ? and  name = ? and gender = ?
	// 		and phone = ? and shop = ? and consultteach = ? and visittime = ? `, customer.Customerid, customer.Name,
	// 	customer.Gender, customer.Phone, customer.Shop, customer.Consultteach, customer.Visittime).Updates(&Customertb{}).Error
	err := models.Conn.Exec(sqlString, customer.Name, customer.Gender, customer.Shop, customer.Consultteach, customer.Visittime, customer.Phone, customer.Customerid).Error
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

// 导出会员信息
// 获取尚美会员信息表,返回文件地址
func ExportSmCustomerList(uid int, username string) (string, error) {
	var customer []Customer
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
		sqlString := `select t_customer.customerid ,t_customer.name ,t_customer.gender ,
		t_customer.phone ,t_customer.shop,t_customer.consultteach , date_format(t_customer.visittime,"%Y-%m-%d") as visittime,
		t_orderlist.treatnum ,t_orderlist.operanum ,t_orderlist.unoperanum  
		from t_customer left join t_orderlist on t_customer.customerid = t_orderlist.customerid `

		// err := models.Conn.Raw(sqlString).Count(&total).Limit(size).Offset(offsetVal).Scan(&customer).Error
		err := models.Conn.Raw(sqlString).Scan(&customer).Error
		if err != nil {
			log.Println(err)
			return "查询失败:" + err.Error(), err
		}
		// models.Conn.Raw(sqlString, userID).Scan(&user)
		// db.Raw(sqlString, userID).Create(&user)
		if len(customer) == 0 {
			log.Println("数据库无数据")
			smMap["message"] = "数据库无数据"
			return "", errors.New("数据库无数据")
		} else {
			var filePath string = "./assets/exportFile/" + username + strconv.FormatInt(time.Now().Unix(), 10) + ".xlsx"
			//数据库数据整理到表格中
			f := excelize.NewFile()
			// 设置单元格的值
			f.SetCellValue("Sheet1", "A1", "会员编号")
			f.SetCellValue("Sheet1", "B1", "姓名")
			f.SetCellValue("Sheet1", "C1", "性别")
			f.SetCellValue("Sheet1", "D1", "手机号码")
			f.SetCellValue("Sheet1", "E1", "门店")
			f.SetCellValue("Sheet1", "F1", "咨询师")
			f.SetCellValue("Sheet1", "G1", "见诊日期")
			// f.SetCellValue("Sheet1", "H1", "见诊日期")
			// f.SetCellValue("Sheet1", "I1", "见诊日期")
			//数据库取值
			sqlStringTotal := `select * from t_customer`
			err := models.Conn.Raw(sqlStringTotal).Scan(&customer).Error
			if err != nil {
				log.Println(err)
				return "", err
			}
			var i int = 2
			for key, value := range customer {
				log.Println(key, "---", value)
				f.SetCellValue("Sheet1", "A"+strconv.Itoa(i), value.Customerid)
				f.SetCellValue("Sheet1", "B"+strconv.Itoa(i), value.Name)
				f.SetCellValue("Sheet1", "C"+strconv.Itoa(i), value.Gender)
				f.SetCellValue("Sheet1", "D"+strconv.Itoa(i), value.Phone)
				f.SetCellValue("Sheet1", "E"+strconv.Itoa(i), value.Shop)
				f.SetCellValue("Sheet1", "F"+strconv.Itoa(i), value.Consultteach)
				f.SetCellValue("Sheet1", "G"+strconv.Itoa(i), strings.Split(value.Visittime, "T")[0])
				i = i + 1

			}
			err = f.SaveAs(filePath)
			if err != nil {
				log.Println(err)
				return "会员写入excel失败", err
			}
			// excelReader, _ := os.Open(filePath)
			// fileInfo, _ := excelReader.Stat()
			return filePath, nil
		}
	} else {
		//limit offset
		sqlString := `select t_customer.customerid ,t_customer.name ,t_customer.gender ,
		t_customer.phone ,t_customer.shop,t_customer.consultteach , date_format(t_customer.visittime,"%Y-%m-%d") as visittime ,
		t_orderlist.treatnum,t_orderlist.operanum ,t_orderlist.unoperanum 
		from t_customer left join t_orderlist on t_customer.customerid = t_orderlist.customerid where t_customer.consultteach = ? `

		err := models.Conn.Raw(sqlString, username).Scan(&customer).Error
		if err != nil {
			log.Println(err)
			return "下载会员信息失败", err
		}
		// models.Conn.Raw(sqlString, userID).Scan(&user)
		// db.Raw(sqlString, userID).Create(&user)
		if len(customer) == 0 {
			return "无数据,无法下载", errors.New("无数据,无法下载")
		} else {
			var filePath string = "./assets/exportFile/" + username + strconv.FormatInt(time.Now().Unix(), 10) + ".xlsx"
			//数据库数据整理到表格中
			f := excelize.NewFile()
			// 设置单元格的值
			f.SetCellValue("Sheet1", "A1", "会员编号")
			f.SetCellValue("Sheet1", "B1", "姓名")
			f.SetCellValue("Sheet1", "C1", "性别")
			f.SetCellValue("Sheet1", "D1", "手机号码")
			f.SetCellValue("Sheet1", "E1", "门店")
			f.SetCellValue("Sheet1", "F1", "咨询师")
			f.SetCellValue("Sheet1", "G1", "见诊日期")
			//数据库取值
			// sqlStringTotal := `select * from t_customer`
			// err := models.Conn.Raw(sqlStringTotal).Scan(&customer).Error
			// if err != nil {
			// 	log.Println(err)
			// 	return "", err
			// }
			var i int = 2
			for key, value := range customer {
				log.Println(key, "---", value)
				f.SetCellValue("Sheet1", "A"+strconv.Itoa(i), value.Customerid)
				f.SetCellValue("Sheet1", "B"+strconv.Itoa(i), value.Name)
				f.SetCellValue("Sheet1", "C"+strconv.Itoa(i), value.Gender)
				f.SetCellValue("Sheet1", "D"+strconv.Itoa(i), value.Phone)
				f.SetCellValue("Sheet1", "E"+strconv.Itoa(i), value.Shop)
				f.SetCellValue("Sheet1", "F"+strconv.Itoa(i), value.Consultteach)
				f.SetCellValue("Sheet1", "G"+strconv.Itoa(i), strings.Split(value.Visittime, "T")[0])
				i = i + 1

			}
			err = f.SaveAs(filePath)
			if err != nil {
				log.Println(err)
				return "会员写入excel失败", err
			}
			// excelReader, _ := os.Open(filePath)
			// fileInfo, _ := excelReader.Stat()
			return filePath, nil
		}
	}

}
