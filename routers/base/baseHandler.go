package base

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"qiudaoyu/middleWare"
	"qiudaoyu/models/achieve"
	"qiudaoyu/models/basic"
	"qiudaoyu/models/customer"
	"qiudaoyu/models/menuInfo"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Data struct {
	Situ01 []string
	Situ02 []string
	Situ03 []string
	Situ04 []string
	Situ05 []string
}

// ----token部分----
type UserInfo struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// 上传业绩表
func UploadSyAchieveTb(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println(err)
	} else {
		log.Println(file.Filename)

		dst := "./assets/fileRec/" + file.Filename
		// 上传文件至指定的完整文件路径
		c.SaveUploadedFile(file, dst)

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	}
}

// 登录验证
func Login(c *gin.Context) {
	// username := c.PostForm("username")
	// password := c.PostForm("password")
	var user UserInfo
	err := c.ShouldBind(&user)
	fmt.Println("userInfo:", user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    2001,
			"message": "无效的参数",
		})
		return
	}
	// 校验用户名和密码是否正确,数据库取用户信息
	//调用数据库查询，根据返回值，判定是否为t_admin表里的用户
	res, err2 := menuInfo.LoginConfirm(user.Username, user.Password)
	if err2 == nil {
		// 生成Token
		fmt.Println("res:", res)
		// var userStruct menuInfo.User
		// err3 := json.Unmarshal(res, &userStruct)
		// resByte, _ := reflect.TypeOf(res["userInfo"]).FieldByName("ID")
		// fmt.Println("res:", reflect.TypeOf(res))
		// fmt.Println(reflect.TypeOf(res["userInfo"]).FieldByName("ID"))
		// fmt.Println(reflect.ValueOf(res["userInfo"]).FieldByName("ID").Int())  res["userInfo"].ID
		fmt.Println(res["userInfo"])
		// res.
		// fmt.Println(res["userInfo"].ID)

		//类型断言
		userId, ok := res["userInfo"].(menuInfo.User)

		if ok {
			fmt.Println(userId.ID)
			// tokenString, _ := middleWare.GenToken(user.Username, reflect.ValueOf(res["userInfo"]).FieldByName("ID").Int())
			tokenString, _ := middleWare.GenToken(user.Username, int64(userId.ID))
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "登录成功",
				"data":    gin.H{"token": tokenString, "userInfo": res["userInfo"]},
			})
			return

		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":    2002,
				"message": "用户名或密码错误",
			})
		}
		// tokenString, _ := middleWare.GenToken(user.Username, reflect.ValueOf(res["userInfo"]).FieldByName("ID").Int())

		return
	}
	log.Println("用户名或密码错误")
	c.JSON(http.StatusOK, gin.H{
		"code":    2002,
		"message": "用户名或密码错误",
	})
}

// 获取菜单信息
func HomeMenuHandler(c *gin.Context) {
	// username := c.PostForm("username")
	// password := c.PostForm("password")
	userID, err1 := strconv.Atoi(c.Request.Header.Get("userID"))
	if err1 != nil {
		c.JSON(5001, gin.H{
			"code":    5001,
			"message": "获取菜单失败",
			"data":    gin.H{},
		})
		return
	}
	res, err := menuInfo.GetMenuDb(userID)
	fmt.Println(res)
	if err != nil {
		c.JSON(5001, gin.H{
			"code":    5001,
			"message": "获取菜单失败",
			"data":    gin.H{},
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取菜单成功",
			"data":    res,
		})
		return
	}
}

// 获取用户信息
func GetUserInfoHandler(c *gin.Context) {
	//获取menu菜单
	res, err := menuInfo.GetUserInfo(c)
	fmt.Println(res)
	if err != nil {
		c.JSON(5001, gin.H{
			"code":    5001,
			"message": "获取用户信息失败",
			"data":    gin.H{},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取用户信息成功",
			"data":    res,
		})
		return
	}
}

// 用户资料信息
func GetUserBasicInfoHandler(c *gin.Context) {
	//获取menu菜单
	userID, _ := strconv.Atoi(c.Request.Header.Get("userID"))
	res, err := menuInfo.GetUserBasicInfo(c, userID)
	fmt.Println(res)
	if err != nil {
		c.JSON(5001, gin.H{
			"code":    5001,
			"message": "获取用户基本信息失败",
			"data":    nil,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取用户基本信息成功",
			"data":    res,
		})
		return
	}
}

// 塑颜业绩新增（导入）
func AddSyAchieveInfoHandler(c *gin.Context) {

	// userID, _ := strconv.Atoi(c.Request.Header.Get("userID"))
	err := achieve.SyAchieveExcelize(c)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    5001,
			"message": "插入塑颜业绩数据库失败,原因:" + err.Error(),
			"data":    nil,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "插入塑颜业绩数据库成功",
		})
		return
	}
}

// 获取塑颜业绩表信息(管理员_id1和测试角色_id2默认可以返回所有数据)
// func GetSyAchieveInfoHandler(c *gin.Context) {
// 	var syMap = make(map[string]interface{})
// 	currentPage, err := strconv.Atoi(c.Query("currentPage"))
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	size, err := strconv.Atoi(c.Query("size"))
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	userID, _ := strconv.Atoi(c.Request.Header.Get("userID"))
// 	//类型断言
// 	userName, _ := c.Get("username")
// 	userNameAssert, ok := userName.(string)
// 	// log.Println("测试:", userID, userName)
// 	// userName, _ := c.Get("username")
// 	if !ok {
// 		c.JSON(200, gin.H{
// 			"code":    5003,
// 			"message": "用户名断言失败",
// 			"data":    nil,
// 		})
// 		return
// 	}
// 	syMap, err = achieve.GetSyAchieve(userID, userNameAssert, currentPage, size)
// 	if err != nil {
// 		c.JSON(200, gin.H{
// 			"code":    5002,
// 			"message": syMap["message"],
// 			"data":    nil,
// 		})
// 		return
// 	} else {
// 		c.JSON(http.StatusOK, gin.H{
// 			"code":    200,
// 			"message": syMap["message"],
// 			"data":    syMap["data"],
// 			"total":   syMap["total"],
// 		})
// 		return
// 	}
// }

// 导出尚美会员信息
func ExportSmCustomerHandler(c *gin.Context) {
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

// 更新尚美会员信息
func UpdateSmCustomerHandler(c *gin.Context) {
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

// 删除尚美会员信息 deleteSmCustomerHandler
func DeleteSmCustomerHandler(c *gin.Context) {
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

// 添加尚美会员信息
func AddSmCustomerHandler(c *gin.Context) {
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

// 获取项目信息
func GetItemHandler(c *gin.Context) {
	// var smMap = make(map[string]interface{})
	resMap, err := basic.GetSmItem()
	if err != nil {
		c.JSON(200, gin.H{
			"code":    5002,
			"message": resMap["message"],
			"data":    nil,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": resMap["message"],
			"data":    resMap["data"],
			"total":   resMap["total"],
		})
		return
	}
}

// GetConsumetypeHandler
// 获取消费类型信息
func GetConsumetypeHandler(c *gin.Context) {
	// var smMap = make(map[string]interface{})
	resMap, err := basic.GetSmConsumetype()
	if err != nil {
		c.JSON(200, gin.H{
			"code":    5002,
			"message": resMap["message"],
			"data":    nil,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": resMap["message"],
			"data":    resMap["data"],
			"total":   resMap["total"],
		})
		return
	}
}

// 获取咨询师信息
func GetConsultteachHandler(c *gin.Context) {
	// var smMap = make(map[string]interface{})
	resMap, err := basic.GetSmConsultteach()
	if err != nil {
		c.JSON(200, gin.H{
			"code":    5002,
			"message": resMap["message"],
			"data":    nil,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": resMap["message"],
			"data":    resMap["data"],
			"total":   resMap["total"],
		})
		return
	}
}

// 获取门店信息
func GetShopHandler(c *gin.Context) {
	// var smMap = make(map[string]interface{})
	resMap, err := basic.GetSmShop()
	if err != nil {
		c.JSON(200, gin.H{
			"code":    5002,
			"message": resMap["message"],
			"data":    nil,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": resMap["message"],
			"data":    resMap["data"],
			"total":   resMap["total"],
		})
		return
	}
}

// 获取尚美会员信息列表(管理员_id1和测试角色_id2默认可以返回所有数据)
func GetSmCustomerListHandler(c *gin.Context) {
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
	smMap, err = customer.GetSmCustomerList(userID, userNameAssert, currentPage, size)
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

// 重名计算
func CalDupName(c *gin.Context) {
	// 注意：下面为了举例子方便，暂时忽略了错误处理
	b, err := c.GetRawData() // 从c.Request.Body读取请求数据
	fmt.Println(b)
	if err == nil {
		// 定义map或结构体
		var m map[string]string
		// 反序列化
		// fmt.Printf()
		_ = json.Unmarshal(b, &m)
		arg2 := m["textArea"]
		arg1 := m["label"]
		fmt.Println("----以下为接收到的参数-----")
		fmt.Println(arg1)
		fmt.Println(arg2)
		// fmt.Println("---------")
		//执行python脚本-------
		cmd := exec.Command("python", "D:/tempdownload/重名判断.py", arg1, arg2)
		out, _ := cmd.CombinedOutput()
		fmt.Println("concatenation: ", string(out))
		msg := string(out)
		var data Data
		_ = json.Unmarshal([]byte(msg), &data)
		fmt.Println(string(out))
		//-----------------------
		// fmt.Println(m)
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "计算成功",
			"data":    data,
		})

	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    50001,
			"message": "输入内容非法",
			"data":    gin.H{},
		})
	}

}

func Page404(c *gin.Context) {
	// c.HTML(http.StatusNotFound, "views/404.html", nil)
	c.Redirect(http.StatusMovedPermanently, "http://localhost:8080/page404")
}

// func AuthHandler(c *gin.Context) {
// 	// 用户发送用户名和密码过来
// 	var user UserInfo
// 	err := c.ShouldBind(&user)
// 	if err != nil {
// 		c.JSON(http.StatusOK, gin.H{
// 			"code": 2001,
// 			"msg":  "无效的参数",
// 		})
// 		return
// 	}
// 	// 校验用户名和密码是否正确
// 	if user.Username == "q1mi" && user.Password == "q1mi123" {
// 		// 生成Token
// 		tokenString, _ := middleWare.GenToken(user.Username)
// 		c.JSON(http.StatusOK, gin.H{
// 			"code": 2000,
// 			"msg":  "success",
// 			"data": gin.H{"token": tokenString},
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"code": 2002,
// 		"msg":  "鉴权失败",
// 	})
// 	return
// }

//---------------
