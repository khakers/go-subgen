package main

import (
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
	config "go-subgen/pkg"
	"go-subgen/pkg/configuration"
	"go-subgen/web/webhooks"
)

func main() {
	conf, err := configuration.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.SetLevel(conf.LogLevel)
	log.Printf("using model type %v for language %s", conf.ModelType, conf.TargetLang)

	log.Debugf("%+v", configuration.Cfg)

	downloaded, err := config.IsModelDownloaded(conf.ModelType)
	if err != nil {
		log.WithError(err).Errorln("Model check failed (this is likely normal and can be ignored)")
	}
	if !downloaded {
		err := config.DownloadModel(conf.ModelType)
		if err != nil {
			log.WithError(err).Fatalln("failed to download model")
		}
	}

	http.Handle("/webhooks/tautulli", http.HandlerFunc(webhooks.ServeTautulli))
	http.Handle("/webhooks/generic", http.HandlerFunc(webhooks.ServeGeneric))
	http.Handle("/webhooks/radarr", http.HandlerFunc(webhooks.ServeRadarr))
	http.Handle("/webhooks/sonarr", http.HandlerFunc(webhooks.ServeSonarr))

	config.StartWorkers(conf)
	err = http.ListenAndServe(":"+strconv.Itoa(int(conf.Port)), nil)
	if err != nil {
		log.WithError(err).Fatal("web server failure")
	}

}
