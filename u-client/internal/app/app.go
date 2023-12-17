package app

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pullya/unique_server/u-client/internal/config"
	"github.com/pullya/unique_server/u-client/internal/model"
)

type App struct {
	address string
}

func InitApp() App {
	return App{
		address: buildAddress(config.Config.Endpoint, config.Config.Port),
	}
}

func (a *App) Run(ctx context.Context) error {
	wg := sync.WaitGroup{}
	for i := 0; i < config.Config.ConnCnt; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Устанавливаем соединение с сервером
			dialer := websocket.DefaultDialer
			conn, _, err := dialer.Dial(a.address, nil)
			if err != nil {
				config.Logger.Errorf("Error connecting to WebSocket server: %v", err)
				return
			}
			defer conn.Close()

			if err = conn.SetReadDeadline(time.Now().Add(time.Duration(config.Config.Timeout) * time.Millisecond)); err != nil {
				config.Logger.Errorf("Error while setting timeout: %v", err)
				return
			}

			// Формируем запрос в формате JSON
			message := model.PrepareMessage(model.MessageTypeRequest, "").AsJson()
			// Отправляем запрос
			err = conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				config.Logger.Errorf("Error writing message: %v", err)
				return
			}

			// Читаем полученный ответ
			mt, message, err := conn.ReadMessage()
			if err != nil || mt == websocket.CloseMessage {
				config.Logger.Errorf("Error reading message: %v", err)
				return
			}
			config.Logger.Infof("Received message: %s\n", message)
		}()
		time.Sleep(time.Millisecond * time.Duration(config.Config.ConnInterval))
	}

	wg.Wait()
	return nil
}

func buildAddress(endpoint string, port int) string {
	return fmt.Sprintf("ws://u_server:%d%s", port, endpoint)
}
