package main

import (
	// "fmt"
	// "log"
	// "net/http"

	"github.com/JanusDG/software-architecture-course-project/config"
	"github.com/JanusDG/software-architecture-course-project/messages-service/messageService"
)

func main() {
	var cfg = config.GetConf()
	var facadeServer_port = cfg.FacadeServer.Port
	var loggingService_port = cfg.LoggingService.Port
	var messageService_port = cfg.MessageService.Port
	var debug = cfg.MessageService.DEBUG_ON

	serv := messageService.NewServer(facadeServer_port, loggingService_port, messageService_port, debug)

	serv.RunServer()

}
