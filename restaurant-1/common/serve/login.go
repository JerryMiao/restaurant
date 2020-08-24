package serve

import (
	"bytes"
	"net/http"
	"restaurant/pkg"
	"time"

	"github.com/pborman/uuid"
	"github.com/sirupsen/logrus"
)

//RegisterUser Set user information, userkey is the token in the corresponding header
func RegisterUser(userkey, username string) error {
	if err := pkg.Rds.Set(userkey, username, time.Minute*5).Err(); err != nil {
		logrus.WithFields(logrus.Fields{"set": err}).Error("redisErr")
		return err
	}
	return nil
}

//BuildPas Parse the hex cipher text of the front end and call the cipher text generation function
func BuildPas(pwd, salt string) []byte {
	bPwd, err := buildUserPassword(huexEncode(pwd), []byte(salt))
	if err != nil {
		logrus.WithFields(logrus.Fields{"pwd": err}).Error("validPwdMd5")
	}
	return bPwd
}

type user struct {
	Pwd    string
	UserID string
}

//Login Use MD5 encryption password
func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name, pwd := r.Form["name"][0], r.Form["pwd"][0]
	if name == "" || pwd == "" {
		SendJSON(w, http.StatusOK, SendMap{"success": false, "msg": "Required content cannot be empty"})
		return
	}
	if !IsExistUser(name) {
		SendJSON(w, http.StatusOK, SendMap{"success": false, "msg": "Incorrect username or password"})
		return
	}
	user := &user{}
	if err := pkg.Ddb.Table("user").Select("user_id,pwd").Where("name = ?", name).Scan(&user).Error; err != nil {
		SendJSON(w, http.StatusOK, SendMap{"success": false, "msg": err.Error})
		return
	}
	bPwd := BuildPas(pwd, user.Pwd)
	if bytes.Equal(bPwd, []byte(user.Pwd)) {
		token := uuid.NewUUID().String()
		if err := RegisterUser(token, user.UserID); err != nil {
			SendJSON(w, http.StatusOK, SendMap{"success": false, "msg": err.Error})
			return
		}
		SendJSON(w, http.StatusOK, SendMap{"token": token, "success": true})
		return
	}
	SendJSON(w, http.StatusOK, SendMap{"success": false, "msg": "Incorrect username or password"})
}
