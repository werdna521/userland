package response

import "net/http"

func OK(w http.ResponseWriter, v interface{}) httpResponse {
	return httpResponse{
		statusCode: http.StatusOK,
		w:          w,
		v:          v,
	}
}
