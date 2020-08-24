package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"restaurant/common/serve"
	"restaurant/pkg"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

func init() {
	pkg.SQLInit()
}

func main() {
	AutoLoader()
}

//AutoLoader
func AutoLoader() {
	http.HandleFunc("/register", serve.Register)
	http.HandleFunc("/login", serve.Login)
	http.HandleFunc("/setFavourite", serve.SetFavourite)
	http.HandleFunc("/unsetFavourite", serve.UnsetFavourite)
	http.HandleFunc("/logout", serve.Logout)
	http.HandleFunc("/get_businesses", serve.GetBusinesses)
	http.HandleFunc("/reserve", serve.Reserve)
	go func() {
		//set listening port
		if err := http.ListenAndServe(":9090", nil); err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()
	//avoid shutdown by accident
	lend := make(chan bool)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			logrus.Info("all back job finished,now shutdown http server...")
			Shutdown()
			logrus.Info("success shutdown")
			lend <- true
			break
		}
	}()
	<-lend
}

// Shutdown operations before quit
func Shutdown() {
	pkg.Ddb.Close()
	pkg.Rds.Close()
}
