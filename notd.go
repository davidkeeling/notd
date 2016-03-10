package notd

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func init() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", homeHandler).
		Methods("GET")

	rtr.HandleFunc("/login", loginHandler).
		Methods("GET")

	rtr.HandleFunc("/bipedaler/{userID}", bipedalerHandler).
		Methods("GET")

	rtr.HandleFunc("/blog", blogHandler).
		Methods("GET")

	rtr.HandleFunc("/blog/post", writeBlogPostHandler).
		Methods("GET")

	rtr.HandleFunc("/blog/post/{postID}", blogPostHandler).
		Methods("GET")

	rtr.HandleFunc("/blog/post/create", blogPostCreateHandler).
		Methods("POST")

	rtr.HandleFunc("/blog/post/edit/{postID}", blogPostEditHandler).
		Methods("GET")

	rtr.HandleFunc("/blog/post/update/{postID}", blogPostUpdateHandler).
		Methods("POST")

	rtr.HandleFunc("/blog/post/delete/{postID}", deleteBlogPostHandler).
		Methods("POST")

	rtr.HandleFunc("/media", mediaHandler).
		Methods("GET")

	rtr.HandleFunc("/media/upload", uploadMediaHandler).
		Methods("POST")

	rtr.HandleFunc("/media/delete/{mediaID}", deleteMediaHandler).
		Methods("POST")

	rtr.HandleFunc("/media/{mediaID}", serveBlogImageHandler).
		Methods("GET")

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
	// dimensions
	Name   string
	Left   int
	Top    int
	Width  int
	Height int
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

func mult(x int, y float64) float64 {
	return float64(x) * y
}

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
	p.Scale = .5
	funcMap := template.FuncMap{"mult": mult}
	t, err := template.New("index.html").Funcs(funcMap).ParseFiles("views/index.html")
	if err != nil {
		exit(w, fmt.Sprintf("Error parsing template: %s", err.Error()))
		return
	}
	err = t.Execute(w, p)
	if err != nil {
		log.Errorf(c, "%s", err)
	}
}
