package serve

import (
	"net/http"
	"restaurant/pkg"
)

//UnsetFavourite delete a geo member
func UnsetFavourite(w http.ResponseWriter, r *http.Request) {
	if pkg.Rds.Get(r.Header.Get("token")).Val() == "" {
		SendJSON(w, http.StatusOK, SendMap{"success": false, "msg": "token验证失败"})
		return
	}
	key, member := r.Form["key"][0], r.Form["member"][0]
	if _, err := DeleteGeoMember(key, member); err != nil {
		SendJSON(w, http.StatusOK, SendMap{"success": false, "msg": err.Error()})
		return
	}
	SendJSON(w, http.StatusOK, SendMap{"success": true})
}

//because the implementation of GEO is based on zset, use zrem to delete geo information
func DeleteGeoMember(key string, member ...string) (int64, error) {
	zm := []interface{}{}
	for _, v := range member {
		zm = append(zm, v)
	}
	return pkg.Rds.ZRem(key, zm...).Result()
}
