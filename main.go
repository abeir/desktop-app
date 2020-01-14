package main

import (
	"github.com/abeir/desktop-app/core"
	"github.com/abeir/desktop-app/restful"
)

func openBrowser(serv *restful.Server){
	_ = core.OpenBrowser("http://localhost:" + serv.Port)
}

func main() {
	restful.NewServer().ServerStartedListener(openBrowser).Start()
}
