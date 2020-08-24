package serve

import (
	"net/http"
	"restaurant/pkg"
)

//Reserve Reserve
// Get user id from token, get restaurant id and save
// Need to verify user appointment information
// Use transactions to prevent simultaneous modification
// No one makes an appointment at most three
func Reserve(w http.ResponseWriter, r *http.Request) {
	userID := pkg.Rds.Get(r.Header.Get("token")).Val()
	if userID == "" {
		SendJSON(w, http.StatusOK, SendMap{"success": false, "msg": "token verification failed"})
		return
	}
	businessID := r.Form["business_id"][0]

	db := pkg.Ddb.Begin()
	count := 0
	if err := db.Table("business").Where("user_id = ?", userID).Count(&count).Error; err != nil {
		db.Rollback()
		SendJSON(w, http.StatusOK, SendMap{"success": false, "msg": err.Error})
		return
	}
	if count >= 3 {
		db.Rollback()
		SendJSON(w, http.StatusOK, SendMap{"success": false, "msg": "Up to three appointments"})
		return
	}

	if err := db.Exec("INSERT INTO business(user_id,business_id)VALUES(?,?)", userID, businessID).Error; err != nil {
		db.Rollback()
		SendJSON(w, http.StatusOK, SendMap{"success": false, "msg": err.Error})
		return
	}
	db.Commit()
	SendJSON(w, http.StatusOK, SendMap{"success": true})
}
