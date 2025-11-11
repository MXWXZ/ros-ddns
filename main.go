package main

import (
	"log"
	"net/http"
	"os"
	"ros-ddns/service"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-yaml"
)

var cfg struct {
	URL      string `yaml:"url"`
	Secret   string `yaml:"secret"`
	IPHeader string `yaml:"ip_header"`
}

func ddns(c *gin.Context) {
	secret := c.Query("secret")
	if secret == "" || secret != cfg.Secret {
		c.AbortWithStatus(403)
		return
	}
	ip := c.Query("ip")
	if ip == "" {
		ip = c.ClientIP()
	}
	service.Update(ip)
	c.AbortWithStatus(http.StatusOK)
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	log.Println("Reading config file")
	content, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	if err := yaml.Unmarshal([]byte(content), &cfg); err != nil {
		log.Fatal(err)
	}

	if cfg.Secret == "" {
		log.Fatal("No secret key provided")
	}

	err = service.ParseConfig(content)
	if err != nil {
		log.Fatal(err)
	}
	if !service.CheckEnabled() {
		log.Fatal("No service enabled")
	}
	if err = service.Init(); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.SetTrustedProxies(nil)
	if cfg.IPHeader != "" {
		r.TrustedPlatform = cfg.IPHeader
	}
	r.GET(cfg.URL, ddns)
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "It works!")
	})
	r.Run()
}
