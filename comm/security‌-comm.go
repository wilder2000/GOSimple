package comm

import (
	"crypto/md5"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wilder2000/GOSimple/glog"
	"golang.org/x/crypto/bcrypt"
)

func UUID() string {
	uniqueID, err := uuid.NewRandom()
	if err != nil {
		glog.Logger.ErrorF("could not generate UUIDv4,for:$s", err.Error())
		return ""
	}
	return strings.ToUpper(uniqueID.String())
}
func LowerUUID() string {
	uniqueID, err := uuid.NewRandom()
	if err != nil {
		glog.Logger.ErrorF("could not generate UUIDv4,for:$s", err.Error())
		return ""
	}
	return strings.ToLower(uniqueID.String())
}

//e2c569be17396eca2a2e3c11578123ed

func MD5(messages ...string) string {
	h := md5.New()
	for _, msg := range messages {
		io.WriteString(h, msg)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

type RandomObject struct {
}

func (r *RandomObject) Init() {
	rand.Seed(time.Now().UnixNano()) // 初始化随机种子
}
func (r *RandomObject) Next() string {
	randomNumber := rand.Intn(999999) // 生成一个0到999999之间的随机数
	return fmt.Sprintf("%06d", randomNumber)
}
func EPassword(pwd string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
}
