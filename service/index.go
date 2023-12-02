package service

import (
	"fmt"
	"ginchat/models"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"strconv"
)

// GetIndex
// @Tags 首页
// @Seccess 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	//ind, err := template.ParseFiles("index.html")
	//
	ind, err := template.ParseFiles("index.html", "views/chat/head.html")
	if err != nil {
		panic(err)
	}

	fmt.Println("Template parsed successfully")
	err = ind.Execute(c.Writer, "index")
	if err != nil {
		panic(err)
	}

	fmt.Println("Template executed successfully")
	//c.JSON(http.StatusOK, gin.H{
	//	"hhh": "welcome!",
	//})
}
func ToRegister(c *gin.Context) {
	if c.Writer.Status() == http.StatusOK {
		//ind, err := template.ParseFiles("index.html", "views/chat/head.html",
		//	"views/user/register.html"){{template "/chat/tabmenu.shtml"}}
		ind, err := template.ParseFiles("views/user/register.html")
		if err != nil {
			panic(err)
		}
		err = ind.Execute(c.Writer, "register")
		if err != nil {
			panic(err)
		}
	} else {
		// 重定向到 register 页面
		c.Redirect(http.StatusFound, "/user/register")
	}
}
func ToChat(c *gin.Context) {
	ind, err := template.ParseFiles("views/chat/index.html",
		"views/chat/group.html", "views/chat/concat.html", "views/chat/profile.html",
		"views/chat/main.html", "views/chat/createcom.html", "views/chat/userinfo.html",
		"views/chat/foot.html", "views/chat/tabmenu.html", "views/chat/head.html")
	if err != nil {
		panic(err)
	}
	userId, _ := strconv.Atoi(c.Query("userId"))
	token := c.Query("token")
	user := models.User{}
	user.ID = uint(userId)
	user.Identity = token
	//校验token
	err = ind.Execute(c.Writer, user)
	if err != nil {
		panic(err)
	}

}
func Chat(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
