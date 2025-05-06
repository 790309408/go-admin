package models

import (
	"go-admin/common/models"
	"time"
)

/*
json: 用于控制结构体字段在JSON序列化和反序列化时的行为
例如，如果你有一个字段不想序列化，可以使用json:"-"来忽略它。
如果字段的名称在JSON中与结构体中的字段名不同，可以通过json:"field_name"来指定JSON中的字段名
gorm:

	GORM是一个Go语言的ORM库，它通过标签来定义结构体字段与数据库表之间的映射关系
	例如，你可以使用gorm:"column:name"来指定数据库中的列名，
	或者使用gorm:"primaryKey"来指定主键

comment: 在使用GORM进行数据库迁移时，comment 标签用于为数据库表的字段添加注释
*/
type SysUserSalary struct {
	models.Model
	UserId            int       `json:"userId" gorm:"size:20;comment:用户ID"`
	Username          string    `json:"username" gorm:"size:64;comment:用户名"`
	DeptId            int       `json:"deptId" gorm:"size:20;comment:部门ID"`
	RoleId            int       `json:"roleId" gorm:"size:20;comment:角色ID"`
	BasicSalary       float32   `json:"basicSalary" gorm:"size:20;comment:基础工资"`
	PerformanceBonus  float32   `json:"performanceBonus" gorm:"size:20;comment:绩效奖金"`
	OvertimePay       float32   `json:"overtimePay" gorm:"size:20;comment:加班费"`
	SocialInsurance   float32   `json:"socialInsurance" gorm:"size:20;comment:社保"`
	PersonalIncomeTax float32   `json:"personalIncomeTax" gorm:"size:20;comment:个人所得税"`
	ActualSalary      float32   `json:"actualSalary" gorm:"size:20;comment:实际工资"`
	ReleaseDate       time.Time `json:"releaseDate" gorm:"size:20;comment:发薪日期"`
	models.ControlBy
	models.ModelTime
}

/*
表名
*/
func (*SysUserSalary) TableName() string {
	return "sys_user_salary"
}
func (e *SysUserSalary) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysUserSalary) GetId() interface{} {
	return e.UserId
}
