package server

type WsServer struct {
	Port int
}

func NewWsServer(port int) WsServer {
	return WsServer{
		Port: port,
	}
}

type IServer interface {
}
