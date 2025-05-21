package http

import (
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/wilder2000/GOSimple/comm"
	"github.com/wilder2000/GOSimple/config"
	"github.com/wilder2000/GOSimple/dbmodel"
	"github.com/wilder2000/GOSimple/glog"
)

// EfKilSag~2`
var (
	KEY       = "WilderBWASTARTUP_s3cr3t_k3y"
	SecretKey = []byte(KEY)
	//CachedUsers = make(cmap[string]string)
	//CachedUsers = cmap.New[string, string]()
)

const (
	KeyUser   = "currentUser"
	BeaKey    = "Bearer"
	AuthorKey = "Authorization"
	TokenUser = "user_id"
	TokenIP   = "user_ip"
	TokenData = "flag"
	//SecretKey = "BWASTARTUP_s3cr3t_k3y"
)

func ValidTokenHttpRequest(c gin.Context) (userid string, err error) {

	var userID string
	authHeader := c.GetHeader(AuthorKey)
	glog.Logger.InfoF("found GET URL token:%s", authHeader)
	token, err := ValidateToken(authHeader)
	if err != nil {
		return userID, errors.New("invalid http request. for " + err.Error())
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return userID, errors.New("http request token invalid.")
	}
	userID = claim[TokenUser].(string)
	_, err2 := RefreshToken(c, userID)
	if err2 != nil {
		glog.Logger.Warn(err2.Error())
		return userID, errors.New("http request token refresh failed.")
	}
	return userID, nil
}
func ParseHttpRequest(c gin.Context) (userid string, err error) {
	authHeader := c.GetHeader(AuthorKey)
	var userID string
	glog.Logger.InfoF("found token:%s", authHeader)
	if !strings.Contains(authHeader, BeaKey) {
		return userID, errors.New("Bad http request.")
	}
	tokenString := ""
	arrayToken := strings.Split(authHeader, " ")
	if len(arrayToken) == 2 {
		tokenString = arrayToken[1]
	}
	token, err := ValidateToken(tokenString)
	if err != nil {
		return userID, errors.New("invalid http request. for " + err.Error())
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return userID, errors.New("http request token invalid.")
	}
	userID = claim[TokenUser].(string)
	//tokid, ok := CachedUsers[userID]
	//if !ok {
	//	//找不到缓存的，则认为需要重新登录
	//	glog.Logger.InfoF("Not found cached user:", tokid)
	//	return userID, errors.New("http request token timeout.")
	//}
	//refresh token
	//to, err := RefreshToken(c, userID)
	_, err2 := RefreshToken(c, userID)
	if err2 != nil {
		glog.Logger.Warn(err2.Error())
		return userID, errors.New("http request token refresh failed.")
	}

	//CachedUsers[userID] = to
	//CachedUsers.Set(userID, to)
	return userID, nil
}
func RefreshToken(c gin.Context, userID string) (string, error) {
	newToken, err := GenerateToken(userID)
	if err != nil {
		glog.Logger.Warn("refresh token failed. for:" + err.Error())
		return "", err
	}
	c.Header(AuthorKey, newToken)
	return newToken, nil
}
func ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		//mc := token.Claims.(jwt.MapClaims)

		if !ok {
			return nil, errors.New("Invalid token")
		}
		return []byte(SecretKey), nil
	})
	if err != nil {
		return token, err
	}

	return token, nil
}
func GenerateToken(userID string) (string, error) {
	nowTime := comm.LocalTime()
	expireTime := nowTime.Add(config.AConfig.ExpireTime * time.Minute)
	claim := jwt.MapClaims{}
	claim[TokenUser] = userID
	claim["exp"] = expireTime.Unix()
	claim[TokenIP] = expireTime.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString(SecretKey)
	if err != nil {
		return signedToken, err
	}
	return signedToken, nil
}
func GenerateToken2(userID string, data string) (string, error) {
	nowTime := comm.LocalTime()
	expireTime := nowTime.Add(config.AConfig.ExpireTime * time.Minute)
	claim := jwt.MapClaims{}
	claim[TokenUser] = userID
	claim["exp"] = expireTime.Unix()
	claim[TokenIP] = expireTime.Unix()
	claim[TokenData] = data

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString(SecretKey)
	if err != nil {
		return signedToken, err
	}
	return signedToken, nil
}

var (
	usersMap sync.Map
)

// LookupUser 从会话中读取用户
func LookupUser(userid string) (*dbmodel.SUser, bool) {
	user, exist := usersMap.Load(userid)
	if exist {
		return user.(*dbmodel.SUser), true
	}
	return nil, false
}
func SaveUser(u *dbmodel.SUser) {
	usersMap.Store(u.ID, u)
}
