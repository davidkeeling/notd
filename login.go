package notd

import (
	"net/http"
	"text/template"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	data := struct {
		LoginURL  string
		LogoutURL string
	}{}
	t, err := template.New("login.html").ParseFiles("views/login.html")
	if err != nil {
		exit(w, err.Error())
		return
	}

	usr := user.Current(c)
	if usr != nil {
		logoutURL, err := user.LogoutURL(c, "/")
		if err != nil {
			log.Errorf(c, "%s", err)
			http.Redirect(w, r, "/", http.StatusInternalServerError)
			return
		}
		data.LogoutURL = logoutURL
		err = t.Execute(w, data)
		if err != nil {
			log.Errorf(c, "%s", err)
		}
		return
	}

	loginURL, err := user.LoginURL(c, "/")
	if err != nil {
		log.Errorf(c, "%s", err)
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}
	data.LoginURL = loginURL
	err = t.Execute(w, data)
	if err != nil {
		log.Errorf(c, "%s", err)
	}
}

func userInfo(c context.Context) *bipedaler {
	usr := user.Current(c)
	if usr == nil {
		return nil
	}
	bper := &bipedaler{}
	q := datastore.NewQuery("bipedaler").Filter("Email =", usr.Email).Limit(1)
	t := q.Run(c)
	if key, err := t.Next(bper); err == nil {
		bper.UserID = key.Encode()
		return bper
	}
	bper.Email = usr.Email
	bper.Name = usr.String()
	key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "bipedaler", nil), bper)
	if err != nil {
		log.Errorf(c, "%s", err)
		return nil
	}
	bper.UserID = key.Encode()
	return bper
}
