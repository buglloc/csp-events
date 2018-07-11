package main

import (
	"net/http"
	"strings"

	"github.com/buglloc/csp-events/evil"
	"github.com/buglloc/csp-events/victim"
)

type HostSwitch map[string]http.Handler

func (hs HostSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.Host, "victim") {
		hs["victim"].ServeHTTP(w, r)
	} else {
		hs["evil"].ServeHTTP(w, r)
	}
}

const (
	victimBaseUri = "http://victim-csp.buglloc.com:9001"
	evilBaseUri   = "http://evil-csp.buglloc.com:9001"
)

func main() {
	hs := make(HostSwitch)
	hs["victim"] = victim.NewVictimRouter()
	hs["evil"] = evil.NewEvilRouter(victimBaseUri, evilBaseUri)

	http.ListenAndServe(":9001", hs)
}
