package victim

import (
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/karlseguin/ccache"
	"github.com/pkg/errors"
)

var posts = ccache.New(ccache.Configure())

func updateMessage(id string, content string) (ok bool) {
	posts.Set(id, content, 5*time.Minute)
	return true
}

func getMessage(id string) (content string, ok bool) {
	if item := posts.Get(id); item != nil && !item.Expired() && item.Value() != nil {
		content = item.Value().(string)
		ok = true
	}
	return
}

func nonce() string {
	return uuid.New().String()
}

func NewVictimRouter() http.Handler {
	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/victim/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.min.html", gin.H{"nonce": nonce()})
	})
	r.POST("/message", func(c *gin.Context) {
		content := c.PostForm("content")
		if content == "" {
			c.AbortWithError(400, errors.New("empty content"))
			return
		}

		id := uuid.New().String()
		if !updateMessage(id, content) {
			c.AbortWithError(500, errors.New("failed to update message"))
			return
		}
		c.Redirect(302, "/message/"+id)
	})
	r.GET("/message/:uuid", func(c *gin.Context) {
		content, ok := getMessage(c.Param("uuid"))
		if !ok || content == "" {
			c.AbortWithError(400, errors.New("no message"))
			return
		}

		c.HTML(200, "message.min.html", gin.H{
			"content": template.HTML(content),
			"nonce":   nonce(),
		})
	})
	r.POST("/message/:uuid", func(c *gin.Context) {
		content := c.PostForm("content")
		if content == "" {
			c.AbortWithError(400, errors.New("empty content"))
			return
		}

		id := c.Param("uuid")
		if !updateMessage(id, content) {
			c.AbortWithError(500, errors.New("failed to update message"))
			return
		}
		c.Redirect(302, "/message/"+id)
	})

	return r
}
