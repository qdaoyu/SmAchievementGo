package middleWare

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// gin jwt-go生成token及jwt中间件验证
// Jwtkey 秘钥 可通过配置文件配置
var Jwtkey = []byte("blog_jwt_key")

type MyClaims struct {
	UserId   int    `json:"user_id"`
	UserName string `json:"username"`
	jwt.StandardClaims
}

// CreateToken 生成token
func CreateToken(userId int, userName string) (string, error) {
	expireTime := time.Now().Add(2 * time.Hour) //过期时间
	nowTime := time.Now()                       //当前时间
	claims := MyClaims{
		UserId:   userId,
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间戳
			IssuedAt:  nowTime.Unix(),    //当前时间戳
			Issuer:    "blogLeo",         //颁发者签名
			Subject:   "userToken",       //签名主题
		},
	}
	tokenStruct := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenStruct.SignedString(Jwtkey)
}

// CheckToken 验证token
func CheckToken(token string) (*MyClaims, bool) {
	tokenObj, _ := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return Jwtkey, nil
	})
	if key, _ := tokenObj.Claims.(*MyClaims); tokenObj.Valid {
		return key, true
	} else {
		return nil, false
	}
}

// JwtMiddleware jwt中间件
func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//从请求头中获取token
		tokenStr := c.Request.Header.Get("Authorization")
		//用户不存在
		if tokenStr == "" {
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "用户不存在"})
			c.Abort() //阻止执行
			return
		}
		//token格式错误
		tokenSlice := strings.SplitN(tokenStr, " ", 2)
		if len(tokenSlice) != 2 && tokenSlice[0] != "Bearer" {
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "token格式错误"})
			c.Abort() //阻止执行
			return
		}
		//验证token
		tokenStruck, ok := CheckToken(tokenSlice[1])
		if !ok {
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "token不正确"})
			c.Abort() //阻止执行
			return
		}
		//token超时
		if time.Now().Unix() > tokenStruck.ExpiresAt {
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "token过期"})
			c.Abort() //阻止执行
			return
		}
		c.Set("username", tokenStruck.UserName)
		c.Set("user_id", tokenStruck.UserId)

		c.Next()
	}
}

//------------------------------------------

// StatCost 是一个统计耗时请求耗时的中间件：记录接口耗时的中间件
func StatCost() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Set("name1", "小王子") // 可以通过c.Set在请求上下文中设置值，后续的处理函数能够取到该值
		// 调用该请求的剩余处理程序
		c.Next()
		// 不调用该请求的剩余处理程序
		// c.Abort()
		// 计算耗时
		cost := time.Since(start)
		log.Println(cost)
	}
}

// 处理跨域
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
