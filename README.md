# 🤖 nupnp

Discovery broker for IoT devices. 🤖

![screen](screen.png)

## API
Register device with:
```
curl -H "Content-Type: application/json" -X POST -d '{"name":"Testdevice","address":"192.168.100.151"}' http://localhost:8180/api/register
```

Optional parameters:
* port

List device with:
```
http://localhost:8180/api/devices
```

## Inspiration
>After about 1 minute open a web browser and point to find.z-wave.me. Below the login screen you will see the IP address of your RaZberry system. Click on the IP address link to open the configuration dialog.

* http://www.meethue.com/api/nupnp
* http://find.z-wave.me

## Questions
* Should it filter IP addresses -> just prevent simple loopback error
* Should port be inside address or separate -> separate makes scripting easier
* Should the user be able to provide full fledged address?

## TODO
- [ ] Proxy should be a flag

## Restarting demon
`go install && killall nupnp && nohup nupnp &`

## Security
Never allow another IP address to access the data. Remove the entries after 24h. If you use a proxy prevent external access to the API server.

## Caddy Proxy configuration
```
proxy /api/register localhost:8180 {
        proxy_header X-Forwarded-Proto {scheme}
        proxy_header X-Forwarded-For {host}
        proxy_header X-Real-IP {remote}
        proxy_header Host {host}
}
```

## License
[MIT](https://tldrlegal.com/license/mit-license)
