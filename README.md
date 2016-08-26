# nupnp

Discovery broker for IoT devices. 🤖

![screen](screen.png)

## API
Register device with:
```
curl -H "Content-Type: application/json" -X POST -d '{"id":"41945125","name":"Testdevice","address":"192.168.100.151"}' http://localhost:8180/api/register
```

List device with:
```
http://localhost:8180/api/devices
```

## Register to nupnp.com once a day
Put this script into `/etc/cron.daily/nupnp`.
```
#!/bin/sh
curl -H "Content-Type: application/json" -X POST -d "{\"name\":\"$(hostname)\",\"address\":\"$(hostname -I)\"}" https://nupnp.com/api/register
```
`chmod +x /etc/cron.daily/nupnp`
Test with `run-parts /etc/cron.daily`

## Inspiration
>After about 1 minute open a web browser and point to find.z-wave.me. Below the login screen you will see the IP address of your RaZberry system. Click on the IP address link to open the configuration dialog.

* http://www.meethue.com/api/nupnp
* http://find.z-wave.me

## TODO
- [ ] Check and validate query parameters.
- [ ] Add a secret to limit access.
- [ ] Are we going to support devices that register with local host?
- [ ] rate limit requests
- [ ] Create access keys for email address, which are not rate limited
- [ ] let user create custom namespaces (paired with access keys)
- [ ] Do I need to use an in memory database?
- [ ] Tests
- [ ] sanitize input
- [ ] sort devices by date added
- [ ] return success message "device added, visit https://nupnp.com"
- [ ] if id is missing, generate uuid
- [ ] Add support for port parameter
- [ ] expose date, sort by date, convert with time ago
- [ ] make a copy paste install script, LUL
- [ ] fix NUPNP logo

## restarting demon
killall nupnp && nohup nupnp &

## Security
Never allow another IP address to access the data. Remove the entries after 24h.

## Notes
Users should not use this service directly, they should not bookmark it. But they will...
Users should only use it to discover their device and then bookmark it.

The device should not publish his IP address all the time, best would be only after a startup or after user requests it.

Font used is Days.

## License
[MIT](https://tldrlegal.com/license/mit-license)
