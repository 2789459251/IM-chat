package models

// 人员关系
import (
	"ginchat/utils"
	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	OwnerId  int64 //谁的关系
	TargetId uint  //对应的谁
	Type     int   //对应的关系类型 1好友 2群聊 3
	Desc     string
}

func (table *Contact) TableName() string {
	return "contact"
}

// 查找某人的好友
func (c Contact) SearchFriends(userId uint) []User {
	contacts := make([]Contact, 0)
	objIds := make([]uint64, 0)
	utils.DB.Where("owner_id=? and type=1", userId).Find(&contacts)
	for _, v := range contacts {
		objIds = append(objIds, uint64(v.TargetId))
	}
	users := make([]User, 0)
	utils.DB.Where("id IN ?", objIds).Find(&users)
	return users
}
