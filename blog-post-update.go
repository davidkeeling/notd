package notd

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func blogPostEditHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	user := userInfo(c)
	if user == nil || !user.IsAdmin {
		exit(w, "Unauthorized")
		return
	}
	params := mux.Vars(r)
	postID := params["postID"]
	blogPost, err := getBlogPost(c, postID)
	if err != nil {
		exit(w, err.Error())
		return
	}

	t, err := template.New("blog-edit.html").ParseFiles("views/blog-edit.html")
	if err != nil {
		exit(w, err.Error())
		return
	}

	blogImages, err := getBlogImages(c)
	if err != nil {
		exit(w, err.Error())
		return
	}

	err = t.Execute(w, struct {
		Post       *BlogPost
		User       *bipedaler
		BlogImages []*BlogImage
	}{
		blogPost,
		user,
		blogImages,
	})
	if err != nil {
		log.Errorf(c, "%s", err)
	}
}

func blogPostUpdateHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID := params["postID"]
	c := appengine.NewContext(r)
	user := userInfo(c)
	if user == nil || !user.IsAdmin {
		exit(w, "Unauthorized")
		return
	}
	newTitle := r.FormValue("title")
	newFeaturedImage := r.FormValue("featuredImage")
	newBody := r.FormValue("body")
	if newTitle == "" || newBody == "" || newFeaturedImage == "" {
		exit(w, "Need title, body, and featured image")
		return
	}
	key, err := datastore.DecodeKey(postID)
	if err != nil {
		exit(w, err.Error())
		return
	}
	blogPost := &BlogPost{}
	err = datastore.Get(c, key, blogPost)
	if err != nil || blogPost == nil {
		exit(w, "Not found")
		return
	}
	blogPost.Body = template.HTML(newBody)
	blogPost.Title = newTitle
	blogPost.FeaturedImageID = newFeaturedImage

	if r.FormValue("posted") != "" {
		milliseconds, err := strconv.Atoi(r.FormValue("posted"))
		if err == nil {
			created := time.Unix(int64(milliseconds/1000), 0)
			blogPost.Created = created
		} else {
			exit(w, err.Error())
			return
		}
	}

	_, err = datastore.Put(c, key, blogPost)
	if err != nil {
		exit(w, err.Error())
		return
	}
	http.Redirect(w, r, "/blog", http.StatusFound)
}
