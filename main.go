package main

import (
	"github.com/abeir/desktop-app/core/sys"
	"github.com/abeir/desktop-app/restful"
)

func openBrowser(serv *restful.Server){
	_ = sys.OpenBrowser("http://localhost:" + serv.Port)
}

func main() {
	restful.NewServer().ServerStartedListener(openBrowser).Start()
}
