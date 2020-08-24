package serve

import (
	"encoding/base64"
	"math/rand"
	"net/http"
	"restaurant/pkg"
	"time"

	"github.com/pborman/uuid"
)

//BuildIserSalt Randomly obtain a segment of +uuid from the user to generate random salt to prevent the code leakage
//and password generation process from being cracked
func BuildIserSalt(user string) string {
	rand.Seed(time.Now().UnixNano())
	sl := rand.Intn(len(user))
	return user[sl:] + base64.RawURLEncoding.EncodeToString(uuid.NewUUID())
}

//Register Register
func Register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name, pwd := r.Form["name"][0], r.Form["pwd"][0]

	if name == "" || pwd == "" {
		SendJSON(w, http.StatusOK, SendMap{"success": false, "msg": "Required content cannot be empty"})
		return
	}
	if IsExistUser(name) {
		SendJSON(w, http.StatusOK, SendMap{"success": false, "msg": "User already exists"})
		return
	}

	salt := BuildIserSalt(name)
	bPwd := BuildPas(pwd, salt)
	if err := pkg.Ddb.Exec("INSERT INTO user(name,password,salt)VALUES(?,?,?)", name, bPwd, salt).Error; err != nil {
		SendJSON(w, http.StatusOK, err.Error())
		return
	}
	SendJSON(w, http.StatusOK, SendMap{"success": true})
}
