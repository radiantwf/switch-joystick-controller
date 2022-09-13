package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	controller "github.com/radiantwf/switch-joystick-controller"

	"github.com/gorilla/websocket"
	"github.com/naoina/denco"
)

// HTTPService 定义
type HTTPService struct {
	addr       string
	upgrader   *websocket.Upgrader
	wsMessages chan []byte
	controller *controller.JoyStick
}

func (s *HTTPService) Init() (err error) {
	port := 8888
	s.addr = fmt.Sprintf(":%d", port)
	s.upgrader = &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	s.wsMessages = make(chan []byte)
	s.controller = controller.NewJoyStick()
	return
}

func (s *HTTPService) Run() (err error) {
	err = s.controller.Open()
	if err != nil {
		return
	}
	defer s.controller.Close()
	go func() {
		for message := range s.wsMessages {
			s.controller.SyncSendKey(controller.NewJoyStickInput(string(message)), 0)
		}
	}()

	mux := denco.NewMux()
	handler, err := mux.Build([]denco.Handler{
		mux.GET("/", s.serveFile),
		mux.GET("/controller", s.switchController),
		mux.GET("/assets/*", s.serveFile),
	})
	if err != nil {
		return
	}
	server := http.Server{Addr: s.addr,
		Handler: handler}
	log.Fatal(server.ListenAndServe())
	return
}

func (s *HTTPService) serveFile(w http.ResponseWriter, r *http.Request, params denco.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin,Accept,Cache-Control,Pragma,Content-Type,Authorization")
	fileName := ""
	if r.URL.Path[1:] == "" {
		fileName = "resources/web/index.html"
	} else {

		filepath := path.Join("resources/web", r.URL.Path[1:])
		s, err := os.Stat(filepath)
		if err == nil {
			if !s.IsDir() {
				fileName = filepath
			}
		}
	}
	if fileName == "" {
		fileName = path.Join("resources/web", "404.html")
	}
	http.ServeFile(w, r, fileName)
}

func (s *HTTPService) switchController(w http.ResponseWriter, r *http.Request, params denco.Params) {
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			switch e := err.(type) {
			case *websocket.CloseError:
				log.Println("WebSocket closed:", e.Code)
			default:
				log.Println("WebSocket read error:", e)
			}
			break
		}
		s.wsMessages <- message
	}
}
