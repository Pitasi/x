package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	addr    = flag.String("addr", ":8080", "Address to listen to (default :8080)")
	scripts = flag.String("scripts", "./scripts", "Path where to find the executable scripts")
	secret  = flag.String("secret", "", "Webhook secret")
)

func main() {
	flag.Parse()

	http.HandleFunc("POST /github", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("incoming github webhook")

		if r.Header.Get("content-type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		payload, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("reading body", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		r.Body.Close()

		sig := r.Header.Get("X-Hub-Signature-256")
		if !validateSignature([]byte(*secret), payload, sig) {
			slog.Warn("invalid signature")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		evtType := r.Header.Get("X-GitHub-Event")
		if evtType != "push" {
			slog.Warn("ignoring webhook event", "type", evtType)
			w.WriteHeader(http.StatusOK)
			return
		}

		var pushPayload GithubPush
		if err := json.Unmarshal(payload, &pushPayload); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

		go func() {
			repo := pushPayload.Repository.FullName
			scriptPath := getScriptPath(*scripts, repo)
			slog.Info("push", "repository", repo, "script", scriptPath)

			_, err := os.Stat(scriptPath)
			if errors.Is(err, os.ErrNotExist) {
				slog.Info("no script for repo, skipping", "script", scriptPath)
				return
			}
			if err != nil {
				slog.Error("os stat", "err", err)
				return
			}

			cmd := exec.Command(scriptPath)
			out, err := cmd.CombinedOutput()
			if err != nil {
				slog.Error("error executing script", "script", scriptPath, "err", err, "out", string(out))
			} else {
				slog.Info("done", "out", string(out))
			}
		}()
	})

	slog.Info("listening", "addr", *addr)
	panic(http.ListenAndServe(*addr, nil))
}

func getScriptPath(scriptsPath, repo string) string {
	return filepath.Join(scriptsPath, strings.ToLower(strings.ReplaceAll(repo, "/", "_")))
}

func validateSignature(secret, payload []byte, expected string) bool {
	sig := hmac.New(sha256.New, secret)
	sig.Write(payload)
	digest := sig.Sum(nil)
	return "sha256="+hex.EncodeToString(digest) == expected
}

type GithubPush struct {
	Repository GithubRepository `json:"repository"`
}

type GithubRepository struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
}
