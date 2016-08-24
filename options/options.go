package options

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

type NodeConf struct {
	Seeds                 []string `toml:"seed_servers"`
	UpdateInterval        string   `toml:"update_interval"`
	HeartbeatInterval     string   `toml:"heartbeat_interval"`
	IgnoreVersionMismatch bool     `toml:"node_ignore_version_mismatch"`
}

type Conf struct {
	Root              string   `toml:"root_dir"`
	ApiPort           string   `toml:"api_port"`
	RunNode           bool     `toml:"run_as_node"`
	NodeConfig        NodeConf `toml:"node"`
	Cores             int      `toml:"cores"`
	Seed              string   `toml:"seed"`
	Cert              string   `toml:"tls_cert"`
	Key               string   `toml:"tls_key"`
	Ssl               bool     `toml:"use_ssl"`
	AuthorizedClients []string `toml:"authorized_clients"`
}

var Config Conf

func GetOptions() {
	var configFile string

	flag.StringVar(&configFile, "config", "", "Configuration file")
	flag.IntVar(&Config.Cores, "cores", 2, "Amount of cores to pass to GOMAXPROC (experimental)")
	flag.StringVar(&Config.Root, "root", "", "Root directory to serve (required). Must be absolute path")
	flag.StringVar(&Config.ApiPort, "api-port", "8081", "Port that the API listens on")
	flag.StringVar(&Config.Seed, "seed", "", "Seed server to query")
	flag.BoolVar(&Config.RunNode, "node", false, "Run as a node")
	flag.StringVar(&Config.Cert, "tls-cert", "", "Path to TLS certificate to use")
	flag.StringVar(&Config.Key, "tls-key", "", "Path to TLS key to use")
	flag.BoolVar(&Config.Ssl, "ssl", true, "Use TLS/SSL")

	flag.StringVar(&Config.NodeConfig.HeartbeatInterval, "heartbeat-interval", "30s", "How often to send a heartbeat to the server")
	flag.StringVar(&Config.NodeConfig.UpdateInterval, "update-interval", "1m", "How often to update with the other servers")
	flag.BoolVar(&Config.NodeConfig.IgnoreVersionMismatch, "node-ignore-version-mismatch", false,
		"Ignore a mismatch in server and client versions")

	flag.Parse()

	log.Println(Config.AuthorizedClients)
	if configFile != "" {
		if _, err := toml.DecodeFile(configFile, &Config); err != nil {
			fmt.Printf("Error reading config %s: %s\n", configFile, err.Error())
			os.Exit(-1)
		}
		fmt.Printf("Configration file options in %s overriding command line options\n", configFile)
	}

	if Config.Root == "" {
		fmt.Println("Must specify root directory")
		os.Exit(-1)
	}

	if Config.RunNode == true && Config.NodeConfig.Seeds == nil && Config.Seed == "" {
		fmt.Println("Must specify seed server when running as node")
		os.Exit(-1)
	} else if Config.RunNode == true && Config.NodeConfig.Seeds == nil && Config.Seed != "" {
		Config.NodeConfig.Seeds = append(Config.NodeConfig.Seeds, Config.Seed)
	}
}
