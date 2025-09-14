package handler

import (
	"fmt"
	"go.uber.org/zap"
	"log"
	"net/http"
	logger2 "pg-summary-service/internal/logger"
	"runtime/debug"
)

type requestLog struct {
	Method     string `json:"method"`
	Path       string `json:"path"`
	URL        string `json:"url"`
	RemoteAddr string `json:"remoteAddr"`
}

func methodChecker(expected string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != expected {
			logger2.Log.Error("Method not allowed", zap.String("method", r.Method))
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		next(w, r)
	}
}

func logger(next http.HandlerFunc) http.HandlerFunc {
	src := "middleware-logger"
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("[Logger] Request received")
		logInfo := requestLog{
			Method:     r.Method,
			Path:       r.Host,
			URL:        r.URL.String(),
			RemoteAddr: r.RemoteAddr,
		}
		logger2.Log.Info(src+" incoming req details", zap.Any("request", logInfo))
		next(w, r)
	}
}

func panicRecovery(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Recovered from panic:", err)
				log.Println("Stack trace:", string(debug.Stack()))
				logger2.Log.Error("System Panic!!! recover ", zap.Any("err", err), zap.Any("stack", string(debug.Stack())))
				http.Error(w, "INTERNAL ERROR", http.StatusInternalServerError)
				return
			}
		}()

		next(w, r)
	}
}
func ApplyMiddlewares(method string, authType AuthType, next http.HandlerFunc) http.HandlerFunc {
	// all the middlewares goes here including auth middleware
	handlers := methodChecker(method, next)
	handlers = logger(handlers)
	handlers = panicRecovery(handlers)
	return handlers
}
