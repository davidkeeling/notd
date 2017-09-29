package notd

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

var (
	homeTpl = template.Must(template.New("index.html").ParseFiles("views/index.html"))
)

func init() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", homeHandler).Methods("GET")
	rtr.HandleFunc("/login", loginHandler).Methods("GET")
	rtr.HandleFunc("/bipedaler/{userID}", bipedalerHandler).Methods("GET")
	rtr.HandleFunc("/blog", blogHandler).Methods("GET")
	rtr.HandleFunc("/blog/post", writeBlogPostHandler).Methods("GET")
	rtr.HandleFunc("/blog/post/{postID}", blogPostHandler).Methods("GET")
	rtr.HandleFunc("/blog/post/create", blogPostCreateHandler).Methods("POST")
	rtr.HandleFunc("/blog/post/edit/{postID}", blogPostEditHandler).Methods("GET")
	rtr.HandleFunc("/blog/post/update/{postID}", blogPostUpdateHandler).Methods("POST")
	rtr.HandleFunc("/blog/post/delete/{postID}", deleteBlogPostHandler).Methods("POST")
	rtr.HandleFunc("/media", mediaHandler).Methods("GET")
	rtr.HandleFunc("/media/upload", uploadMediaHandler).Methods("POST")
	rtr.HandleFunc("/media/delete/{mediaID}", deleteMediaHandler).Methods("POST")
	rtr.HandleFunc("/media/{mediaID}", serveBlogImageHandler).Methods("GET")

	http.Handle("/", rtr)
}

type page struct {
	Head      template.HTML
	Body      template.HTML
	Scale     float64
	Tapestry  dimensions
	Bandmates []bandmate
	IsAdmin   bool
}

type dimensions struct {
	Width  int
	Height int
}

type bandmate struct {
	Name   string `json:"name"`
	Left   int    `json:"left"`
	Top    int    `json:"top"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

var (
	bandmates = []bandmate{
		bandmate{
			Name:   "dave",
			Left:   602,
			Top:    308,
			Width:  121,
			Height: 325,
		},
		bandmate{
			Name:   "taylor",
			Left:   718,
			Top:    310,
			Width:  152,
			Height: 335,
		},
		bandmate{
			Name:   "jack",
			Left:   840,
			Top:    278,
			Width:  193,
			Height: 370,
		},
		bandmate{
			Name:   "ben",
			Left:   1060,
			Top:    280,
			Width:  156,
			Height: 339,
		},
		bandmate{
			Name:   "grant",
			Left:   1200,
			Top:    282,
			Width:  142,
			Height: 340,
		},
	}
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var p page
	bipedaler := userInfo(c)
	if bipedaler != nil {
		p.IsAdmin = bipedaler.IsAdmin
	}
	p.Bandmates = bandmates
	p.Tapestry = dimensions{
		Width:  2048,
		Height: 906,
	}
	err := homeTpl.Execute(w, p)
	if err != nil {
		log.Errorf(c, "%s", err)
	}
}
