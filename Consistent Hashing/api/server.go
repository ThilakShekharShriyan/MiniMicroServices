package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	hashring "github.com/thilakshekharshriyan/hashring"
)


// Wraps the HashRing with HTTP handlers.
// Encapsulates ring logic within a web server.
type Server struct {
	Ring *hashring.HashRing
}


// Creates a new hash ring with the specified number of virtual node replicas.
func NewServer(replicas int) *Server {
	return &Server{Ring: hashring.New(replicas)}
}

// API to Handle Add Node
func (s *Server) AddNodeHandler(w http.ResponseWriter, r *http.Request) {

	// Expects JSON Body as { "name": "NodeA" }
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	log.Printf("Adding node: %s", body.Name)

	//Adds the node to the ring.
	s.Ring.Add(body.Name)

	//Returns 200 OK on success.
	w.WriteHeader(http.StatusOK)
}

// Expects URL path param: /nodes/{node}
func (s *Server) RemoveNodeHandler(w http.ResponseWriter, r *http.Request) {
	node := mux.Vars(r)["node"]
	log.Printf("Removing node: %s", node)
	s.Ring.Remove(node)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) LookupHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}
	node := s.Ring.Get(key)
	log.Printf("Lookup for key '%s' → node '%s'", key, node)
	json.NewEncoder(w).Encode(map[string]string{"node": node})
}

// Returns all real nodes in the ring:
func (s *Server) ListNodesHandler(w http.ResponseWriter, r *http.Request) {
	nodes := s.Ring.Nodes()
	json.NewEncoder(w).Encode(nodes)
}

func (s *Server) Routes() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/nodes", s.AddNodeHandler).Methods("POST")
	r.HandleFunc("/nodes/{node}", s.RemoveNodeHandler).Methods("DELETE")
	r.HandleFunc("/lookup", s.LookupHandler).Methods("GET")
	r.HandleFunc("/nodes", s.ListNodesHandler).Methods("GET")
	return r
}
