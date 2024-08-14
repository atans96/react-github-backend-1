package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"os"
	"server_go/src/service"
)

type Starts struct {
	Cookie     []byte
	jsonStruct jsonStruct
}
type jsonStruct struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

func (Starts *Starts) Start(c *fiber.Ctx) error {
	store, err := service.SessionStore.DB.Get(c)
	if err != nil {
		err := c.Status(fiber.StatusForbidden).SendString("Something wrong with session")
		if err != nil {
			panic(err)
		}
	}
	if c.Query("end") == "true" {
		url := os.Getenv("UWEBSOCKET_HOST") + ":" + os.Getenv("UWEBSOCKET_PORT") + "/server_uwebsocket/logout?end=true&username=" + c.Query("username")
		makeNew, err := http.NewRequest("GET", url, nil)
		if err != nil {
			panic(err)
		}
		makeNew.Header.Set("User-Agent", "Chrome: 91")
		makeNew.Header.Set("Origin", "https://"+os.Getenv("GOLANG_HOST")+":"+os.Getenv("GOLANG_PORT"))
		request := service.MakeRequest{Req: makeNew}
		request.Request()
		request.Clear()
		c.Append(fiber.HeaderClearSiteData, "\"cache\", \"cookies\", \"storage\"")
		return c.SendStatus(fiber.StatusNoContent)
	}
	if Starts.Cookie = c.Request().Header.Cookie("session_id"); Starts.Cookie == nil {
		send := `{
			"data": "",
			"username": ""
		}`
		return c.SendString(send)
	}
	cookie := string(c.Request().Header.Cookie("session_id"))
	j := store.Get(cookie)
	if j == nil {
		c.Append(fiber.HeaderClearSiteData, "\"cache\", \"cookies\", \"storage\"")
		err := c.Status(fiber.StatusForbidden).SendString("Your session is invalid")
		if err != nil {
			panic(err)
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
	err = json.Unmarshal([]byte(j.(string)), &Starts.jsonStruct)
	if err != nil {
		panic(err)
	}
	rid := fmt.Sprintf(`{
		"data": %v,
		"username": "%v"
	}`, true, Starts.jsonStruct.Username)
	return c.SendString(rid)
}
