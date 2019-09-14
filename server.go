package main

import (
	"bounzr/iam/config"
	"bounzr/iam/logger"
	"bounzr/iam/pages"
	"bounzr/iam/repository"
	packageRouter "bounzr/iam/router"
	"bounzr/iam/token"
	"bounzr/iam/utils"
	"github.com/gorilla/mux"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	configFilePath string
)

func main() {
	print("" +
		" ___   ___   _     _     ____  ___\n" +
		"| |_) / / \\ | | | | |\\ |  / / | |_)\n" +
		"|_|_) \\_\\_/ \\_\\_/ |_| \\| /_/_ |_| \\\n\n")

	app := cli.NewApp()
	app.Name = "bounzr server"
	app.Usage = "identity and access management server"
	app.Version = "1.0.0-SNAPSHOT"

	startCommand := cli.Command{
		Name:    "start",
		Aliases: []string{"s"},
		Usage:   "start the bounzr server",
		Action:  startFunc,
	}
	startCommand.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "load configuration from `FILE`",
			Value:       "./config.yml",
			Destination: &configFilePath,
		},
	}
	app.Commands = []cli.Command{
		startCommand,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func startFunc(c *cli.Context) error {

	config.Init(configFilePath)
	repository.Init()
	packageRouter.Init()
	token.Init()
	pages.Init()

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
		return err
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
		return err
	}
	return nil
}
