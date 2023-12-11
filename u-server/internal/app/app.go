package app

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pullya/unique_server/u-server/internal/config"
	"github.com/pullya/unique_server/u-server/internal/model"
	"github.com/pullya/unique_server/u-server/internal/repository"
	log "github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type App struct {
	IpRepo  repository.IIpRepo
	NumRepo repository.INumRepo

	wg *sync.WaitGroup
}

func InitApp() App {
	ipRepo := repository.NewIpRepo()
	numRepo := repository.NewNumRepo()

	wg := sync.WaitGroup{}

	return App{
		IpRepo:  &ipRepo,
		NumRepo: &numRepo,
		wg:      &wg,
	}
}

func (a *App) Run(ctx context.Context) error {
	http.HandleFunc(config.Endpoint, func(w http.ResponseWriter, r *http.Request) {
		a.wg.Add(1)
		go a.handleConnection(ctx, w, r)
	})

	err := http.ListenAndServe(preparePort(config.WsPort), nil)
	if err != nil {
		return err
	}
	log.WithField("service", config.ServiceName).Infof("WebSocket server is running on :%d", config.WsPort)

	a.wg.Wait()
	return nil
}

func (a *App) handleConnection(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	defer a.wg.Done()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()
	log.WithField("service", config.ServiceName).Info("Client connected")

	ip := conn.RemoteAddr().String()
	if !a.IpRepo.IsNewIp(ctx, ip) {
		log.WithField("service", config.ServiceName).Errorf("Not a new ip addess %s. Close connection", ip)
		return
	}

	num := a.NumRepo.GenUniqueNum(ctx)
	msg := model.PrepareMessage("OK", num.String())

	if err := conn.WriteMessage(websocket.TextMessage, msg.AsJson()); err != nil {
		log.WithField("service", config.ServiceName).Errorf("Failed to send response: %v", err)
		return
	}
}

func preparePort(port int) string {
	return fmt.Sprintf(":%d", port)
}
