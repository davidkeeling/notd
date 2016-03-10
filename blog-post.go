package notd

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func blogPostHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	urlVars := mux.Vars(r)
	blogPost, err := getBlogPost(c, urlVars["postID"])
	if err != nil {
		exit(w, err.Error())
		return
	}

	t, err := template.New("blog-post.html").ParseFiles("views/blog-post.html")
	if err != nil {
		exit(w, err.Error())
		return
	}

	user := userInfo(c)

	err = t.Execute(w, struct {
		Post *BlogPost
		User *bipedaler
	}{
		blogPost,
		user,
	})
	if err != nil {
		log.Errorf(c, "%s", err)
	}
}

func deleteBlogPostHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	user := userInfo(c)
	if user == nil || !user.IsAdmin {
		exit(w, "Unauthorized")
		return
	}
	params := mux.Vars(r)
	postID := params["postID"]
	if postID == "" {
		exit(w, "Not found")
		return
	}
	key, err := datastore.DecodeKey(postID)
	if err != nil {
		exit(w, err.Error())
		return
	}
	err = datastore.Delete(c, key)
	if err != nil {
		log.Errorf(c, "%s", err)
	}
	http.Redirect(w, r, "/blog", http.StatusFound)
}
