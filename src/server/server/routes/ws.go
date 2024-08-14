package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"server_go/src/service/ws"
) // websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index

func WS() fiber.Handler {
	wsClients := make(map[string]string)
	ws.On(ws.EventDisconnect, func(ep *ws.EventPayload) {
		// Remove the user from the local wsClients
		delete(wsClients, ep.Kws.GetStringAttribute("user_id"))
		fmt.Println(fmt.Sprintf("Disconnection event - User: %s", ep.Kws.GetStringAttribute("user_id")))
	})
	ws.On(ws.EventConnect, func(ep *ws.EventPayload) {
		fmt.Println(fmt.Sprintf("Connection event 1 - User: %s", ep.Kws.GetStringAttribute("user_id")))
	})
	return ws.New(func(kws *ws.Websocket) {
		// Retrieve the user id from endpoint
		userId := kws.Query("id")

		// Add the connection to the list of the connected wsClients
		// The UUID is generated randomly and is the key that allow
		// ws to manage Emit/EmitTo/Broadcast
		wsClients[userId] = kws.UUID
		kws.SetAttribute("user_id", userId)
		//Broadcast to all the connected users the newcomer
		kws.Broadcast([]byte(fmt.Sprintf("New user connected: %s and UUID: %s", userId, kws.UUID)), true)
		//Write welcome message
		kws.Emit([]byte(fmt.Sprintf("Hello user: %s with UUID: %s", userId, kws.UUID)))
	})
}
