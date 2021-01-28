package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//Controller controller
type Controller struct {
	config DmuxConf
	router *mux.Router
}

func (c *Controller) start() {
	c.router = mux.NewRouter()
	c.initializeRoutes()
}

func (c *Controller) initializeRoutes() {
	c.router.HandleFunc("/v1/config", c.getConfig).Methods("GET")
	log.Fatal(http.ListenAndServe(":1234", c.router))
}
func (c *Controller) getConfig(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, c.config)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
