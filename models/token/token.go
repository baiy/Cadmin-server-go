package token

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/baiy/Cadmin-server-go/models"
	"github.com/baiy/Cadmin-server-go/models/utils"
	"github.com/doug-martin/goqu/v9"
	"math/rand"
	"time"
)

type Model struct {
	models.Model
	Token       string     `json:"token"`
	AdminUserId int        `json:"admin_user_id"`
	ExpireTime  utils.Time `json:"expire_time"`
}

func (m *Model) IsExpire() bool {
	return time.Time(m.ExpireTime).Before(time.Now())
}

func GetByToken(token string) (model *Model, err error) {
	model = new(Model)
	found, err := models.Db.From("admin_token").Where(goqu.Ex{
		"token": token,
	}).ScanStruct(model)
	if err == nil {
		if !found {
			err = errors.New("token不存在")
		}
	}
	return
}

func Clear() {
	_, _ = models.Db.Delete("admin_token").
		Where(goqu.Ex{"expire_time": goqu.Op{"lt": time.Now().Format("2006-01-02 15:04:05")}}).
		Executor().Exec()
}

func Remove(t string) {
	_, _ = models.Db.Delete("admin_token").
		Where(goqu.Ex{"token": t}).
		Executor().Exec()
}

func Add(userId int) string {
	tokenStr := generate()
	dd, _ := time.ParseDuration("48h")
	_, _ = models.Db.Insert("admin_token").Rows(map[string]interface{}{
		"admin_user_id": userId,
		"token":         tokenStr,
		"expire_time":   time.Now().Add(dd).Format("2006-01-02 15:04:05"),
	}).Executor().Exec()
	return tokenStr
}

func generate() string {
	rand.Seed(time.Now().UnixNano())
	str := fmt.Sprintf("%d%d", time.Now().UnixNano(), rand.Intn(9999))

	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
