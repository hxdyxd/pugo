package model

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

type UserToken struct {
	Id     int64
	Hash   string `xorm:"VARCHAR(200) notnull index(hash)"`
	UserId int64  `xorm:"notnull"`

	CreateTime int64 `xorm:"INT(12) created"`
	ExpireTime int64 `xorm:"INT(12)"`

	From string `xorm:"VARCHAR(50) notnull index(origin)"`
	Note string
}

func (ut *UserToken) SetHash(hash string) {
	m := md5.New()
	m.Write([]byte(hash))
	ut.Hash = hex.EncodeToString(m.Sum(nil))
}

func (ut *UserToken) IsExpired() bool {
	return time.Now().Unix() > ut.ExpireTime
}
