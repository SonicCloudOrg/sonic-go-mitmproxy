package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lqqyt2423/go-mitmproxy/addon"
	"github.com/lqqyt2423/go-mitmproxy/addon/flowmapper"
	"github.com/lqqyt2423/go-mitmproxy/addon/web"
	"github.com/lqqyt2423/go-mitmproxy/proxy"
	log "github.com/sirupsen/logrus"
)

const version = "0.1.3"

type Config struct {
	version bool

	addr         string
	webAddr      string
	ssl_insecure bool

	dump      string // dump filename
	dumpLevel int    // dump level

	mapperDir string
}

func loadConfig() *Config {
	config := new(Config)

	flag.BoolVar(&config.version, "version", false, "show version")
	flag.StringVar(&config.addr, "addr", ":9080", "proxy listen addr")
	flag.StringVar(&config.webAddr, "web_addr", ":9081", "web interface listen addr")
	flag.BoolVar(&config.ssl_insecure, "ssl_insecure", false, "not verify upstream server SSL/TLS certificates.")
	flag.StringVar(&config.dump, "dump", "", "dump filename")
	flag.IntVar(&config.dumpLevel, "dump_level", 0, "dump level: 0 - header, 1 - header + body")
	flag.StringVar(&config.mapperDir, "mapper_dir", "", "mapper files dirpath")
	flag.Parse()

	return config
}

func main() {
	config := loadConfig()

	if config.version {
		fmt.Println("go-mitmproxy: " + version)
		os.Exit(0)
	}

	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(false)
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	log.Infof("go-mitmproxy version %v\n", version)

	opts := &proxy.Options{
		Addr:              config.addr,
		StreamLargeBodies: 1024 * 1024 * 5,
		SslInsecure:       config.ssl_insecure,
	}

	p, err := proxy.NewProxy(opts)
	if err != nil {
		log.Fatal(err)
	}

	p.AddAddon(&addon.Log{})
	p.AddAddon(web.NewWebAddon(config.webAddr))

	if config.dump != "" {
		dumper := addon.NewDumper(config.dump, config.dumpLevel)
		p.AddAddon(dumper)
	}

	if config.mapperDir != "" {
		mapper := flowmapper.NewMapper(config.mapperDir)
		p.AddAddon(mapper)
	}

	log.Fatal(p.Start())
}
