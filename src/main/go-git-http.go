package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os/exec"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/{repo}/info/refs", index)
	//r.HandleFunc("/qwe/info/refs", outdex)
	http.Handle("/", r)
	http.ListenAndServe(":9090", nil)
}

func outdex(w http.ResponseWriter, r *http.Request) {
	test := fmt.Sprintf("%04x# service=upload-pack", 20)
	log.Println(test)
}
func index(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoName := vars["repo"]
	log.Println("repo_name: ", repoName)
	log.Println("req url: ", r.URL.Path)

	service := r.URL.Query().Get("service")
	log.Println("service: ", service)
	pFirst := fmt.Sprintf("# service=%s\n", service)
	log.Println(pFirst)

	version := r.Header.Get("Git-Protocol")
	log.Println("version", version)

	repoPath := "/Users/freddie/private/study-docs/test/1.git"
	// git upload-pack --stateless-rpc --advertise-refs /Users/freddie/private/study-docs/test/qwea
	cmdRefs := exec.Command("git", service[4:], "--stateless-rpc", "--advertise-refs", repoPath)
	refsBytes, _ := cmdRefs.Output()

	responseBody := fmt.Sprintf("%04x# service=%s\n0000%s", len(pFirst)+4, service, string(refsBytes)) // 拼接 Body
	fmt.Printf(responseBody)
	handleRefsHeader(&w, service)
	_, _ = w.Write([]byte(responseBody))
}

func handleRefsHeader(w *http.ResponseWriter, service string) {
	cType := fmt.Sprintf("application/x-%s-advertisement", service)
	(*w).Header().Add("Content-Type", cType)
	(*w).Header().Set("Expires", "Fri, 01 Jan 1980 00:00:00 GMT")
	(*w).Header().Set("Pragma", "no-cache")
	(*w).Header().Set("Cache-Control", "no-cache, max-age=0, must-revalidate")
}

//todo new feature:  will support "git clone" soon
func nothing() {

}
