package main

import (
	"embed"
	_ "embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	config "test/config"
	db "test/db"
	"test/service"

	"github.com/gorilla/mux"

	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/spf13/viper"
)

//go:embed  static
var static embed.FS

func main() {
	fmt.Println("entered into main method")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	config.LoadConfig()
	db.Connect()
	db.Init()
	count := db.CreateDefaultUser()
	log.Println("total number of users in usercollection", count)
	rpcServer := rpc.NewServer()
	rpcServer.RegisterCodec(json.NewCodec(), "application/json")
	rpcService := new(service.RPCService)
	rpcServer.RegisterService(rpcService, "rpcService")
	//step 2 : get the values from http and assign the required functions
	router := mux.NewRouter()
	//step 3
	rrpc := router.PathPrefix("/api").Subrouter()
	webapp, errs := fs.Sub(static, "static")
	if errs != nil {
		log.Println(errs)
	}
	//step 4
	router.PathPrefix("/").Handler(http.FileServer(http.FS(webapp)))
	//step 5
	rrpc.Handle("", rpcServer)
	//step 6
	rrpc.Use(service.Middleware)
	//step 1 : getting the url,body from the http and pass into router
	log.Println("listening on.....", viper.GetString("port"))
	err := http.ListenAndServe(viper.GetString("port"), router)
	if err != nil {
		log.Println(err)
	}

}

//step 1 : hit from browswer
//step 2 : reach the port to the server (main method) and router plays it's role
//step 3 :  if the url is /api means it serve the angular file
//step 4 : in browser login page displayed
//step 5 : signin from login it reaches the main method
//step 6 : again router play it's role
//step 7 : in rrpc handle function we given the rpcserver ..so that they handle the request
//step 8 : and then goes to middle ware from there all the api's called
