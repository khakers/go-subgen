package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	config "go-subgen/pkg"
	"go-subgen/web/webhooks"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("using model type %v for language %s", conf.ModelType, conf.TargetLang)

	log.SetLevel(conf.LogLevel)
	log.Debugf("%+v", config.Cfg)

	http.Handle("/webhooks/tautulli", http.HandlerFunc(webhooks.ServeTautulli))
	http.Handle("/webhooks/generic", http.HandlerFunc(webhooks.ServeGeneric))

	config.StartWorkers(conf)

	err = http.ListenAndServe(":8095", nil)
	if err != nil {
		log.WithError(err).Fatal("web server failure")
	}

}
