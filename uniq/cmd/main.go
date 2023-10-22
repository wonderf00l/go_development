package main

import (
	"log"

	uniq "github.com/wonderf00l/go-park-vk/uniq/pkg"
)

func main() {
	opts, err := uniq.NewOptions()
	if err != nil {
		log.Fatalf("Parsing options for the application: %s\n", err)
	}
	cfg, err := uniq.NewConfig(opts)
	if err != nil {
		log.Printf("Assembling config for the application: %s\n", err)
		return
	}
	defer cfg.Close()
	uniq.Uniquify(cfg.InputStream, cfg.OutputStream, cfg)
}
