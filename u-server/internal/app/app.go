package app

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pullya/unique_server/u-server/internal/config"
	"github.com/pullya/unique_server/u-server/internal/model"
	"github.com/pullya/unique_server/u-server/internal/storage"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type App struct {
	IpStorage  storage.IpStorageer
	NumStorage storage.NumStorageer

	endpoint string
	port     int

	wg *sync.WaitGroup
}

func InitApp() App {
	ipStorage := storage.NewIpStorage()
	numStorage := storage.NewNumStorage()

	wg := sync.WaitGroup{}

	return App{
		IpStorage:  &ipStorage,
		NumStorage: &numStorage,
		endpoint:   config.Config.Endpoint,
		port:       config.Config.Port,
		wg:         &wg,
	}
}

func (a *App) Run(ctx context.Context) error {
	http.HandleFunc(a.endpoint, func(w http.ResponseWriter, r *http.Request) {
		a.handleConnection(ctx, w, r)
	})

	config.Logger.Debugf("WebSocket server is running on :%d", a.port)

	err := http.ListenAndServe(preparePort(a.port), nil)
	if err != nil {
		config.Logger.Errorf("Failed to run server: %v", err)
		return err
	}

	return nil
}

func (a *App) handleConnection(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		config.Logger.Errorf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	if err = conn.SetReadDeadline(time.Now().Add(time.Duration(config.Config.Timeout) * time.Millisecond)); err != nil {
		config.Logger.Errorf("Error while setting timeout: %v", err)
		return
	}

	// Получаем IP-адрес нового соединения
	ip := conn.RemoteAddr().String()
	config.Logger.Infof("Client connected from ip: %s", ip)

	// Если установлен strictMode=on, то проверяем только IP-адрес без учета номера порта.
	if config.Config.StrictMode == "on" {
		ip = parseIp(ip)
	}
	// Проверяем, было ли уже ранее соединение с этого IP-адреса. Если было, то - отказ
	if !a.IpStorage.IsNewIp(ctx, ip) {
		config.Logger.Errorf("Not a new ip addess %s. Close connection", ip)
		return
	}

	// Читаем входящий запрос от клиента
	mt, message, err := conn.ReadMessage()
	if err != nil || mt == websocket.CloseMessage {
		return
	}
	config.Logger.Infof("Message from client received: %s", message)

	// Генерируем уникальный big.Int для ответа
	num := a.NumStorage.GenUniqueNum(ctx)

	// Формируем ответ в формате JSON
	msg := model.PrepareMessage(model.MessageTypeResponse, num.String())

	// Отправляем ответное сообщение
	if err := conn.WriteMessage(websocket.TextMessage, msg.AsJson()); err != nil {
		config.Logger.Errorf("Failed to send response: %v", err)
		return
	}
}

func preparePort(port int) string {
	return fmt.Sprintf(":%d", port)
}

func parseIp(ip string) string {
	return strings.Split(ip, ":")[0]
}
