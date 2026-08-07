package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade(r, w)
	if err != nil {
		log.Println(err)
		return
	}
	//Handle websocket messages and keep the connection alive
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Println("websocket read error:", err)
			}
			break
		}
		//Handle the message and send the response back to the client
		if messageType == websocket.TextMessage {
			textMessage := string(message)
			response := "Received message: " + textMessage
			if err := conn.WriteMessage(websocket.TextMessage, []byte(response)); err != nil {
				log.Println("websocket write error:", err)
				break
			}
		} else if messageType == websocket.CloseMessage {
			break
		}
	}
}

func upgrade(r *http.Request, w http.ResponseWriter) (any, any) {
	panic("unimplemented")
}

func main() {
	//Initialize the api Gateway
	apiGateway := openapi.NewAPIGateway()

	//Define routes
	routes := []Route{
		{"POST", "/users", "/services/user-service",
			func(w http.ResponseWriter, r *http.Request) {
				apiGateway.proxyUserService(r, w)
			},
		},
	}
	apiGateway.SetRoutes(routes)

	//start the server
	http.ListenAndServe(":8080", apiGateway.Server())
}
