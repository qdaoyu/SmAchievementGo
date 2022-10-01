package base

import (
	"log"
	"net/http"
	"qiudaoyu/models/customer"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// 导出尚美订单信息-待改
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
	filePath, err := customer.ExportSmCustomerList(userID, userNameAssert)
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
		c.Header("Content-Disposition", "attachment; filename="+userNameAssert+"customerInfo.xlsx") // 用来指定下载下来的文件名
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

// 更新尚美订单信息-待改
func UpdateSmOrderHandler(c *gin.Context) {
	var customertb customer.Customertb
	// log.Println("log:", c.Query("customerid"))
	err := c.ShouldBind(&customertb)
	log.Println("addCustomer:", customertb)
	log.Println(err)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    2001,
			"message": "会员信息更新失败,请联系系统管理员(江昌杰)" + err.Error(),
		})
		return
	}

	err = customer.UpdateSmCustomer(customertb)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    2002,
			"message": "会员信息更新失败,请联系系统管理员(江昌杰)" + err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "会员信息更新成功",
		})
	}
}

// 删除尚美订单信息-待改
func DeleteSmOrderHandler(c *gin.Context) {
	phone := c.Param("id")
	log.Println("log:", phone)
	if len(phone) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    2003,
			"message": "删除失败，无删除参数",
		})
		return
	}

	err := customer.DeleteSmCustomer(phone)
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

// 添加尚美订单信息-待改
func AddSmOrderHandler(c *gin.Context) {
	var errMsg string
	var customertb customer.Customertb
	log.Println("log:", c.Query("customerid"))
	err := c.ShouldBind(&customertb)
	log.Println("addCustomer:", customertb)
	log.Println(err)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    2001,
			"message": "会员信息新增失败,请联系系统管理员(江昌杰)" + err.Error(),
		})
		return
	}

	err = customer.AddSmCustomer(customertb)
	if err != nil {
		ok := strings.Contains(err.Error(), "PRIMARY")
		if ok {
			errMsg = "添加失败,此会员已注册！"
		} else {
			errMsg = "新增会员失败:" + err.Error()
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    2002,
			"message": errMsg,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "新增会员成功",
		})
	}
}

// 获取尚美订单信息列表(管理员_id1和测试角色_id2默认可以返回所有数据)-待改
func GetSmOrderListHandler(c *gin.Context) {
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

	item := c.Query("item")
	consumetype := c.Query("consumetype")
	consultteach := c.Query("consultteach")
	unoperanum, err := strconv.Atoi(c.Query("unoperanum"))
	if err != nil {
		log.Println("剩余未操作次数转化int失败")
		unoperanum = -1
	}
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
	smMap, err = customer.GetSmOrderList(userID, userNameAssert, currentPage, size, order)
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
