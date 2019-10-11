package admin

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// 密码生成器接口
type Password interface {
	Hash(password []byte) []byte
	Verify(password, hash []byte) bool
}

// 注册密码生成器
func RegisterPassword(password Password) {
	Passworder = password
}

var Passworder Password

// 默认密码生成器
type PasswordDefault struct {
}

func (p *PasswordDefault) Hash(password []byte) []byte {
	return p.hashWithSalt(password, p.salt())
}

func (p *PasswordDefault) hashWithSalt(password, salt []byte) []byte {
	sha := p.sha265(append(p.sha265(append(password, salt...)), salt...))
	src := append(
		sha[:],
		append([]byte("|"), salt...)...
	)

	return []byte(base64.StdEncoding.EncodeToString(src))
}

func (p *PasswordDefault) Verify(password, hash []byte) bool {
	dst := make([]byte, 0)
	dst, err := base64.StdEncoding.DecodeString(string(hash))
	if err != nil {
		return false
	}
	hashItem := strings.Split(string(dst), "|")
	if len(hashItem) != 2 {
		return false
	}
	return string(hash) == string(p.hashWithSalt(password, []byte(hashItem[1])))
}

func (p *PasswordDefault) salt() []byte {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := make([]byte, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 8; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return result
}

func (p *PasswordDefault) sha265(data []byte) []byte {
	return []byte(fmt.Sprintf("%x", sha256.Sum256(data)))
}

func init() {
	RegisterPassword(new(PasswordDefault))
}
