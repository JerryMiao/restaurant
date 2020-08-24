package serve

import (
	"net/http"
	"restaurant/pkg"
	"strconv"

	"github.com/go-redis/redis"
)

//GeoRadiusQuerys geo parameters
//key，member，unit，sort（asc/desc），count（limitation of sql）,radius，longitude, latitude
//withdist withcoord: Returns the center distance and longitude and latitude of the location name
//WithGeoHash: 52-bit signed integer
//store: Store the geographic location information of the returned results to the specified key
//StoreDist: Store the distance of the returned result from the center node to the specified key
type GeoRadiusQuerys struct {
	Key       string
	Member    string
	Longitude float64
	Latitude  float64
	Radius    float64
	// Can be m, km, ft, or mi. Default is km.
	Unit        string
	WithCoord   bool
	WithDist    bool
	WithGeoHash bool
	Count       int
	// Can be ASC or DESC. Default is no sort order.
	Sort      string
	Store     string
	StoreDist string
}

//GetRadiusMember get numbers that near members by name，has the same parameters with GetRadius
func GetRadiusMember(geoQuery GeoRadiusQuerys) ([]redis.GeoLocation, error) {
	gq := redis.GeoRadiusQuery{
		Radius: geoQuery.Radius,
		// Can be m, km, ft, or mi. Default is km.
		Unit:        geoQuery.Unit,
		WithCoord:   geoQuery.WithCoord,
		WithDist:    geoQuery.WithDist,
		WithGeoHash: geoQuery.WithGeoHash,
		Count:       geoQuery.Count,
		// Can be ASC or DESC. Default is no sort order.
		Sort:      geoQuery.Sort,
		Store:     geoQuery.Store,
		StoreDist: geoQuery.StoreDist,
	}
	return pkg.Rds.GeoRadiusByMember(geoQuery.Key, geoQuery.Member, &gq).Result()
}
func getFormData(r *http.Request) (GeoRadiusQuerys, error) {
	var err error
	longitude, err := strconv.ParseFloat(r.Form["longitude"][0], 64)
	latitude, err := strconv.ParseFloat(r.Form["latitude"][0], 64)
	radius, err := strconv.ParseFloat(r.Form["radius"][0], 64)
	count, err := strconv.Atoi(r.Form["count"][0])
	if err != nil {
		return GeoRadiusQuerys{}, err
	}
	withcoord, withdist, withgeoHash := false, false, false
	if r.Form["withcoord"][0] == "true" {
		withcoord = true
	}
	if r.Form["withdist"][0] == "true" {
		withdist = true
	}
	if r.Form["withgeoHash"][0] == "true" {
		withgeoHash = true
	}
	geo := GeoRadiusQuerys{
		Key:       r.Form["key"][0],
		Member:    r.Form["member"][0],
		Longitude: longitude,
		Latitude:  latitude,
		Radius:    radius,
		// Can be m, km, ft, or mi. Default is km.
		Unit:        r.Form["unit"][0],
		WithCoord:   withcoord,
		WithDist:    withdist,
		WithGeoHash: withgeoHash,
		Count:       count,
		// Can be ASC or DESC. Default is no sort order.
		Sort:      r.Form["sort"][0],
		Store:     r.Form["store"][0],
		StoreDist: r.Form["storedist"][0],
	}

	return geo, nil
}

//SetFavourite
func SetFavourite(w http.ResponseWriter, r *http.Request) {
	if pkg.Rds.Get(r.Header.Get("token")).Val() == "" {
		SendJSON(w, http.StatusOK, SendMap{"success": false, "msg": "token验证失败"})
		return
	}
	geo, err := getFormData(r)
	if err != nil {
		SendJSON(w, http.StatusOK, SendMap{"success": false, "msg": err.Error()})
		return
	}
	res, err := GetRadiusMember(geo)
	if err != nil {
		SendJSON(w, http.StatusOK, SendMap{"success": false, "msg": err.Error()})
		return
	}
	SendJSON(w, http.StatusOK, SendMap{"success": true, "data": res})
}
