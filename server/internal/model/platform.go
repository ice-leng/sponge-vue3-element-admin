package model

import (
	"admin/internal/types"
	"github.com/zhufuyi/sponge/pkg/ggorm"
	"time"
)

// Platform 管理员
type Platform struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	Username string              `gorm:"column:username;type:varchar(32);NOT NULL" json:"username"`                                                                                          // 账号
	Password string              `gorm:"column:password;type:varchar(64);NOT NULL" json:"password"`                                                                                          // 密码
	Avatar   string              `gorm:"column:avatar;type:varchar(255);default:https://oss.youlai.tech/youlai-boot/2023/05/16/811270ef31f548af9cffc026dfc3777b.gif;NOT NULL" json:"avatar"` // 头像
	RoleID   types.LocalIntArray `gorm:"column:role_id;type:json;NOT NULL" json:"roleID"`                                                                                                    // 角色
	Status   *int                `gorm:"column:status;type:tinyint(4);NOT NULL" json:"status"`                                                                                               // 状态
	LastTime *time.Time          `gorm:"column:last_time;type:datetime" json:"lastTime"`                                                                                                     // 上次登录时间
}

// TableName table name
func (m *Platform) TableName() string {
	return "t_platform"
}
