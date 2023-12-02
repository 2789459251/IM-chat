package models

import (
	"fmt"
	"ginchat/utils"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Name          string
	Passwd        string
	Email         string `valid:"email"`
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9})"`
	Salt          string
	Identity      string
	ClentIp       string
	ClentPort     string
	LoginTime     time.Time `gorm:"default:NULL"`
	HeartbeatTime time.Time `gorm:"default:NULL"`
	LoginOutTime  time.Time `gorm:"column:login_out_time;default:NULL" json:"login_out_time"`
	IsLogout      bool
	DeviceInfo    string
}

func (table *User) TableName() string {
	return "uer_basic"
}
func GetUserList() []*User {
	date := make([]*User, 10)
	utils.DB.Find(&date)
	for i, user := range date {
		fmt.Println(i, user)
	}
	return date
}
func FindUserByNameAndPasswd(Name, Passwd string) User {
	user := User{}
	//utils.DB.Where("name=?and passwd=?", Name, Passwd).First(&user)
	utils.DB.Where("name = ? AND passwd = ? "+
		"AND `uer_basic`.`deleted_at` IS NULL", Name, Passwd).
		Order("id").Limit(1).Find(&user)

	//token加密
	str := fmt.Sprintf("%d", time.Now().Unix())
	temp := utils.MD5Encode(str)
	utils.DB.Model(&User{}).
		Where("id = ?", user.ID).Update("identity", temp)
	return user

}
func FindUserByName(Name string) User {
	user := User{}
	utils.DB.Where("name=?", Name).First(&user)
	return user
}
func FindUserByPhone(Phone string) User {
	user := User{}
	utils.DB.Where("phone=?", Phone).First(&user)
	return user
}
func FindUserByEmail(Email string) User {
	user := User{}
	utils.DB.Where("Email=?", Email).First(&user)
	return user
}
func CreateUser(user User) *gorm.DB {
	return utils.DB.Create(&user)
}
func DeleteUser(user User) *gorm.DB {
	return utils.DB.Delete(&user)
}
func UpdateUser(user User) *gorm.DB {
	return utils.DB.Model(&user).Updates(User{Name: user.Name, Passwd: user.Passwd, Phone: user.Phone, Email: user.Email})
}
