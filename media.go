package notd

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/blobstore"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

//BlogImage is the record in the datastore of a blob
type BlogImage struct {
	BlobKey     appengine.BlobKey
	BlogImageID string
}

func mediaHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	user := userInfo(c)
	if user == nil || !user.IsAdmin {
		exit(w, "Unauthorized")
		return
	}

	blogImages, err := getBlogImages(c)
	if err != nil {
		exit(w, err.Error())
		return
	}

	t, err := template.New("media.html").ParseFiles("views/media.html")
	if err != nil {
		exit(w, err.Error())
		return
	}

	uploadURL, err := blobstore.UploadURL(c, "/media/upload", nil)
	if err != nil {
		exit(w, err.Error())
		return
	}

	err = t.Execute(w, struct {
		UploadURL string
		IsAdmin   bool
		Images    []*BlogImage
	}{
		uploadURL.String(),
		true,
		blogImages,
	})
	if err != nil {
		log.Errorf(c, "%s", err)
	}
}

func getBlogImages(c context.Context) ([]*BlogImage, error) {
	query := datastore.NewQuery("BlogImage")
	imgs := []*BlogImage{}
	iterator := query.Run(c)
	for {
		img := &BlogImage{}
		_, err := iterator.Next(img)
		if err == datastore.Done {
			break
		} else if err != nil {
			return nil, err
		}
		imgs = append(imgs, img)
	}
	return imgs, nil
}

func uploadMediaHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	user := userInfo(c)
	if user == nil || !user.IsAdmin {
		exit(w, "Unauthorized")
		return
	}
	blobCollections, _, err := blobstore.ParseUpload(r)
	if err != nil {
		exit(w, err.Error())
		return
	}
	blogImages := []*BlogImage{}
	blogImageKeys := []*datastore.Key{}
	for key, blobs := range blobCollections {
		log.Infof(c, "key: %v; blobs: %+v", key, blobs)
		for _, blob := range blobs {
			blogImageKey := datastore.NewKey(c, "BlogImage", string(blob.BlobKey), 0, nil)
			blogImages = append(blogImages, &BlogImage{BlobKey: blob.BlobKey, BlogImageID: blogImageKey.Encode()})
			blogImageKeys = append(blogImageKeys, blogImageKey)
		}
	}
	_, err = datastore.PutMulti(c, blogImageKeys, blogImages)
	if err != nil {
		log.Errorf(c, "%s", err)
	}
	http.Redirect(w, r, "/media", http.StatusFound)
}

func deleteMediaHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	user := userInfo(c)
	if user == nil || !user.IsAdmin {
		exit(w, "Unauthorized")
		return
	}
	params := mux.Vars(r)
	blogImageID := params["mediaID"]
	if blogImageID == "" {
		exit(w, "Not found")
		return
	}
	key, err := datastore.DecodeKey(blogImageID)
	if err != nil {
		exit(w, err.Error())
		return
	}
	blogImage := &BlogImage{}
	err = datastore.Get(c, key, blogImage)
	if err != nil {
		exit(w, err.Error())
		return
	}
	err = blobstore.Delete(c, blogImage.BlobKey)
	if err != nil {
		exit(w, err.Error())
		return
	}
	err = datastore.Delete(c, key)
	if err != nil {
		log.Errorf(c, "%s", err)
	}
	http.Redirect(w, r, "/media", http.StatusFound)
}

func serveBlogImageHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	mediaID := params["mediaID"]
	if mediaID == "" {
		exit(w, "no media")
		return
	}
	blobstore.Send(w, appengine.BlobKey(mediaID))
}
