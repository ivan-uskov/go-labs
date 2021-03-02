package orderservice

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func helloWorld(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprint(w, "hello-world")
	if err != nil {
		log.Error(err)
	}
}

func orders(w http.ResponseWriter, _ *http.Request) {
	orders := OrdersList{
		Orders: []Order{
			{ID: "3fa85f64-5717-4562-b3fc-2c963f66afa6", MenuItems: []MenuItem{{ID: "3fa85f64-5717-4562-b3fc-2c963f66afa6", Quantity: 0}}},
		},
	}

	renderJson(w, orders)
}

func renderJson(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func order(w http.ResponseWriter, r *http.Request) {
	id, found := mux.Vars(r)["ID"]
	if found {
		if id == "3fa85f64-5717-4562-b3fc-2c963f66afa6" {
			order := OrderDetails{
				Order: Order{ID: "3fa85f64-5717-4562-b3fc-2c963f66afa6", MenuItems: []MenuItem{{ID: "3fa85f64-5717-4562-b3fc-2c963f66afa6", Quantity: 0}}},
				Cost:  1,
				Time:  1,
			}

			renderJson(w, order)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	_, err := fmt.Fprint(w, "Not found")
	if err != nil {
		log.Error(err)
	}
}

func logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		h.ServeHTTP(w, r)
		log.WithFields(log.Fields{
			"method":     r.Method,
			"url":        r.URL,
			"remoteAddr": r.RemoteAddr,
			"userAgent":  r.UserAgent(),
			"duration":   time.Now().Sub(startTime).String(),
			"at":         startTime,
		}).Info("got request")
	})
}

func Router() http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/hello-world", helloWorld).Methods(http.MethodGet)
	s.HandleFunc("/orders", orders).Methods(http.MethodGet)
	s.HandleFunc("/order/{ID:[0-9a-zA-Z]+}", order).Methods(http.MethodGet)

	return logMiddleware(r)
}
