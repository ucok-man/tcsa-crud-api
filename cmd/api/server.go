package main

import (
	"context"
	"errors"
	"fmt"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ucok-man/tcsa/internal/tlog"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.Port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     stdlog.New(app.logger, "", 0),
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// Blocking until receive quit signal
		s := <-quit

		app.logger.Infoj(tlog.JSON{"message": "shutting down server", "signal": s.String()})

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
		app.logger.Infoj(tlog.JSON{"message": "completing background tasks", "addr": srv.Addr})

		app.wg.Wait()
		shutdownError <- nil
	}()

	app.logger.Infoj(tlog.JSON{"message": "starting server", "addr": srv.Addr, "env": app.config.Env})

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Info(tlog.JSON{"message": "server stopped", "addr": srv.Addr})
	return nil
}
