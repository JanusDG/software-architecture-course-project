package facadeServer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/satori/go.uuid"
)

type FacadeServer struct {
	FacPort  int
	LogPort  int
	MesPort  int
	Debug_on bool
}

// func NewServer - constructor for server instance
func NewServer(facport int, logport int, mesport int, debug_on bool) *FacadeServer {
	return &FacadeServer{FacPort: facport, LogPort: logport, MesPort: mesport, Debug_on: debug_on}
}

func (s *FacadeServer) MessageHandler(w http.ResponseWriter, r *http.Request) {
	if s.Debug_on && r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		resp, err := http.Get("http://localhost:" + strconv.Itoa(s.LogPort))

		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}

		if s.Debug_on {
			fmt.Printf("Recieved from Logging server:\n%s\n", string(body))
		}

		resp1, err1 := http.Get("http://localhost:" + strconv.Itoa(s.MesPort))

		if err1 != nil {
			log.Fatal(err1)
		}

		defer resp1.Body.Close()

		body1, err1 := ioutil.ReadAll(resp1.Body)

		if err1 != nil {
			log.Fatal(err1)
		}

		if s.Debug_on {
			fmt.Printf("Recieved from Message server:\n%s\n", string(body1))
		}
		fmt.Fprintf(w, "\nLogging:\n%s\nMessage:\n%s\n", string(body), string(body1))

	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		message := r.FormValue("message")
		if s.Debug_on && message != "" {
			fmt.Fprintf(w, "Message recieved\n")
			fmt.Printf("Recieved message: %s\n\n", message)
		}

		var message_uuid = uuid.NewV1().String()

		data := url.Values{
			"uuid":    {message_uuid},
			"message": {message},
		}

		resp, err := http.PostForm("http://localhost:"+strconv.Itoa(s.LogPort), data)

		if err != nil {
			log.Fatal(err)
		}

		var res map[string]interface{}

		json.NewDecoder(resp.Body).Decode(&res)

	default:
		if s.Debug_on {
			fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
		}
	}
}

func (s *FacadeServer) RunServer() {
	http.HandleFunc("/", s.MessageHandler)
	if s.Debug_on {
		log.Printf("Server started at port %d with Debug ON", s.FacPort)
	}
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(s.FacPort), nil))
}
