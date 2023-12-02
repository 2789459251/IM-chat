package service

import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// GetUserList
// @Tags  用户模块
// @Summary 用户列表
// @Seccess 200 {string} json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	date := make([]*models.User, 10)
	date = models.GetUserList()
	c.JSON(http.StatusOK, gin.H{
		"code":    0, //0成功-1失败
		"message": "获取成功",
		"data":    date,
	})
}

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Seccess 200 {string} json{"code","message"}
// @Router /user/createUser [post]
func CreateUser(c *gin.Context) {
	user := models.User{}
	user.Name = c.Request.FormValue("name")

	passwd := c.Request.FormValue("password")
	repasswd := c.Request.FormValue("Identity")
	salt := fmt.Sprintf("%06d", rand.Int31())
	date := models.FindUserByName(user.Name)
	if user.Name == "" || passwd == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, //0成功-1失败
			"message": "用户或密码不能为空",
			"data":    user,
		})
		return
	}
	if date.Name != "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, //0成功-1失败
			"message": "用户名已注册",
			"data":    user,
		})
		return
	}
	fmt.Println(passwd)
	fmt.Println(repasswd)
	if passwd != repasswd {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, //0成功-1失败
			"message": "两次密码不一致",
			"data":    user,
		})
	} else {
		//user.Passwd = passwd
		user.Passwd = utils.MakePasswd(passwd, salt)
		user.Salt = salt
		models.CreateUser(user)
		c.JSON(http.StatusOK, gin.H{
			"code":    0, //0成功-1失败
			"message": "新增用户成功",
			"data":    user,
		})
	}

}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id query string false "id"
// @Seccess 200 {string} json{"code","message"}
// @Router /user/deleteUser [delete]
func DeleteUser(c *gin.Context) {
	user := models.User{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)

	models.DeleteUser(user)

	c.JSON(http.StatusOK, gin.H{
		"code":    0, //0成功-1失败
		"message": "删除用户成功",
		"data":    user,
	})
}

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @param id query string false "id"
// @param name query string false "name"
// @param passwd query string false "passwd"
// @param phone query string false "phone"
// @param email query string false "email"
// @Seccess 200 {string} json{"code","message"}
// @Router /user/updateUser [put]
func UpdateUser(c *gin.Context) {
	user := models.User{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	user.Name = c.Query("name")
	user.Passwd = c.Query("passwd")
	user.Phone = c.Query("phone")
	user.Email = c.Query("email")

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, //0成功-1失败
			"message": "修改参数不匹配",
			"data":    user,
		})
	} else {
		models.UpdateUser(user)
		c.JSON(http.StatusOK, gin.H{
			"code":    0, //0成功-1失败
			"message": "修改用户成功",
			"data":    user,
		})
	}
}

// FindUserByNameAndPasswd
// @Tags  用户模块
// @Summary 用户登录
// @param name query string false "name"
// @param passwd query string false "password"
// @Seccess 200 {string} json{"code","message"}
// @Router /user/findUserByNameAndPwd [post]
func FindUserByNameAndPasswd(c *gin.Context) {
	name := c.Request.FormValue("name")
	passwd := c.Request.FormValue("password")
	//fmt.Println("name:", name, "passwd", passwd)
	user2 := models.FindUserByName(name)
	if user2.Name == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, //0成功-1失败
			"message": "用户名不存在",
			"data":    nil,
		})
	}
	flag := utils.ValidPassword(passwd, user2.Salt, user2.Passwd)
	if !flag {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1, //0成功-1失败
			"message": "密码错误",
			"data":    nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0, //0成功-1失败
		"message": "欢迎回来" + user2.Name,
		"data":    user2,
	})

}

// 防止跨域伪装请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(ws)
	MsgHandler(ws, c)
}
func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	for {
		msg, err := utils.Subscribe(c, utils.PublishKey)
		if err != nil {
			fmt.Println(err)
			return
		}
		tm := time.Now().Format("2006-01-02 15:04:05")
		msg = fmt.Sprintf("[ws][%s]:%s", tm, msg)
		err = ws.WriteMessage(1, []byte(msg))
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}
func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
