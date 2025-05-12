package antopt

import "net/http"

func (Website) x(mux *http.ServeMux) {
	mux.HandleFunc("GET /x", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMovedPermanently)
		w.Write([]byte(`<head><meta name="go-import" content="anto.pt/x git https://github.com/Pitasi/x">
<meta http-equiv="refresh" content="0;URL='https://github.com/Pitasi/x'">
<body>Redirecting you to the <a href="https://github.com/Pitasi/x">project page</a>...`))
	})
	mux.HandleFunc("GET /x/{path...}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMovedPermanently)
		w.Write([]byte(`<head><meta name="go-import" content="anto.pt/x git https://github.com/Pitasi/x">
<meta http-equiv="refresh" content="0;URL='https://github.com/Pitasi/x'">
<body>Redirecting you to the <a href="https://github.com/Pitasi/x">project page</a>...`))
	})
}
