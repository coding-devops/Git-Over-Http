package main

import (
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"syscall"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.Printf("========================================================\n")
	log.Printf("|  smart-git-server Info")
	log.Printf("|  Server Name: local-git")
	log.Printf("|  VERSION: v0.1")
	log.Printf("========================================================\n")

	r := mux.NewRouter()
	r.HandleFunc("/{repo}/info/refs", index)
	r.HandleFunc("/{repo}/git-upload-pack", HandleGitUpload)

	r.HandleFunc("/http/test", outdex)
	http.Handle("/", r)
	http.ListenAndServe(":9090", nil)
}

func Kill(cmd *exec.Cmd) {
	if cmd.Process != nil && cmd.Process.Pid > 0 {
		if cmd.SysProcAttr != nil && cmd.SysProcAttr.Setpgid {
			syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
		}
		syscall.Kill(cmd.Process.Pid, syscall.SIGTERM)
	}
}

func TestArg(n int, d func(a string, b string)) {
}

// 完成git upload
func HandleGitUpload(res http.ResponseWriter, req *http.Request) {
	repoPath := "/Users/freddie/private/study-docs/test-client/1.git"
	cmdPack := exec.Command("git", "upload-pack", "--stateless-rpc", repoPath)
	cmdStdin, err := cmdPack.StdinPipe()
	cmdStdout, err := cmdPack.StdoutPipe()
	err = cmdPack.Start()
	if err != nil {
		_, _ = res.Write([]byte(err.Error()))
		return
	}
	// transfer data
	go func() {
		_, _ = io.Copy(cmdStdin, req.Body)
		_ = cmdStdin.Close()
	}()
	_, _ = io.Copy(res, cmdStdout)
	_ = cmdPack.Wait() // wait for std complete
	Kill(cmdPack)
}

func outdex(w http.ResponseWriter, r *http.Request) {
	//test := fmt.Sprintf("%04x# service=upload-pack", 20)
	//log.Println(test)
	buf := make([]byte, 1024) // 输入流缓存数组
	n, _ := r.Body.Read(buf)
	fmt.Println("123： ", string(buf[0:n]))
	// 或者
	//body, _ := ioutil.ReadAll(r.Body)
	//fmt.Println("456: ", string(body))
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

	repoPath := "/Users/freddie/private/study-docs/test-client/1.git"
	// git upload-pack --stateless-rpc --advertise-refs /Users/freddie/private/study-docs/test-client/1.git
	cmdRefs := exec.Command("git", service[4:], "--stateless-rpc", "--advertise-refs", repoPath)
	refsBytes, _ := cmdRefs.Output()
	log.Println("tttttt")
	fmt.Print(cmdRefs.Run())
	responseBody := fmt.Sprintf("%04x# service=%s\n0000%s", len(pFirst)+4, service, string(refsBytes)) // 拼接 Body
	log.Println("tttttt")
	log.Println("response body ----以下")
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
func statusCodeWithMessage(w *http.ResponseWriter, code int, message string) {
	(*w).WriteHeader(code)

}
