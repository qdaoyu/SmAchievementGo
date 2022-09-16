package menuInfo

import (
	"errors"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// 用户信息
type User struct {
	// gorm.Model         //自动添加主键
	ID        int    `gorm:"type:int(0) " json:"id" binding:"required"`
	Name      string `gorm:"type:varchar(64)" json:"name" binding:"required"`
	Phone     string `gorm:"type:varchar(64)" json:"phone" binding:"-"`
	TelePhone string `gorm:"type:varchar(64) " json:"telephone" binding:"-"`
	Address   string `gorm:"type:varchar(64) " json:"address" binding:"-"`
	Enabled   string `gorm:"type:varchar(64) " json:"enabled" binding:"required"`
	Username  string `gorm:"type:int(0) " json:"username" binding:"required"`
	Password  string `gorm:"type:int(0) " json:"password" binding:"required"`
	Userface  string `gorm:"type:varchar(255) " json:"userface" binding:"-"`
	Remark    string `gorm:"type:varchar(255) " json:"remark" binding:"-"`
}

// 定义数据地址
var Dsn string = "root:qiudaoyu@tcp(127.0.0.1:3306)/qiudaoyu?charset=utf8mb4&parseTime=True&loc=Local"

func DbCommonOperation(dsn string) (*gorm.DB, error) {
	//连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic("failed to connect database")
	}
	sqlDB, _ := db.DB()

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	return db, err
}

func LoginConfirm(userName string, passWord string) (gin.H, error) {

	//连接数据库
	var user User
	db, err := DbCommonOperation(Dsn)
	if err != nil {
		panic("failed to connect database")
	}

	// AND
	db.Table("t_admin").Where("username = ? AND password = ?", userName, passWord).Find(&user)
	if !reflect.DeepEqual(user, User{}) {
		return gin.H{"userInfo": user}, nil
	} else {
		return gin.H{}, errors.New("账号密码错误")
	}

}

func GetUserInfo(c *gin.Context) (gin.H, error) {
	//从token中获取用户名username
	// dsn := "root:qiudaoyu@tcp(127.0.0.1:3306)/qiudaoyu?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := DbCommonOperation(Dsn)

	//获取个人信息
	var user User
	username, exist := c.Get("username")
	db.Table("t_admin").Where("username = ? ", username).Find(&user)
	// fmt.Println(t_menu)
	if !exist {
		return gin.H{"message": "用户信息获取失败"}, errors.New("用户信息获取失败")
	} else {
		return gin.H{
			"data":    user,
			"message": "用户信息获取成功",
		}, nil
	}

}
