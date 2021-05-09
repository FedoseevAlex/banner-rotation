package server

import (
	"log"
	"net/http"
	"time"

	"github.com/FedoseevAlex/banner-rotation/internal/types"
	"github.com/julienschmidt/httprouter"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriterWrapper) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func loggingMiddleware(next httprouter.Handle, logger types.Logger) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		wrappedResponseWriter := &responseWriterWrapper{ResponseWriter: w, status: http.StatusOK}
		log.Println("wrappedResponseWriter ", wrappedResponseWriter)

		begin := time.Now()
		next(wrappedResponseWriter, r, params)
		duration := time.Since(begin)

		info := map[string]interface{}{
			"ip":          r.RemoteAddr,
			"timestamp":   begin.Format(time.RFC822Z),
			"method":      r.Method,
			"path":        r.URL.Path,
			"HTTP ver.":   r.Proto,
			"status code": wrappedResponseWriter.status,
			"latency":     duration.String(),
			"user agent":  r.UserAgent(),
		}
		logger.Trace("Request", info)
	}
}
