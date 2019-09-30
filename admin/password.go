package admin

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Passwrod interface {
	Hash(password []byte) []byte
	Verify(password, hash []byte) bool
}

type PasswrodDefault struct {
}

func (p *PasswrodDefault) Hash(password []byte) []byte {
	return p.hashWithSalt(password, p.salt())
}

func (p *PasswrodDefault) hashWithSalt(password, salt []byte) []byte {
	sha := p.sha265(append(p.sha265(append(password, salt...)), salt...))
	src := append(
		sha[:],
		append([]byte("|"), salt...)...
	)

	return []byte(base64.StdEncoding.EncodeToString(src))
}

func (p *PasswrodDefault) Verify(password, hash []byte) bool {
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

func (p *PasswrodDefault) salt() []byte {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := make([]byte, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 8; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return result
}

func (p *PasswrodDefault) sha265(data []byte) []byte {
	return []byte(fmt.Sprintf("%x", sha256.Sum256(data)))
}