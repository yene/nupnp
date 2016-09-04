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

## Register to nupnp.com
```
#!/bin/sh
curl -H "Content-Type: application/json" -X POST -d "{\"name\":\"$(hostname)\",\"address\":\"$(hostname -I)\"}" https://nupnp.com/api/register
```

## Inspiration
>After about 1 minute open a web browser and point to find.z-wave.me. Below the login screen you will see the IP address of your RaZberry system. Click on the IP address link to open the configuration dialog.

* http://www.meethue.com/api/nupnp
* http://find.z-wave.me

## TODO
- [ ] Test with external IP
- [ ] improve Go mutex
- [ ] add anonymous logging of the requests
- [ ] Support for IPv6
- [ ] Improve tests
- [ ] Improve [private address](https://en.wikipedia.org/wiki/Private_network) check
- [ ] Fix the responsive resize

## Branding
Create subdomains for companies, and let their devices register with custom secret.

## restarting demon
go install && killall nupnp && nohup nupnp &

## Security
Never allow another IP address to access the data. Remove the entries after 24h. If you use a proxy prevent external access to the API server.

## Font
Font used is Days.

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
