package RequestListener

import (
	"fmt"
	"gladiators/src/ProxyConnector/Messages"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

type RequestListener struct {
	localAddr string
	proxyAddr string
}

func NewRequestListener(localAddr, proxyAddr string) *RequestListener {
	return &RequestListener{
		localAddr: localAddr,
		proxyAddr: proxyAddr,
	}
}
func (rl *RequestListener) Listen(internalMessageChan chan messages.ProxyClientMsg) {
	proxy := &httputil.ReverseProxy{Director: func(req *http.Request) {
		rl.handleDirectorRequest(req, internalMessageChan)
	}} // 创建代理服务器对象
	http.ListenAndServe(rl.localAddr, proxy)
}

func (rl *RequestListener) handleDirectorRequest(req *http.Request, internalMessageChan chan messages.ProxyClientMsg) {
	// Only intercept requests from the Chrome
	if strings.Contains(req.Header.Get("User-Agent"), "Chrome") {

		request, err := messages.CreateNewMessage(messages.ClientRequestType, req)

		if err != nil {
			log.Fatal(err)
		}

		req.URL.Scheme = "https"    // 修改请求协议
		req.URL.Host = rl.proxyAddr // 修改请求主机名

		fmt.Println("收到请求：", request.GetStringifyData())

		internalMessageChan <- request

	}
}
