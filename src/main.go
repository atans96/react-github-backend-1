package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	mainServer "server_go/src/server/server"
	serverCORS "server_go/src/server/server_cors_proxy"
	serverImg "server_go/src/server/server_image"
	"server_go/src/service"
)

// init is invoked before main()
func init() {
	service.ReadFiles()
	_ = service.SessionStore.NewSessionStore()

	//localIP := service.GetLocalIP()
	err := os.Setenv("DATABASE", fmt.Sprintf("mongodb://%[1]v:27017,%[1]v:27018,%[1]v:27019/github?replicaSet=re0", "192.168.100.192"))
	if err != nil {
		panic(err)
	}
	err = os.Setenv("REDIS_LOCAL_ENDPOINT", fmt.Sprintf("%s", "192.168.100.192"))
	if err != nil {
		return
	}
	if err != nil {
		panic(err)
	}
	dir, err := filepath.Abs("./")
	if err != nil {
		panic(err)
	}
	err = os.Chdir(filepath.Join(dir, "src"))
	if err != nil {
		panic(err)
	}
	err = godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	_ = service.Mongo.NewDatastore()
	cert, err := tls.X509KeyPair([]byte(os.Getenv("localhost.crt")), []byte(os.Getenv("localhost.decrypted.key")))
	proxy.WithTlsConfig(&tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
	})
}
func main() {
	go func() {
		serverCORS.Run()
	}()
	go func() {
		serverImg.Run()
	}()
	mainServer.Run()
}
