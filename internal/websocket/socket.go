package websocket

import (
	"github.com/Seann-Moser/BaseGoAPI/pkg/response"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

// https://github.com/gorilla/websocket/tree/master/examples/echo

type WebSocketHandler struct {
	logger   *zap.Logger
	upgrader websocket.Upgrader
	Resp     *response.Response
}

const (
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	writeWait = 10 * time.Second
)

type WSHandler func(w http.ResponseWriter, r *http.Request)

func (ws *WebSocketHandler) Connect(readChan, writeChan chan interface{}) WSHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := ws.ConnectToWebsocket(w, r)
		if err != nil {
			ws.Resp.Error(w, err, http.StatusBadRequest, "failed to read ws")
			return
		}
		wg := &sync.WaitGroup{}
		wg.Add(2)
		go func() {
			ws.WriteData(c, writeChan, wg)
			close(readChan)
		}()
		go func() {
			ws.ReadData(c, readChan, wg)
			close(writeChan)
		}()
		wg.Wait()
	}
}

func (ws *WebSocketHandler) ConnectToWebsocket(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	c, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		ws.logger.Error("failed to upgrade connection to ws", zap.Error(err))
		return nil, err
	}
	return c, err
}

func (ws *WebSocketHandler) WriteData(c *websocket.Conn, dataChan chan interface{}, wg *sync.WaitGroup) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.Close()
		wg.Done()
	}()
	for {
		select {
		case message, ok := <-dataChan:
			_ = c.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.WriteMessage(websocket.CloseMessage, []byte{})
				close(dataChan)
				return
			}
			err := c.WriteJSON(message)
			if err != nil {
				ws.logger.Error("failed to write json to ws connection", zap.Error(err))
				_ = c.WriteMessage(websocket.CloseMessage, []byte("failed to write json to ws connection"))
				close(dataChan)
				continue
			}
		case <-ticker.C:
			_ = c.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.WriteMessage(websocket.PingMessage, nil); err != nil {
				close(dataChan)
				return
			}
		}
	}
}

func (ws *WebSocketHandler) ReadData(c *websocket.Conn, dataChan chan interface{}, wg *sync.WaitGroup) {
	defer func() {
		_ = c.Close()
		wg.Done()
	}()
	c.SetPongHandler(func(string) error { _ = c.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var msg interface{}
		err := c.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				ws.logger.Error("failed reading data from ws", zap.Error(err))
			}
			break
		}
		dataChan <- msg
	}
}
