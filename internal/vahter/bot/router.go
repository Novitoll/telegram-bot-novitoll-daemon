package bot

import (
	"log"
	"bytes"
	"net/http"
	"encoding/json"

	cfg "github.com/novitoll/novitoll_daemon_bot/config"
)

type App struct {
	Features *cfg.FeaturesConfig
	Lang string
}

func (app *App) RegisterHandlers() {
	http.HandleFunc("/process", app.ProcessMessageHandlerFunc)
	http.HandleFunc("/flushQueue", app.FlushQueueHandlerFunc)

	log.Printf("[+] Handlers for HTTP end-points are registered")
}

func (app *App) ProcessMessageHandlerFunc(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		msg := &RouteError{w, 400, nil, "Please send a request body"}
		log.Fatalln(msg.Error())
		return
	}

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		msg := &RouteError{w, 400, nil, "Could not parse request body"}
		log.Fatalln(msg.Error())
		return
	}
	log.Println(buf.String()) // TODO: remove

	var br BotIngressRequest
	err = json.Unmarshal([]byte(buf.String()), &br)
	if err != nil {
		msg := &RouteError{w, 400, nil, "Please send a valid JSON"}
		log.Fatalln(msg.Error())
		return
	}

	go br.Process(app)
    w.WriteHeader(http.StatusAccepted)
}

func (app *App) FlushQueueHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
}