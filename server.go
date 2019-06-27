package main

import (
	"./config"
	"./logger"
	"./pages"
	"./repository"
	packageRouter "./router"
	"./token"
	"./utils"
	"github.com/gorilla/mux"

	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
)

func main() {

	print("" +
		" ___   ___   _     _     ____  ___\n" +
		"| |_) / / \\ | | | | |\\ |  / / | |_)\n" +
		"|_|_) \\_\\_/ \\_\\_/ |_| \\| /_/_ |_| \\\n\n")

	config.Init("./config.yml")
	repository.Init()
	packageRouter.Init()
	token.Init()
	pages.LoadPages("html/*.html")

	log := logger.GetLogger()
	hostname := config.IAM.Server.Hostname
	port := config.IAM.Server.Port
	host := hostname + ":" + port
	router := packageRouter.NewRouter()

	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, _ := route.GetPathTemplate()
		met, err := route.GetMethods()
		if err == nil {
			log.Info(host+tpl,
				zap.Strings("verbs", met),
			)
		}
		return nil
	})
	if err != nil {
		log.Error("router.Walk", zap.Error(err))
	}
	srv := &http.Server{
		Handler:      router,
		Addr:         host,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	//cert generation
	//verify if certificate exists
	certFile := config.IAM.Server.Certificate
	if len(certFile) == 0 {
		certFile = "./cert.pem"
	}
	keyFile := config.IAM.Server.PrivateKey
	if len(keyFile) == 0 {
		keyFile = "./key.pem"
	}
	genCert := false
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		genCert = true
	}
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		genCert = true
	}
	if genCert {
		utils.GetNewSSLCert()
	}
	//end of cert generation

	//start tls server
	err = srv.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		log.Error("srv.ListenAndServeTLS", zap.Error(err))
	}
}
