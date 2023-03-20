package app

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"websocket/domain"
	"websocket/service"
	"websocket/utils/auth"

	"github.com/gin-gonic/gin"
)

var HTTPS bool = false

// 環境変数チェック
func sanityCheck() {
	if os.Getenv("SERVER_ADDRESS") == "" || os.Getenv("SERVER_PORT") == "" {
		log.Fatal("Environment variable not defined....")
	}
	if os.Getenv("HTTPS") == "true" {
		HTTPS = true
	} else {
		HTTPS = false
	}
}

// TLS情報の作成
func NewTLSConfig(address string) (*tls.Config, error) {
	serverTLSConfig, err := auth.SetupTLSConfig(auth.TLSConfig{
		CertFile:      auth.ServerCertFile,
		KeyFile:       auth.ServerKeyFile,
		CAFile:        auth.CAFile,
		ServerAddress: address,
	})
	if err != nil {
		return nil, err
	}
	return serverTLSConfig, nil
}

func Run() {
	//sanityCheck()
	router := gin.Default()
	chatroom := domain.NewRoom()
	// DI
	// db := driver.NewDBDriver()
	// rp := domain.NewRepositoryDB(db)
	rp := domain.NewRepositoryStub()
	s := service.NewService(rp, chatroom)
	h := NewHandlers(s)
	// handling

	router.GET("/", h.index)
	router.GET("/ws/:id", h.WebSocketConnection)

	go chatroom.Run()

	// address := os.Getenv("SERVER_ADDRESS")
	// port := os.Getenv("SERVER_PORT")
	address := "localhost"
	port := "3000"

	if HTTPS == false {
		err := http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), router)
		if err != nil {
			log.Fatal(err)
		}
	} else if HTTPS == true {
		cfg, err := NewTLSConfig(address)
		if err != nil {
			panic(err)
		}
		srv := &http.Server{
			Addr:         ":" + port,
			Handler:      router,
			TLSConfig:    cfg,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		}
		log.Fatal(srv.ListenAndServeTLS(auth.ServerCertFile, auth.ServerKeyFile))
	}
}
