package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/GoSeoTaxi/tg_echo/internal/config"
	"github.com/GoSeoTaxi/tg_echo/internal/logger"
	"github.com/GoSeoTaxi/tg_echo/internal/service/notifier"
	"github.com/GoSeoTaxi/tg_echo/internal/transport/server"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel)
	logger.ReplaceGlobals(log)

	n, err := notifier.New(cfg.BotToken, cfg.ChatID, log)
	if err != nil {
		log.Fatal("init notifier", zap.Error(err))
	}

	err = n.Send(notifier.Message{Body: "🚀 tg_echo Starting", Time: time.Now().UTC()})
	if err != nil {
		log.Fatal("notifier", zap.Error(err))
	}

	h := server.NewHandler(n, log)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: h.Router(),
	}

	// graceful shutdown context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// run server
	go func() {
		log.Info("listen", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("http", zap.Error(err))
		}
	}()

	// wait for signal
	<-ctx.Done()
	log.Info("shutdown signal received")

	shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shCtx); err != nil {
		log.Error("server shutdown", zap.Error(err))
	} else {
		log.Info("server gracefully stopped")
	}

	err = n.Send(notifier.Message{Body: "⛔ tg_echo Stopped", Time: time.Now().UTC()})
	if err != nil {
		log.Fatal("notifier", zap.Error(err))
	}
}
