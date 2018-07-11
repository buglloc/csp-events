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

func main() {
	hs := make(HostSwitch)
	hs["victim"] = victim.NewVictimRouter()
	hs["evil"] = evil.NewEvilRouter()

	http.ListenAndServe(":9001", hs)
}
