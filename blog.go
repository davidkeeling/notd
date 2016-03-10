package notd

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/russross/blackfriday"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

//BlogPost for DB
type BlogPost struct {
	Created              time.Time
	AuthorID             string
	AuthorLink           template.HTML
	Title                string
	Body                 template.HTML `datastore:",noindex"`
	RenderedBody         template.HTML
	PostID               string
	FeaturedImageID      string
	FeaturedImageBlobKey appengine.BlobKey
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	user := userInfo(c)
	isAdmin := user != nil && user.IsAdmin
	query := datastore.NewQuery("BlogPost").Order("-Created")
	blogPosts := []*BlogPost{}
	t := query.Run(c)
	for {
		blogPost := &BlogPost{}
		key, err := t.Next(blogPost)
		if err == datastore.Done {
			break
		} else if err != nil {
			log.Errorf(c, "%s", err)
			continue
		}
		blogPost.PostID = key.Encode()
		err = prepareBlogPostForTpl(c, blogPost)
		if err != nil {
			log.Errorf(c, "%s", err)
			continue
		}
		blogPosts = append(blogPosts, blogPost)
	}
	data := struct {
		BlogEntries []*BlogPost
		IsAdmin     bool
	}{
		blogPosts,
		isAdmin,
	}
	tpl, err := template.New("blog.html").ParseFiles("views/blog.html")
	if err != nil {
		log.Errorf(c, "%s", err)
	}

	err = tpl.Execute(w, data)
	if err != nil {
		log.Errorf(c, "%s", err)
	}
}

func getBlogPost(c context.Context, postID string) (*BlogPost, error) {
	key, err := datastore.DecodeKey(postID)
	if err != nil {
		return nil, err
	}
	blogPost := &BlogPost{}
	err = datastore.Get(c, key, blogPost)
	if err != nil {
		return nil, err
	}
	blogPost.PostID = postID
	err = prepareBlogPostForTpl(c, blogPost)
	return blogPost, err
}

func prepareBlogPostForTpl(c context.Context, blogPost *BlogPost) error {
	blogPost.RenderedBody = template.HTML(blackfriday.MarkdownBasic([]byte(blogPost.Body)))
	authorKey, err := datastore.DecodeKey(blogPost.AuthorID)
	if err != nil {
		return err
	}
	author := &bipedaler{}
	err = datastore.Get(c, authorKey, author)
	if err != nil {
		return err
	}
	blogPost.AuthorLink = template.HTML(fmt.Sprintf("<a href='/bipedaler/%s'>%s</a>", blogPost.AuthorID, author.Name))

	featuredImageKey, err := datastore.DecodeKey(blogPost.FeaturedImageID)
	if err != nil {
		return err
	}
	blogImage := &BlogImage{}
	err = datastore.Get(c, featuredImageKey, blogImage)
	if err != nil {
		return err
	}
	blogPost.FeaturedImageBlobKey = blogImage.BlobKey
	return nil
}
