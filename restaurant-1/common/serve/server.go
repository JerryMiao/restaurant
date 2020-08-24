package serve

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"restaurant/pkg"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/scrypt"
)

//SendMap send structure
type SendMap map[string]interface{}

//buildUserPassword Generate cipher text based on cipher text and salt
func buildUserPassword(pwdMd5, salt []byte) ([]byte, error) {
	return scrypt.Key(pwdMd5, salt, 16384, 8, 1, 32)
}

//IsExistUser Determine whether the user exists, the existence is true
func IsExistUser(username string) bool {
	num := 0
	if err := pkg.Ddb.Exec("SELECT count(*) FROM user where name=?", username).Scan(&num).Error; err != nil {
		logrus.WithFields(logrus.Fields{"get user": err}).Error("user")
		return false
	}
	if num > 0 {
		return true
	}
	return false
}

//Front-end hex string
func huexEncode(md5Pwd string) []byte {
	decoded, err := hex.DecodeString(md5Pwd)
	if err != nil {
		logrus.WithFields(logrus.Fields{"decode": err}).Error("hex")
		return []byte{}
	}
	return decoded
}

//SendJSON SendJSON
func SendJSON(w http.ResponseWriter, statuscode int, data interface{}) {
	bt, err := json.Marshal(data)
	if err != nil {
		logrus.WithFields(logrus.Fields{"send err": err}).Error("SendJSON")
		return
	}
	w.WriteHeader(statuscode)
	w.Write(bt)
}
