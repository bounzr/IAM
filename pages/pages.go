package pages

import (
	"bounzr/iam/config"
	"bounzr/iam/logger"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"html/template"
	"net/http"
	"os"
)

var (
	log *zap.Logger
)

var defaultPages = map[string]string{
	"authorize": "./html/authorize.html",
	"index":     "./html/index.html",
	"login":     "./html/login.html",
	"signup":    "./html/signup.html",
}

func Init() {
	log = logger.GetLogger()
}

func RenderPage(w http.ResponseWriter, name string, data interface{}) error {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(dir)

	templates := config.CustomConf["webpages"]
	templatesMap := templates.(map[interface{}]interface{})
	templatePath, found := templatesMap[name].(string)
	if !found {
		log.Error("web template not found", zap.String("name", name))
		return errors.New("missing web template " + name)
	}
	log.Info("web template", zap.String("name", name), zap.String("path", templatePath))
	if len(templatePath) == 0 {
		templatePath = defaultPages[name]
	}
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return errors.New("can not parse web template " + name)
	}
	return t.Execute(w, data)
}
