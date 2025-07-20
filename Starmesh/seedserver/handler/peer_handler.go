package handler

import (
	"encoding/json"
	"net/http"
	"starmesh/common/model"
	"starmesh/seedserver/storage"
)

func RegisterPeerHandler(w http.ResponseWriter, r *http.Request) {
	var p model.Peer
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	storage.AddPeer(p)
	w.WriteHeader(http.StatusOK)
}

func GetPeersHandler(w http.ResponseWriter, r *http.Request) {
	peers := storage.GetPeers(10)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peers)
}
