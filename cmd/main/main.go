package main

import (
	"github.com/craftizmv/rewards/config"
	"github.com/craftizmv/rewards/server"
)

func main() {

	// Initialising the config
	conf := config.GetConfig()
	// start the echo server.
	server.NewEchoServer(conf.EchoCfg).Start()

	// invoke for the consumer ...

	// init rabbitMQ

}
