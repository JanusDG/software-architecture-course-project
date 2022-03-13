package loggingService

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

type LoggingService struct {
	FacPort        int
	LogPort        int
	MesPort        int
	Debug_on       bool
	DB_placeholder map[string]string
}

// func NewServer - constructor for server instance
func NewServer(facport int, logport int, mesport int, debug_on bool) *LoggingService {
	return &LoggingService{FacPort: facport,
		LogPort:        logport,
		MesPort:        mesport,
		Debug_on:       debug_on,
		DB_placeholder: make(map[string]string, 0)}
}

func (s *LoggingService) MessageHandler(w http.ResponseWriter, r *http.Request) {
	if s.Debug_on && r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		for i, k := range s.DB_placeholder {
			fmt.Fprintf(w, "id:%s m:%s\n", i, k)
		}
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		message := r.FormValue("message")
		uuid := r.FormValue("uuid")
		s.DB_placeholder[uuid] = message
		if s.Debug_on && message != "" {
			fmt.Printf("Recieved and logged message: %s with id=%s\n", message, uuid)
		}

	default:
		if s.Debug_on {
			fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
		}
	}
}

func (s *LoggingService) RunServer() {
	http.HandleFunc("/", s.MessageHandler)
	if s.Debug_on {
		log.Printf("Server started at port %d with Debug ON", s.LogPort)
	}
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(s.LogPort), nil))
}
