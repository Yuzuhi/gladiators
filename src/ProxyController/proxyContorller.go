package ProxyController

import (
	"gladiators/src/ProxyConnector"
	"gladiators/src/ProxyConnector/Messages"
	"gladiators/src/RequestListener"
	"sync"
)

type ProxyManager struct {
	proxyAddr      string
	connectionType string
	localAddr      string
}

func NewProxyManager(connectionType, proxyHost, proxyPort, localHost, localPort string) *ProxyManager {
	proxyAddr := proxyHost + ":" + proxyPort
	localAddr := localHost + ":" + localPort
	return &ProxyManager{
		connectionType: connectionType,
		proxyAddr:      proxyAddr,
		localAddr:      localAddr,
	}
}

func (pm *ProxyManager) Listen() error {
	var wg sync.WaitGroup

	pc := ProxyConnector.NewProxyConnector(pm.connectionType, pm.proxyAddr)

	internalMessageChan := make(chan messages.ProxyClientMsg)

	wg.Add(2)

	go pc.HandleProxyConnection(internalMessageChan)

	rl := RequestListener.NewRequestListener(pm.localAddr, pm.proxyAddr)

	go rl.Listen(internalMessageChan)

	wg.Wait()

	return nil

}
