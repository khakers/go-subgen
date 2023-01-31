package webhooks

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"go-subgen/pkg"
)

type GenericWebhookData struct {
	File string `json:"file"`
}

func ServeGeneric(w http.ResponseWriter, r *http.Request) {
	var data GenericWebhookData
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
