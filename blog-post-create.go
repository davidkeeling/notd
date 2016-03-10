package notd

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func writeBlogPostHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	user := userInfo(c)
	if user == nil || !user.IsAdmin {
		exit(w, "Unauthorized")
		return
	}
	t, err := template.New("blog-write.html").ParseFiles("views/blog-write.html")
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
		User       *bipedaler
		BlogImages []*BlogImage
	}{
		user,
		blogImages,
	})
	if err != nil {
		log.Errorf(c, "%s", err)
	}
}

func blogPostCreateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	c := appengine.NewContext(r)
	user := userInfo(c)
	if user == nil || !user.IsAdmin {
		exit(w, "Unauthorized")
		return
	}
	newBlogPost := &BlogPost{
		Created:         time.Now(),
		AuthorID:        user.UserID,
		Title:           r.FormValue("title"),
		Body:            template.HTML(r.FormValue("body")),
		FeaturedImageID: r.FormValue("featuredImage"),
	}
	if newBlogPost.Title == "" || newBlogPost.Body == "" || newBlogPost.FeaturedImageID == "" || newBlogPost.AuthorID == "" {
		exit(w, "Need title, body, featured image and authorID")
		return
	}
	if r.FormValue("posted") != "" {
		milliseconds, err := strconv.Atoi(r.FormValue("posted"))
		if err == nil {
			created := time.Unix(int64(milliseconds/1000), 0)
			newBlogPost.Created = created
		} else {
			exit(w, err.Error())
			return
		}
	}
	key := datastore.NewKey(c, "BlogPost", fmt.Sprintf("%s%s", newBlogPost.Title, newBlogPost.Created.String()), 0, nil)
	newBlogPost.PostID = key.Encode()
	_, err := datastore.Put(c, key, newBlogPost)
	if err != nil {
		exit(w, err.Error())
		return
	}
	http.Redirect(w, r, "/blog", http.StatusFound)
}
