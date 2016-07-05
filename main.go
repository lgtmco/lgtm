package main

import (
	"net/http"
	"time"

	"github.com/lgtmco/lgtm/router"
	"github.com/lgtmco/lgtm/router/middleware"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/ianschenck/envflag"
	_ "github.com/joho/godotenv/autoload"

	_ "github.com/lgtmco/lgtm/approval/org"
)

var (
	addr = envflag.String("SERVER_ADDR", ":8000", "")
	cert = envflag.String("SERVER_CERT", "", "")
	key  = envflag.String("SERVER_KEY", "", "")

	debug = envflag.Bool("DEBUG", false, "")
)

func main() {
	envflag.Parse()

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.WarnLevel)
	}

	handler := router.Load(
		ginrus.Ginrus(logrus.StandardLogger(), time.RFC3339, true),
		middleware.Version,
		middleware.Store(),
		middleware.Remote(),
		middleware.Cache(),
	)

	if *cert != "" {
		logrus.Fatal(
			http.ListenAndServeTLS(*addr, *cert, *key, handler),
		)
	} else {
		logrus.Fatal(
			http.ListenAndServe(*addr, handler),
		)
	}
}
