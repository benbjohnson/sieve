package sieve

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"

	"github.com/gorilla/mux"
)

// NewHandler returns a new root web handler.
func NewHandler(db *DB) http.Handler {
	var r = mux.NewRouter()
	r.Handle("/", &IndexHandler{}).Methods("GET")
	r.Handle("/assets/{filename}", &AssetsHandler{}).Methods("GET")
	r.Handle("/subscribe", NewSubscribeHandler(db)).Methods("GET")
	return r
}

// IndexHandler returns the main index page.
type IndexHandler struct{}

// ServeHTTP sends the index page to the client.
func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := Asset("index.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(b)
}

// AssetsHandler handles HTTP requests for static files.
type AssetsHandler struct{}

// ServeHTTP sends the requested asset file.
func (h *AssetsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	b, err := Asset(vars["filename"])
	if err != nil {
		http.NotFound(w, r)
		return
	}
	switch path.Ext(vars["filename"]) {
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript")
	}
	w.Write(b)
}

// NewSubscribeHandler returns a new subscription handler for a database.
func NewSubscribeHandler(db *DB) *SubscribeHandler {
	return &SubscribeHandler{db}
}

// SubscribeHandler sends data to the client via server sent events.
type SubscribeHandler struct {
	db *DB
}

// ServeHTTP sends the database data.
func (h *SubscribeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	index, err := strconv.Atoi(r.FormValue("index"))
	if err != nil && len(r.FormValue("index")) > 0 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a channel for subscribing to database updates.
	ch := h.db.Subscribe(index)
	closeNotifier := w.(http.CloseNotifier).CloseNotify()

	// Mark this as an SSE event stream.
	w.Header().Set("Content-Type", "text/event-stream")

	// Continually stream updates as they come.
loop:
	for {
		select {
		case <-closeNotifier:
			break loop

		case row := <-ch:
			// Encode row as JSON.
			b, err := json.Marshal(row)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				break loop
			}

			// Send row as server-sent event.
			w.Write([]byte("data: "))
			w.Write(b)
			w.Write([]byte("\n\n"))
			w.(http.Flusher).Flush()
		}
	}

	// Unsubscribe from the database when the connection is lost.
	h.db.Unsubscribe(ch)
}
