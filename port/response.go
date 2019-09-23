package port

import (
	"encoding/json"
	"net/http"
	"strconv"
)

const (
	contentType   = "Content-Type"
	contentLength = "Content-Length"
	appjson       = "application/json"
)

// RespondRestWithError generates a common json error for REST APIs.
// @see #respondRestWithJSON
func RespondRestWithError(w http.ResponseWriter, code int, message string) {
	log.Infof("creatig error response with code: %d and message: %s", code, message)
	RespondRestWithJSON(w, code, map[string]string{"error": message})
}

// RespondRestWithJSON generates a common json response for REST APIs.
func RespondRestWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	log.Debugf("creatig response with code: %d and payload: %v", code, payload)
	response, _ := json.Marshal(payload)

	w.Header().Set(contentType, appjson)
	w.Header().Set(contentLength, strconv.Itoa(len(response)))
	w.WriteHeader(code)
	w.Write(response)
}
