package server

import (
	"net"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/GoSeoTaxi/tg_echo/internal/service/notifier"
)

type Handler struct {
	n   *notifier.Notifier
	log *zap.Logger
	mux *http.ServeMux
}

func NewHandler(n *notifier.Notifier, log *zap.Logger) *Handler {
	h := &Handler{n: n, log: log, mux: http.NewServeMux()}
	h.mux.HandleFunc("/tg", h.send)
	return h
}

func (h *Handler) Router() http.Handler { return h.mux }

func (h *Handler) send(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("msg")
	if text == "" {
		http.Error(w, "msg param required", http.StatusBadRequest)
		return
	}
	text = strings.Trim(text, `"`)

	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}

	if err := h.n.Send(notifier.Message{Body: text, Time: time.Now().UTC(), IP: ip}); err != nil {
		h.log.Error("send", zap.Error(err))
		http.Error(w, "telegram error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
