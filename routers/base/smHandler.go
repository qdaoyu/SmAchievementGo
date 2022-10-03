package base

import (
	"log"
	"net/http"
	"qiudaoyu/models/customer"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// 导出尚美订单信息
func ExportSmOrderHandler(c *gin.Context) {
	log.Println("进入")
	userID, _ := strconv.Atoi(c.Request.Header.Get("userID"))
	//类型断言
	userName, _ := c.Get("username")
	userNameAssert, ok := userName.(string)
	log.Println(userNameAssert)
	// log.Println("测试:", userID, userName)
	// userName, _ := c.Get("username")
	if !ok {
		c.JSON(200, gin.H{
			"code":    5003,
			"message": "用户名断言失败",
			"data":    nil,
		})
		return
	}

	//返回文件地址
	filePath, err := customer.ExportSmOrderList(userID, userNameAssert)
	log.Println(filePath)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    5002,
			"message": "下载失败:" + err.Error(),
			"data":    nil,
		})
		return
	} else {

		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment; filename="+userNameAssert+"OrderInfo.xlsx") // 用来指定下载下来的文件名
		c.Header("Content-Transfer-Encoding", "binary")

		// fileData, err := ioutil.ReadFile(filePath)
		if err != nil {
			c.JSON(200, gin.H{
				"code":    5002,
				"message": "下载失败:" + err.Error(),
				"data":    nil,
			})
			return
		}
		c.File(filePath)
		// c.Data(200, "application/octet-stream", fileData)
		return
	}
}

// 更新尚美订单信息
func UpdateSmOrderHandler(c *gin.Context) {
	var ordertb customer.Ordertb
	// log.Println("log:", c.Query("customerid"))
	err := c.ShouldBind(&ordertb)
	log.Println("addCustomer:", ordertb)
	log.Println(err)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    2001,
			"message": "订单信息更新失败,请联系系统管理员(江昌杰)" + err.Error(),
		})
		return
	}

	err = customer.UpdateSmOrder(ordertb)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    2002,
			"message": "订单信息更新失败,请联系系统管理员(江昌杰)" + err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "订单信息更新成功",
		})
	}
}

// 删除尚美订单信息-待改
func DeleteSmOrderHandler(c *gin.Context) {
	orderID := c.Param("id")
	log.Println("log:", orderID)
	if len(orderID) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    2003,
			"message": "删除失败，无删除参数",
		})
		return
	}

	err := customer.DeleteSmOrder(orderID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    2002,
			"message": "删除失败:" + err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "删除成功",
		})
	}
}

// 添加尚美订单信息
func AddSmOrderHandler(c *gin.Context) {
	var errMsg string
	var ordertb customer.Ordertb
	err := c.ShouldBind(&ordertb)
	log.Println("addOrder:", ordertb)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    2001,
			"message": "订单信息新增失败,请联系系统管理员(江昌杰)" + err.Error(),
		})
		return
	}

	err = customer.AddSmOrder(ordertb)
	if err != nil {
		log.Println("err.Error():-----------", err.Error())
		ok := strings.Contains(err.Error(), "Duplicate")
		if ok {
			errMsg = "添加失败,该顾客当天已登记见诊记录！"
		} else {
			log.Println("err.Error():-----------", err.Error())
			errMsg = "添加订单失败:" + err.Error()
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    2002,
			"message": errMsg,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "添加订单记录成功",
		})
	}
}

// 获取尚美订单信息列表(管理员_id1和测试角色_id2默认可以返回所有数据)
func GetSmOrderListHandler(c *gin.Context) {

	// log.Println("当前日期", time.Now().Format("2006-01-02"))
	var order customer.Ordertb
	var smMap = make(map[string]interface{})
	currentPage, err := strconv.Atoi(c.Query("currentPage"))
	if err != nil {
		log.Println(err)
		return
	}
	size, err := strconv.Atoi(c.Query("size"))
	if err != nil {
		log.Println(err)
		return
	}
	name := c.Query("name")
	item := c.Query("item")
	consumetype := c.Query("consumetype")
	consultteach := c.Query("consultteach")
	unoperanum, err := strconv.Atoi(c.Query("unoperanum"))
	if err != nil {
		log.Println("剩余未操作次数转化int失败")
		unoperanum = -1
	}
	// order.Name = name
	order.Item = item
	order.Consumetype = consumetype
	order.Consultteach = consultteach
	order.Unoperanum = unoperanum
	log.Println("order:", order)
	userID, _ := strconv.Atoi(c.Request.Header.Get("userID"))
	//类型断言
	userName, _ := c.Get("username")
	userNameAssert, ok := userName.(string)
	// log.Println("测试:", userID, userName)
	// userName, _ := c.Get("username")
	if !ok {
		c.JSON(200, gin.H{
			"code":    5003,
			"message": "用户名断言失败",
			"data":    nil,
		})
		return
	}
	smMap, err = customer.GetSmOrderList(userID, userNameAssert, currentPage, size, order, name)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    5002,
			"message": smMap["message"],
			"data":    nil,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": smMap["message"],
			"data":    smMap["data"],
			"total":   smMap["total"],
		})
		return
	}
}
