package httpstring

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func StringHandlerRouter(s []byte) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write(s)
	}
}

func StringHandlerFunc(s []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(s)
	}
}
