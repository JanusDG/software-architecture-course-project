package facadeServer

import (
	"fmt"
	"log"
	"net/http"
	// "net/http/httputil"
	"strconv"
	// "github.com/satori/go.uuid"
)

type FacadeServer struct {
	Port     int
	Debug_on bool
}

// func Init - initializer for server instance
func (s *FacadeServer) Init(port int, DEBUG_ON bool) {
	s.Port = port
	s.Debug_on = DEBUG_ON
}

// func NewServer - constructor for server instance
func NewServer(port int, debug_on bool) *FacadeServer {
	return &FacadeServer{Port: port, Debug_on: debug_on}
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Hello, there\n")
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "form.html")
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		name := r.FormValue("name")
		fmt.Fprintf(w, "Name = %s\n", name)
		// byts, _ := httputil.DumpRequest(r, true)
		// fmt.Println(string(byts))
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func (s *FacadeServer) RunServer() {
	http.HandleFunc("/", HelloHandler)
	if s.Debug_on {
		log.Printf("Server started at port %d", s.Port)
	}
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(s.Port), nil))
}
