package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
)

type Response struct {
	Code    ResponseCode `json:"code"`
	Message string       `json:"message,omitempty"`
	Data    interface{}  `json:"data,omitempty"`
}

type ResponseWithRequestId struct {
	RequestId string `json:"request_id"`
	Response
}

type BuErrorResponse struct {
	HttpCode int
	*Response
}

func (res Response) Error() string {
	t, _ := json.Marshal(res)
	return string(t)
}

func NewResponse(code ResponseCode, msg string, data interface{}) *Response {
	return &Response{
		Code:    code,
		Message: msg,
		Data:    data,
	}
}

func NewResponseWithRequestId(requestId string, r *Response) *ResponseWithRequestId {
	return &ResponseWithRequestId{
		RequestId: requestId,
		Response: Response{
			Code:    r.Code,
			Message: r.Message,
			Data:    r.Data,
		},
	}
}

type ResponseCode int

const (
	CodeOk ResponseCode = iota
	CodeUnKnownReasonErr

	CodeError
	CodeInternalServer
	CodeBadRequest
	CodeNotFound
)

var (
	OK               = &Response{Code: CodeOk, Message: "success."}
	ErrUnKnownReason = &Response{Code: CodeUnKnownReasonErr, Message: "unknown reason."}

	ErrInternalServer = &Response{Code: CodeInternalServer, Message: "server internal error."}
	ErrBadRequest     = &Response{Code: CodeBadRequest, Message: "bad request."}
	ErrNotFound       = &Response{Code: CodeNotFound, Message: "object not found."}
)

type Gin struct {
	C *gin.Context
}

func (g *Gin) HTTPResponseOK(data interface{}) {
	requestId := g.C.Request.Header.Get("Kong-Request-ID")
	g.C.JSON(http.StatusOK, NewResponseWithRequestId(requestId, NewResponse(OK.Code, OK.Message, data)))
}
