package tips

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/byblix/mongotest/storage"
)

const (
	database   = "bas"
	collection = "fez"
)

// CreateTip creates a logging for a tip in MongoDB
func CreateTip(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var tip storage.Tip
	err := json.NewDecoder(r.Body).Decode(&tip)
	if err != nil {
		http.Error(w, "Error decoding from body: "+err.Error(), http.StatusInternalServerError)
		return
	}
	db := &storage.DatabaseRef{Database: "bas", Collection: "fez"}
	res, err := db.InsertOneItem(&tip)
	if err != nil {
		http.Error(w, "Error inserting: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Response
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetAllTips get a single tip
func GetAllTips(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	defer r.Body.Close()
	db := &storage.DatabaseRef{Database: database, Collection: collection}
	ch := make(chan []*storage.Tip)
	go db.GetAllTips(ch, &wg)
	wg.Wait()
	res := <-ch
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
