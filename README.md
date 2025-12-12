# ros-ddns

DDNS framework for RouterOS.

Current support DNS: Aliyun

## RouterOS config

Use this script, add to scheduler whatever you like.

```
:local secret "mysecret"
:local pppoe "pppoe-out1"
:local ipaddr [/ip address get [/ip address find interface=$pppoe] address]
:set ipaddr [:pick $ipaddr 0 ([len $ipaddr] -3)]
:global recip
:if ($ipaddr != $recip) do={
    /tool fetch url="https://example.com/ddns?secret=$secret&ip=$ipaddr" output=none
    :log info "Set DDNS $recip to $ipaddr"
    :set recip "$ipaddr"
}
```

## Common config

- `secret`: secret key for authorization.
- `url`: url for update dns.
- `ip_header`: leave empty if you are not behind a proxy, otherwise set it to the header of real ip.

## Aliyun
Add RAM user, grant `AliyunDNSFullAccess` permission. Normally you only need to change the following configs:

- `aliyun.enable`: set to true to enable Aliyun.
- `aliyun.ak`: access key.
- `aliyun.sk`: secret key.
- `aliyun.domain`: domain name, like `example.com`.
- `aliyun.subdomain`: subdomain, like `ddns`.