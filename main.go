package main

import (
	"crypto/tls"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"
	"web-cemuhook/internal/CemuHook"
	"web-cemuhook/utils"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:embed templates tls
var staticFS embed.FS

var udpServer *CemuHook.UdpServer

var HttpsAddr = "0.0.0.0:8443"

var WebSocketUpgrader = websocket.Upgrader{}

func HomePage(w http.ResponseWriter, req *http.Request) {
	tmpl, _ := template.ParseFS(staticFS, "templates/home.html")
	tmpl.Execute(w, nil)
}

func WebSocketPage(w http.ResponseWriter, req *http.Request) {
	connection, err := WebSocketUpgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal().Err(err)
	}
	defer connection.Close()
	for {
		_, message, err := connection.ReadMessage()
		if err != nil {
			log.Warn().Err(err).Msg("WebSocket Disconnect")
			break
		}

		data := &utils.MessageData{}
		json.Unmarshal(message, data)
		log.Debug().RawJSON("json", message).Msg("WebSocket Received")
		udpServer.SendControllerData(data)
	}
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Stamp, NoColor: false})

	udpServer = CemuHook.NewUdpServer()
	go udpServer.StartReceive()

	log.Info().Msgf("Https server running at %s", HttpsAddr)
	log.Info().Msgf("Http running at %s", "0.0.0.0:8000")

	mux := http.NewServeMux()
	mux.HandleFunc("/", HomePage)
	mux.HandleFunc("/ws", WebSocketPage)
	serverCrt, _ := staticFS.ReadFile("tls/server.crt")
	serverKey, _ := staticFS.ReadFile("tls/server.key")
	cert, _ := tls.X509KeyPair(serverCrt, serverKey)
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
	server := &http.Server{
		Addr:         HttpsAddr,
		TLSConfig:    tlsConfig,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
		Handler:      mux,
	}

	go http.ListenAndServe(":8000", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		target := "https://" + req.Host + req.URL.Path
		target = strings.Replace(target, ":8000", ":8443", 1)
		fmt.Println(target)
		http.Redirect(w, req, target, http.StatusMovedPermanently)
	}))
	err := server.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatal().Err(err)
	}
}
