package webserver

type WebServerStarter struct {
	WebServer WebServer
}

func NewWebServerStarter(WebServer WebServer) *WebServerStarter {
	return &WebServerStarter{
		WebServer: WebServer,
	}
}
