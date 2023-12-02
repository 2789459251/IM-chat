package service

import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

func SearchFriends(c *gin.Context) {
	con := models.Contact{}
	idstr := c.Request.FormValue("userId")
	if idstr == "" {
		// 处理没有提供 userId 的情况，可能返回错误或其他逻辑
		fmt.Println("userId parameter is missing")
		return
	}
	id, err := strconv.Atoi(idstr)
	if err != nil {
		// 处理错误，可能是因为参数不是有效的整数
		fmt.Println("Error converting userId:", err)
		// 返回错误或其他逻辑
		return
	}
	users := con.SearchFriends(uint(id))
	utils.RespOKList(c.Writer, users, len(users))
}
