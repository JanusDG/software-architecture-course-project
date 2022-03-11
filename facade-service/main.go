package main

import (
	// "fmt"
	// "log"
	// "net/http"

	"github.com/JanusDG/software-architecture-course-project/facade-service/config"
	"github.com/JanusDG/software-architecture-course-project/facade-service/facadeServer"
)

func main() {
	var cfg = config.GetConf()
	var self_port = cfg.FacadeServer.Port
	// var loggingService_port = cfg.LoggingService.Port
	// var messageService_port = cfg.MessageService.Port
	var debug = cfg.FacadeServer.DEBUG_ON

	serv := facadeServer.NewServer(self_port, debug)

	serv.RunServer()

}
