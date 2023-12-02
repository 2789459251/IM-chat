package models

import "gorm.io/gorm"

// 群信息
type Group struct {
	gorm.Model
	Name    string
	OwnerId uint
	Icon    string
	Type    int //lev
	Desc    string
}

func (table *Group) TableName() string {
	return "group_basic"
}
