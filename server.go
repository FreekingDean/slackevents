package slackevents

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var DefaultServer = &Server{
	URLVerificationHandler: DefaultURLVerificationHandler,
	CallbackHandler:        DefaultCallbackServer.Handler,
	ErrorHandler:           DefaultErrorHandler,
}

type Server struct {
	URLVerificationHandler func(*URLVerification) (string, error)
	CallbackHandler        CallbackHandler
	ErrorHandler           func(error) (int, string)
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, resp := server.handle(r.Body)
	w.WriteHeader(status)
	w.Write([]byte(resp))
}

func (server *Server) handle(data io.ReadCloser) (int, string) {
	defer data.Close()

	//validate request
	//return error

	decoder := json.NewDecoder(data)
	event := &Event{}
	err := decoder.Decode(event)
	if err != nil {
		return server.ErrorHandler(err)
	}

	switch parsedEvent := event.ParsedEvent.(type) {
	case *URLVerification:
		response, err := server.URLVerificationHandler(parsedEvent)
		if err != nil {
			return server.ErrorHandler(err)
		}

		return 200, response
	case *Callback:
		err = server.CallbackHandler(parsedEvent)
		if err != nil {
			return server.ErrorHandler(err)
		}

		return 200, ""
	}

	return server.ErrorHandler(fmt.Errorf("Couldnt handle event type %T", event.ParsedEvent))
}

func RegisterCallbackHandler(handler CallbackHandler) {
	DefaultServer.RegisterCallbackHandler(handler)
}

func (server *Server) RegisterCallbackHandler(handler CallbackHandler) {
	server.CallbackHandler = handler
}
