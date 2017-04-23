package main

import (
	//	"bufio"
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	//	"path/filepath"
	"runtime"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

// 端口
const (
	HTTP_PORT  string = "8000"
	HTTPS_PORT string = "446"
)

// 目录
const (
	CSS_CLIENT_PATH   = "/css/"
	DART_CLIENT_PATH  = "/js/"
	IMAGE_CLIENT_PATH = "/image/"

//	CSS_SVR_PATH   = "web"
//	DART_SVR_PATH  = "web"
//	IMAGE_SVR_PATH = "web"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

}

func main() {
	//	// 先把css和脚本服务上去
	//	http.Handle(CSS_CLIENT_PATH, http.FileServer(http.Dir(CSS_SVR_PATH)))
	//	http.Handle(DART_CLIENT_PATH, http.FileServer(http.Dir(DART_SVR_PATH)))

	// 网址与处理逻辑对应起来
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/chat", ChatPage)
	http.Handle("/ws", websocket.Handler(OnWebSocket))

	// 开始服务
	err := http.ListenAndServe(":"+HTTP_PORT, nil)
	if err != nil {
		fmt.Println("服务失败 /// ", err)
	}
}

func WriteTemplateToHttpResponse(res http.ResponseWriter, t *template.Template) error {
	if t == nil || res == nil {
		return errors.New("WriteTemplateToHttpResponse: t must not be nil.")
	}
	var buf bytes.Buffer
	err := t.Execute(&buf, nil)
	if err != nil {
		return err
	}
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = res.Write(buf.Bytes())
	return err
}

func ChatPage(res http.ResponseWriter, req *http.Request) {

	t, err := template.ParseFiles("chat.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = WriteTemplateToHttpResponse(res, t)
	if err != nil {
		fmt.Println(err)
		return
	}

}

func HomePage(res http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("homepage.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = WriteTemplateToHttpResponse(res, t)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func OnWebSocket(ws *websocket.Conn) {
	defer ws.Close()

	var err error
	var str string

	buf := bytes.NewBuffer([]byte{})

	pipReader, pipWriter := io.Pipe()

	messageChan := make(chan string)

	//	//独立进程，将内容写入管道，如果没有写入内容则block住，消息内容通过channel传入到线程中
	go PutIntoPipe(pipWriter, messageChan)

	//独立进程，将内容从管道中读出并写入到文件中，如果没有可读取出来的内容则block住，消息内容通过channel传入到线程中
	go GetOutPipe(pipReader)

	for {

		if err = websocket.Message.Receive(ws, &str); err != nil {

			break

		} else {
			//messageChan <- "从客户端收到：" + str

			fmt.Println("websocket loop .... 从客户端收到：", str)

			fmt.Println("write buf .... 从客户端收到：", str)
			buf.WriteString(str)

			//			messageChan <- str
			pipWriter.Write([]byte(str + "\n"))
		}

		//		str = "hello, I'm server.\n"

		if err = websocket.Message.Send(ws, str); err != nil {
			break
		} else {
			pipWriter.Write([]byte("Server Reply:" + str + "\n"))
			fmt.Println("websocket loop .... 向客户端发送：", str)

		}
	}
}

func PutIntoPipe(write *io.PipeWriter, messageChan <-chan string) {

	//这里是一个单独的gorountin
	//里面包含的是 pipeWriter
	//持续不断的将信息写入到pipe中去

	for {
		select {
		case data := <-messageChan:
			write.Write([]byte(data))
			fmt.Println("PutIntoPipe....write conent %s", data)
			//			buffer.Reset()
			time.Sleep(2 * time.Second)
		}
	}
}

func GetOutPipe(read *io.PipeReader) {

	fmt.Printf("GetOutPipe thread is starting......")

	//这里是一个单独的gorountin,因为存在block所以必须是一个goroution
	//这里包含的是 pipeReader
	//持续不断的将pipe中的信息读取出来，并写入到文件中
	//messagebyte := make([]byte,100)

	currentDir := getCurrentDir()

	resultFile, err := os.Create(currentDir + "/test.txt")

	if err != nil {
		panic("create recored file failed!")
	}

	//	buf := bytes.NewBuffer([]byte{})
	//	var temp string

	for {
		data := make([]byte, 100)
		n, err := read.Read(data)
		if err != nil {
			fmt.Errorf("Read customer info get error!....")
		}
		fmt.Println("GetOutput....read conent", string(data[:n]))
		resultFile.WriteString(string(data[:n]))
		resultFile.Sync()
		time.Sleep(2 * time.Second)
	}
}

func getCurrentDir() string {

	currentPath := os.Args[0]
	splitArray := strings.Split(currentPath, "/")
	splitSlide := splitArray[:len(splitArray)-1]
	currentDir := strings.Join(splitSlide, "/")
	fmt.Printf("getCurrentDir() .... current Dir is: %s", currentDir)
	return currentDir
}
