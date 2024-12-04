package main

import "net/http"

func x(mux *http.ServeMux) {
	mux.HandleFunc("/x", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMovedPermanently)
		w.Write([]byte(`<head><meta name="go-import" content="anto.pt/x git https://github.com/Pitasi/x">
<meta http-equiv="refresh" content="0;URL='https://github.com/Pitasi/x'">
<body>Redirecting you to the <a href="https://github.com/Pitasi/x">project page</a>...`))
	})
}
