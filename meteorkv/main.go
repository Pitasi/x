package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	cfg "github.com/cometbft/cometbft/config"
	cmtos "github.com/cometbft/cometbft/libs/os"
	"github.com/spf13/viper"
)

var homeDir string

func init() {
	flag.StringVar(&homeDir, "home", "", "Path to the CometBFT config directory (if empty, uses $HOME/.meteorkv)")
}

func main() {
	flag.Parse()
	config := setupConfig()

	cmd := flag.Arg(0)
	switch cmd {
	case "init":
		initFilesWithConfig(config)
	case "start":
		start(config)
	default:
		log.Fatalf("Usage: %s [init|start]", os.Args[0])
	}
}

func setupConfig() *cfg.Config {
	if homeDir == "" {
		homeDir = os.ExpandEnv("$HOME/.meteorkv")
	}

	// ensure directories exist
	if err := cmtos.EnsureDir(homeDir, cfg.DefaultDirPerm); err != nil {
		panic(err.Error())
	}
	if err := cmtos.EnsureDir(filepath.Join(homeDir, cfg.DefaultConfigDir), cfg.DefaultDirPerm); err != nil {
		panic(err.Error())
	}
	if err := cmtos.EnsureDir(filepath.Join(homeDir, cfg.DefaultDataDir), cfg.DefaultDirPerm); err != nil {
		panic(err.Error())
	}

	// write default config file if not exists
	defaultConfigFilePath := filepath.Join(cfg.DefaultConfigDir, cfg.DefaultConfigFileName)
	configFilePath := filepath.Join(homeDir, defaultConfigFilePath)
	if !cmtos.FileExists(configFilePath) {
		cfg.WriteConfigFile(configFilePath, cfg.DefaultConfig())
	}

	// setup viper
	config := cfg.DefaultConfig()
	config.SetRoot(homeDir)
	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Reading config: %v", err)
	}
	if err := viper.Unmarshal(config); err != nil {
		log.Fatalf("Decoding config: %v", err)
	}

	// validate config
	if err := config.ValidateBasic(); err != nil {
		log.Fatalf("Invalid configuration data: %v", err)
	}

	return config
}
