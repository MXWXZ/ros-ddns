package service

import (
	"fmt"
	"log"

	alidns "github.com/alibabacloud-go/alidns-20150109/v4/client"
	"github.com/alibabacloud-go/darabonba-openapi/v2/utils"
	"github.com/alibabacloud-go/tea/tea"
)

type config struct {
	Enable    bool   `yaml:"enable"`
	AK        string `yaml:"ak"`
	SK        string `yaml:"sk"`
	Endpoint  string `yaml:"endpoint"`
	Domain    string `yaml:"domain"`
	SubDomain string `yaml:"subdomain"`
}

type AliyunService struct {
	ip     string
	id     *string
	client *alidns.Client
	cfg    config
}

func (c *AliyunService) Name() string {
	return "aliyun"
}

func (c *AliyunService) Config() interface{} {
	return &c.cfg
}

func (c *AliyunService) Enabled() bool {
	return c.cfg.Enable
}

func (c *AliyunService) Init() error {
	cfg := &utils.Config{
		AccessKeyId:     tea.String(c.cfg.AK),
		AccessKeySecret: tea.String(c.cfg.SK),
		Endpoint:        tea.String(c.cfg.Endpoint),
	}
	var err error
	c.client, err = alidns.NewClient(cfg)
	if err != nil {
		return err
	}
	req := &alidns.DescribeDomainRecordsRequest{
		DomainName: &c.cfg.Domain,
		KeyWord:    &c.cfg.SubDomain,
		SearchMode: tea.String("EXACT"),
	}
	rsp, err := c.client.DescribeDomainRecords(req)
	if err != nil {
		return err
	}
	body := rsp.GetBody()
	if len(body.DomainRecords.Record) == 0 {
		return fmt.Errorf("Cannot find domain: %v.%v", c.cfg.SubDomain, c.cfg.Domain)
	}
	c.id = body.DomainRecords.Record[0].RecordId
	c.ip = *body.DomainRecords.Record[0].Value
	log.Printf("Aliyun current ip: %v", c.ip)
	return nil
}

func (c *AliyunService) Update(ip string) error {
	if ip == c.ip {
		return nil
	}
	req := &alidns.UpdateDomainRecordRequest{
		RecordId: c.id,
		Value:    tea.String(ip),
		Type:     tea.String("A"),
		RR:       tea.String(c.cfg.SubDomain),
	}
	rsp, err := c.client.UpdateDomainRecord(req)
	if err != nil {
		return err
	}
	c.id = rsp.Body.RecordId
	c.ip = ip
	log.Printf("Aliyun update ip: %v", c.ip)
	return nil
}
