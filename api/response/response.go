package response

import (
	"encoding/json"
	"net/http"
)

type httpResponse struct {
	statusCode int
	w          http.ResponseWriter
	v          interface{}
}

func (r httpResponse) JSON() {
	r.w.Header().Set("Content-Type", "application/json")
	r.w.WriteHeader(r.statusCode)
	json.NewEncoder(r.w).Encode(r.v)
}
