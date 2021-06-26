package handler

import (
	"chat_server/echange"
	"fmt"

	"github.com/gobwas/ws"
	"github.com/labstack/echo/v4"
)

type NoticeHandler struct{}

func NewNoticeHandler() *NoticeHandler {
	echange.EpollerServer()
	return &NoticeHandler{}
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func CreateErrorResponse() *Response {
	return &Response{
		Code: 401,
		Msg:  "auth err",
		Data: nil,
	}
}

func CreateSuccessResponse(data interface{}) *Response {
	return &Response{
		Code: 0,
		Msg:  "ok",
		Data: data,
	}
}

func CreateServerErrRes(msg string) *Response {
	return &Response{
		Code: 500,
		Msg:  msg,
	}
}

func (notice *NoticeHandler) WsHandler(c echo.Context) error {
	userID := c.QueryParam("userID")
	if userID == "" {
		return c.JSON(200, CreateErrorResponse())
	}
	conn, _, _, err := ws.UpgradeHTTP(c.Request(), c.Response())
	if err != nil {
		fmt.Println(err)
	}
	if fd, err := echange.GetSocketPusher().Add(echange.NewSocket(userID, conn)); err != nil {
		return c.JSON(200, CreateErrorResponse())
	} else {
		fmt.Println(fd, userID, ":user connected")
	}
	return c.JSON(200, CreateSuccessResponse(nil))
}
