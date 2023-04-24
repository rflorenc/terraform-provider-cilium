package main

import (
	"context"
	"flag"
	"log"
	"terraform-provider-cilium/cilium"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

const providerName = "hashicorp.com/edu/cilium"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "run the provider with support for debuggers")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: providerName,
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), cilium.New, opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
