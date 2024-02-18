package config

import (
	"fahi/pkg/util"
	"net"
	"os"
	"time"
)

type Config struct {
	Version               string `json:"version"`
	Host                  string `json:"host"`
	Port                  int    `json:"port"`
	AutoSSL               bool   `json:"auto_ssl"`
	Password              string `json:"password"`
	Lang                  string `json:"lang"`
	LogLevel              string `json:"log_level"`
	WgPrivateKey          string `json:"wg_private_key"`
	WgDevice              string `json:"wg_device"`
	WgPort                int    `json:"wg_port"`
	WgMTU                 int    `json:"wg_mtu"`
	WgPersistentKeepalive int    `json:"wg_persistent_keepalive"`
	WgAddress             string `json:"wg_address"`
	WgDNS                 string `json:"wg_dns"`
	WgAllowedIPs          string `json:"wg_allowed_ips"`
}

func LoadOrCreate() (*Config, error) {
	var cfg Config

	cfgPath := util.RootDir + "config.json"
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		wgDevice := "eth0"
		r, err := util.NewRouter()
		if err == nil {
			iface, _, _, err := r.Route(net.IPv4(0, 0, 0, 0))
			if err == nil {
				wgDevice = iface.Name
			}
		}

		host := ""
		ip, err := util.GetExternalIP(7 * time.Second)
		if err == nil {
			host = ip.String()
		}

		cfg = Config{
			Version:               "1",
			Host:                  host,
			Port:                  50121,
			Lang:                  "en",
			LogLevel:              "error",
			Password:              "firefly",
			WgPrivateKey:          util.GeneratePrivateKey(),
			WgDevice:              wgDevice,
			WgPort:                50120,
			WgMTU:                 1280,
			WgPersistentKeepalive: 25,
			WgAddress:             "198.18.0.1/16",
			WgDNS:                 "1.1.1.1",
			WgAllowedIPs:          "0.0.0.0/0, ::/0",
		}

		err = Save(&cfg)
		if err != nil {
			return nil, err
		}

		return &cfg, nil
	}

	err = util.Json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func Save(cfg *Config) error {
	data, err := util.Json.MarshalIndent(cfg, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(util.RootDir+"config.json", data, 0644)
}
