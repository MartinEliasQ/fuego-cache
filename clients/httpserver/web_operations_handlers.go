package httpServer

import (
	"encoding/json"
	cache "github.com/tomiok/fuego-cache/fuego"
	"github.com/tomiok/fuego-cache/internal"
	"github.com/tomiok/fuego-cache/logs"
	"net/http"
	"strings"
)

type OperationsHandler struct {
	GetCallback    func(k interface{}) (string, error)
	SetCallback    func(k interface{}, v string, ttl int) (string, error)
	DeleteCallback func(k interface{}) (string, error)
	//bulks operations
	BulkGetCallback    func()
	BulkSetCallback    func(bulkEntry cache.BulkEntry) cache.BulkResponse
	BulkDeleteCallback func(keys []interface{})
}

func (o *OperationsHandler) GetValueHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		key := strings.TrimPrefix(r.URL.Path, GetUrl)
		logs.Info(cache.ApplyHash(key))
		var _key interface{} = key
		res, err := o.GetCallback(_key)

		//when a response is with error true and value is nil, it means that the key is not present in the cache
		_ = json.NewEncoder(w).Encode(HTTPRes{Response: res, Err: err != nil})
	}
}

func (o *OperationsHandler) SetValueHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == PostMethod {
			body := r.Body
			var b SetRequest
			err := json.NewDecoder(body).Decode(&b)

			if err != nil || b.Key == "" {
				http.Error(w, "cannot process current request", http.StatusBadRequest)
				return
			}

			res, err := o.SetCallback(b.Key, b.Value, b.TTL)

			internal.OnCloseError(body.Close)
			_ = json.NewEncoder(w).Encode(HTTPRes{Response: res, Err: err != nil})
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

func (o *OperationsHandler) DeleteValueHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == DeleteMethod {
			id := strings.TrimPrefix(r.URL.Path, DeleteUrl)
			deleted, err := o.DeleteCallback(id)
			_ = json.NewEncoder(w).Encode(HTTPRes{Response: deleted, Err: err != nil})
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

func (o *OperationsHandler) BulkSetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == PostMethod {
			var cmd []BulkSetCMD
			body := r.Body
			defer body.Close()
			_ = json.NewDecoder(body).Decode(&cmd)
			res := o.BulkSetCallback(toBulkEntry(cmd))

			_ = json.NewEncoder(w).Encode(&res)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

func toBulkEntry(body []BulkSetCMD) cache.BulkEntry {
	var bulkEntry cache.BulkEntry
	for _, entry := range body {
		bulkEntry.Add(entry.Key, entry.Value, entry.TTL)
	}
	return bulkEntry
}

type BulkSetCMD struct {
	Key   interface{} `json:"key"`
	Value string      `json:"value"`
	TTL   int         `json:"ttl,omitempty"`
}