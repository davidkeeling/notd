package notd

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type bipedaler struct {
	Name    string
	Email   string
	UserID  string
	IsAdmin bool
}

func bipedalerHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	params := mux.Vars(r)
	userID := params["userID"]
	key, err := datastore.DecodeKey(userID)
	if err != nil || key == nil {
		exit(w, "User not found")
		return
	}
	user := &bipedaler{}
	err = datastore.Get(c, key, user)
	if err != nil {
		exit(w, "User not found")
		return
	}

	t, err := template.New("bipedaler.html").ParseFiles("views/bipedaler.html")
	if err != nil {
		exit(w, err.Error())
		return
	}

	blogPosts := []*BlogPost{}
	query := datastore.NewQuery("BlogPost").Filter("AuthorID =", userID).Project("Title", "PostID")
	iterator := query.Run(c)
	for {
		blogPost := &BlogPost{}
		_, err := iterator.Next(blogPost)
		if err == datastore.Done {
			break
		} else if err != nil {
			log.Errorf(c, "%s", err)
			continue
		}
		blogPosts = append(blogPosts, blogPost)
	}

	err = t.Execute(w, struct {
		User  *bipedaler
		Posts []*BlogPost
	}{
		user,
		blogPosts,
	})
	if err != nil {
		log.Errorf(c, "%s", err)
	}
}

func exit(w http.ResponseWriter, msg string) {
	w.Write([]byte(msg))
}
