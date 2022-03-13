package messageService

import (
	// "encoding/json"
	"fmt"
	"log"
	"net/http"
	// "net/url"
	// "net/http/httputil"
	"strconv"
	// "github.com/satori/go.uuid"
)

type MessageService struct {
	FacPort  int
	LogPort  int
	MesPort  int
	Debug_on bool
}

// func NewServer - constructor for server instance
func NewServer(facport int, logport int, mesport int, debug_on bool) *MessageService {
	return &MessageService{FacPort: facport,
		LogPort:  logport,
		MesPort:  mesport,
		Debug_on: debug_on}
}

func (s *MessageService) MessageHandler(w http.ResponseWriter, r *http.Request) {
	if s.Debug_on && r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "Not implemented yet\n")
	case "POST":
		if s.Debug_on {
			fmt.Fprintf(w, "Nothing happaned\n")
		}

	default:
		if s.Debug_on {
			fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
		}
	}
}

func (s *MessageService) RunServer() {
	http.HandleFunc("/", s.MessageHandler)
	if s.Debug_on {
		log.Printf("Server started at port %d with Debug ON", s.MesPort)
	}
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(s.MesPort), nil))
}
