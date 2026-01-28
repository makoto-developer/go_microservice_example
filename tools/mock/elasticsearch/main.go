package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

type IndexDocumentRequest struct {
	ID   string                 `json:"id,omitempty"`
	Data map[string]interface{} `json:"data"`
}

type SearchRequest struct {
	Query map[string]interface{} `json:"query"`
	From  int                    `json:"from,omitempty"`
	Size  int                    `json:"size,omitempty"`
}

var documents = make(map[string]map[string]interface{})

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9200"
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/", elasticsearchHandler)

	log.Printf("Elasticsearch Mock Server listening on :%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"name":    "mock-elasticsearch",
		"version": map[string]string{"number": "8.0.0"},
	})
}

func elasticsearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// POST /{index}/_doc/{id} - Index document
	if r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/_doc/") {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) >= 3 {
			index := parts[0]
			docID := parts[2]

			var doc map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			key := index + ":" + docID
			documents[key] = doc

			log.Printf("[Elasticsearch Mock] Indexed document: %s in index: %s\n", docID, index)

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"_index":  index,
				"_id":     docID,
				"result":  "created",
				"_version": 1,
			})
			return
		}
	}

	// POST /{index}/_search - Search documents
	if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/_search") {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		index := parts[0]

		log.Printf("[Elasticsearch Mock] Search request for index: %s\n", index)

		// Mock: return all documents for this index
		hits := []map[string]interface{}{}
		for key, doc := range documents {
			if strings.HasPrefix(key, index+":") {
				docID := strings.TrimPrefix(key, index+":")
				hits = append(hits, map[string]interface{}{
					"_index":  index,
					"_id":     docID,
					"_source": doc,
				})
			}
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"took": 5,
			"hits": map[string]interface{}{
				"total": map[string]int{"value": len(hits)},
				"hits":  hits,
			},
		})
		return
	}

	http.Error(w, "Not found", http.StatusNotFound)
}
