package serve

import (
	"net/http"
	"restaurant/pkg"

	"github.com/sirupsen/logrus"
)

//Logout Delete a user
func Logout(w http.ResponseWriter, r *http.Request) {
	if pkg.Rds.Get(r.Header.Get("token")).Val() == "" {
		SendJSON(w, http.StatusOK, SendMap{"success": false, "msg": "token verification failed"})
		return
	}

	if err := pkg.Rds.Del(r.Header.Get("token")).Err(); err != nil {
		logrus.WithFields(logrus.Fields{"deluser": err}).Error("redisErr")
		SendJSON(w, http.StatusOK, SendMap{"success": false, "msg": err.Error()})
		return
	}
	SendJSON(w, http.StatusOK, SendMap{"success": true})
}
