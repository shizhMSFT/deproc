package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/urfave/cli/v2"
)

var serveCommand = &cli.Command{
	Name:      "serve",
	Usage:     "Listen and serve HTTPS",
	ArgsUsage: "<repository_name>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:      "key",
			Aliases:   []string{"k"},
			Usage:     "TLS key file",
			TakesFile: true,
			Required:  true,
		},
		&cli.StringFlag{
			Name:      "cert",
			Aliases:   []string{"c"},
			Usage:     "TLS cert file",
			TakesFile: true,
			Required:  true,
		},
		&cli.StringFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Usage:   "Listen port",
			Value:   ":443",
		},
	},
	Action: runServe,
}

func runServe(ctx *cli.Context) error {
	port := ctx.String("port")
	if !strings.Contains(port, ":") {
		port = ":" + port
	}
	return http.ListenAndServeTLS(port, ctx.String("cert"), ctx.String("key"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		url := fmt.Sprintf("https://%s%s", r.Host, r.URL.String())
		log := log.WithFields(log.Fields{
			"http_method": r.Method,
			"http_url":    url,
			"http_range":  r.Header.Get("Range"),
		})

		req, err := http.NewRequestWithContext(r.Context(), r.Method, url, r.Body)
		if err != nil {
			log = log.WithDuration(time.Since(startTime))
			log = log.WithError(err)
			log.Error("request failure")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		req.Header = r.Header.Clone()
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log = log.WithDuration(time.Since(startTime))
			log = log.WithError(err)
			log.Error("upstream failure")
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		log = log.WithField("http_status", resp.StatusCode)
		header := w.Header()
		for k, v := range resp.Header {
			header[k] = v
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log = log.WithDuration(time.Since(startTime))
			log = log.WithError(err)
			log.Error("upstream failure")
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		w.WriteHeader(resp.StatusCode)
		w.Write(body)
		log = log.WithDuration(time.Since(startTime))
		log.Info("request proxied")
	}))
}
