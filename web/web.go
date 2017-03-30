package web

import (
	"html/template"
	"log"
	"net/http"

	"github.com/dcb9/boast/config"
	"github.com/dcb9/boast/web/ws"
	assetfs "github.com/elazarl/go-bindata-assetfs"
)

var wsHub = ws.NewHub()

func Serve() {
	go wsHub.Run()
	http.Handle("/static/", http.StripPrefix(
		"/static/", http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "assets/static"}),
	))

	http.HandleFunc("/replay", func(rw http.ResponseWriter, req *http.Request) {
	})

	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		tpl, err := Asset("assets/index.html")
		if err != nil {
			log.Fatal(err)
		}

		tmpl := template.Must(
			template.New("index.html").
				Delims("{%", "%}").
				Parse(string(tpl)),
		)
		data := req.Host
		err = tmpl.Execute(rw, data)
		if err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/ws", func(rw http.ResponseWriter, req *http.Request) {
		ws.Serve(wsHub, rw, req)
	})

	log.Fatal(http.ListenAndServe(config.Config.DebugAddr, nil))
}
