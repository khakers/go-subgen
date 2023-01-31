package webhooks

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"go-subgen/pkg"
)

type TautulliWebhookData struct {
	Event     string `json:"event"`
	File      string `json:"file"`
	Filename  string `json:"filename"`
	Mediatype string `json:"mediatype"`
}

func ServeTautulli(w http.ResponseWriter, r *http.Request) {
	var data TautulliWebhookData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(400)
		fmt.Fprintf(w, "error: %v", err)
		return
	}
	pkg.EnqueueSub(data.File)
	w.WriteHeader(200)
}
