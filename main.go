package main

import "crl/stream/stream"

func main() {
	go stream.Start()

	for {}
}