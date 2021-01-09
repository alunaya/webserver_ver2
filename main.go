package main

import (
	"gws-ver2/cws"
)

func main(){
	server := cws.NewServer(3000)
	server.Run()
}