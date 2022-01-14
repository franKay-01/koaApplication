package main

import (
	"context"
	"firstProject/api"
	"firstProject/config"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/itrepablik/itrlog"
	"github.com/itrepablik/sakto"
)

var CurrentLocalTime = sakto.GetCurDT(time.Now(), config.SiteTimeZone)

func main() {
	os.Setenv("TZ", config.SiteTimeZone)
	fmt.Println("Starting the web servers at ", CurrentLocalTime)

	var dir string
	var wait time.Duration

	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.StringVar(&dir, "dir", "static", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()

	r := mux.NewRouter()

	r.Use(mux.CORSMethodMiddleware(r))

	api.MainRouters(r) // URLs for the main app.

	srv := &http.Server{
		Addr: "127.0.0.1:8001",

		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	go func() {
		msg := `Web server started at `
		fmt.Println(msg, CurrentLocalTime)
		itrlog.Info("Web server started at ", CurrentLocalTime)
		if err := srv.ListenAndServe(); err != nil {
			itrlog.Error(err)
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)
	fmt.Println("Shutdown web server at " + CurrentLocalTime.String())
	itrlog.Warn("Server has been shutdown at ", CurrentLocalTime.String())
	os.Exit(0)
}
