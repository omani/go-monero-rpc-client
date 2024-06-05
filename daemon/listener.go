package daemon

type DaemonListenerHandler interface {
	OnNewBlockInvoke()
	AddListener(l DaemonListener)
}

type DaemonListenerHandlerSync struct {
	listeners []DaemonListener
}

func (dhl *DaemonListenerHandlerSync) OnNewBlockInvoke() {
	for _, l := range dhl.listeners {
		l.onNewBlock()
	}
}
func (dhl *DaemonListenerHandlerSync) AddListener(l DaemonListener) {
	dhl.listeners = append(dhl.listeners, l)
}

type DaemonListenerHandlerAsync struct {
	listeners []DaemonListener
}

func (dhl *DaemonListenerHandlerAsync) OnNewBlockInvoke() {
	for _, l := range dhl.listeners {
		go l.onNewBlock()
	}
}
func (dhl *DaemonListenerHandlerAsync) AddListener(l DaemonListener) {
	dhl.listeners = append(dhl.listeners, l)
}

type DaemonListener struct {
	onNewBlock func()
}
