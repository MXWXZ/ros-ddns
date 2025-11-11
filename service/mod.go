package service

import (
	"bytes"
	"fmt"

	"github.com/goccy/go-yaml"
)

type Service interface {
	Name() string
	Config() interface{}
	Enabled() bool

	Init() error
	Update(ip string) error
}

var services = []Service{
	new(AliyunService),
}

func ParseConfig(s []byte) error {
	for _, srv := range services {
		path, err := yaml.PathString(fmt.Sprintf("$.%v", srv.Name()))
		if err != nil {
			return err
		}
		err = path.Read(bytes.NewReader(s), srv.Config())
		if err != nil {
			return err
		}
	}
	return nil
}

func CheckEnabled() bool {
	for _, srv := range services {
		if srv.Enabled() {
			return true
		}
	}
	return false
}

func Init() error {
	for _, srv := range services {
		if srv.Enabled() {
			err := srv.Init()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Update(ip string) error {
	for _, srv := range services {
		if srv.Enabled() {
			err := srv.Update(ip)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
