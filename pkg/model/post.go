package model

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/vgarvardt/rklotz/pkg/formatter"
)

const (
	postMetaSeparator = "+++"
)

var (
	// ErrorBadPostStructure is the error returned when trying to load a post with bad internal structure
	ErrorBadPostStructure = errors.New("Bad post structure: must be post meta lines, separator, post body. Separator: " + postMetaSeparator)
	// ErrorBadMetaStructure is the error returned when trying to load a post with bad meta structure
	ErrorBadMetaStructure = errors.New("Bad post meta structure, must have the following lines: post title, publishing date, post tags")
)

// Post represents post model
type Post struct {
	Path        string `storm:"id"`
	ID          string `storm:"unique"`
	Title       string
	PublishedAt time.Time `storm:"index"`
	Tags        []string
	Body        string
	Format      string
	HTML        string
}

// NewPostFromFile loads new post instance from file
func NewPostFromFile(basePath, postPath string, f formatter.Formatter) (*Post, error) {
	post := &Post{Path: postPath[len(basePath) : len(postPath)-len(filepath.Ext(postPath))]}

	fileContents, err := ioutil.ReadFile(postPath)
	if err != nil {
		return nil, err
	}

	postParts := strings.SplitN(string(fileContents), postMetaSeparator, 2)
	if len(postParts) != 2 {
		return nil, ErrorBadPostStructure
	}

	postMeta := strings.Split(strings.TrimSpace(postParts[0]), "\n")
	if len(postMeta) < 3 {
		return nil, ErrorBadMetaStructure
	}

	post.Title = postMeta[0]

	post.PublishedAt, err = time.Parse(time.RFC822Z, postMeta[1])
	if nil != err {
		return nil, err
	}

	post.Tags = strings.Split(postMeta[2], ",")
	for i, tag := range post.Tags {
		post.Tags[i] = strings.TrimSpace(tag)
	}

	post.Body = strings.TrimSpace(postParts[1])
	post.Format = strings.ToLower(filepath.Ext(postPath)[1:])

	h := sha1.New()
	h.Write([]byte(post.Path))
	post.ID = fmt.Sprintf("%x", h.Sum(nil))

	post.HTML, err = f(post.Body, post.Format)

	return post, err
}
