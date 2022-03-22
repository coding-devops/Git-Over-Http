package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os/exec"
	"syscall"
)

func main() {

	log.Printf("========================================================\n")
	log.Printf("|  smart-git-server Info")
	log.Printf("|  Server Name: local-git")
	log.Printf("|  VERSION: v0.1")
	log.Printf("========================================================\n")

	r := mux.NewRouter()
	r.HandleFunc("/{repo}/info/refs", infoRefs).Methods("GET")
	r.HandleFunc("/{repo}/git-upload-pack", HandleGitUpload).Methods("POST")
	r.HandleFunc("/{repo}/git-receive-pack", ReceivePack).Methods("POST")

	r.HandleFunc("/test/epoll", epoll)

	r.HandleFunc("/http/test", outdex)
	http.Handle("/", r)
	http.ListenAndServe(":9090", nil)
}
func init() {

}

type Commonchan struct {
	stringchan chan string
}

type ReqJson struct {
	msg       string `json:"msg"`
	PkgId     string `json:"PkgId"`
	VersionId string `json:"VersionId"`
}

func epoll(res http.ResponseWriter, req *http.Request) {
	//维护一个等待被执行队列
	a := make(chan string, 100)
	commonchan := &Commonchan{
		stringchan: a,
	}
	var test *ReqJson
	var by []byte
	if _, err := req.Body.Read(by); err != nil {
		//请求体读取失败
		log.Printf("读取body报错了")
		return
	}
	err := json.Unmarshal(by, test)
	if err != nil {
		log.Printf("反序列化报错了")
		return
	}
	//启动一个 worker 协程
	go func() {
		for {
			select {
			case commonchan.stringchan <- test.msg:
				log.Printf("等待队列写入了一条数据：", test.msg)
			default:

			}
		}
	}()
}

func Kill(cmd *exec.Cmd) {
	if cmd.Process != nil && cmd.Process.Pid > 0 {
		if cmd.SysProcAttr != nil && cmd.SysProcAttr.Setpgid {
			syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
		}
		syscall.Kill(cmd.Process.Pid, syscall.SIGTERM)
	}
}

func ReceivePack(res http.ResponseWriter, req *http.Request) {
	repoPath := "/Users/freddie/private/study-docs/test-client/1.git"
	cmdPack := exec.Command("git", "receive-pack", "--stateless-rpc", repoPath)
	cmdStdin, err := cmdPack.StdinPipe()
	cmdStdout, err := cmdPack.StdoutPipe()
	err = cmdPack.Start()
	if err != nil {
		_, _ = res.Write([]byte(err.Error()))
		return
	}
	go func() {
		_, _ = io.Copy(cmdStdin, req.Body)
		_ = cmdStdin.Close()
	}()
	_, _ = io.Copy(res, cmdStdout)
	_ = cmdPack.Wait() // wait for std complete
	Kill(cmdPack)
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
	res.Header().Set("Content-Type", fmt.Sprintf("application/x-%s-result", "git-upload-pack"))
	res.Header().Set("Connection", "Keep-Alive")
	res.Header().Set("Transfer-Encoding", "chunked")
	res.Header().Set("X-Content-Type-Options", "nosniff")

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
	//fmt.Println("456: ", string(body))
}
func infoRefs(w http.ResponseWriter, r *http.Request) {
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
func statusCodeWithMessage(w *http.ResponseWriter, code int, message string) {
	(*w).WriteHeader(code)

}
