package main

import (
	// go-micro plugins
	_ "github.com/micro/go-plugins/client/grpc"
	_ "github.com/micro/go-plugins/server/grpc"

	_ "github.com/micro/go-plugins/registry/kubernetes"
)
