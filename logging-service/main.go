package main

import (
	// "fmt"
	// "log"
	// "net/http"

	"github.com/JanusDG/software-architecture-course-project/config"
	"github.com/JanusDG/software-architecture-course-project/logging-service/loggingService"
)

func main() {
	var cfg = config.GetConf()
	var facadeServer_port = cfg.FacadeServer.Port
	var loggingService_port = cfg.LoggingService.Port
	var messageService_port = cfg.MessageService.Port
	var debug = cfg.LoggingService.DEBUG_ON

	serv := loggingService.NewServer(facadeServer_port, loggingService_port, messageService_port, debug)

	serv.RunServer()

}
