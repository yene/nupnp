# nupnp

Discovery broker for IoT devices. 🤖

![screen](screen.png)

## API
Register device with:
```
curl -H "Content-Type: application/json" -X POST -d '{"id":"41945125","name":"Testdevice","address":"192.168.100.151"}' https://nupnp.com/register
```
`http://localhost:8080/register?id=2323&name=device&address=192.168.100.151`

List device with:
```
http://localhost:8080/devices
```

Calling the app without parameter redirects to the first internal address.

curl -H "Content-Type: application/json" -X POST -d '{"id":"41945125","name":"Testdevice","address":"192.168.100.151"}' https://nupnp.com/register


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

## Security
Never allow another IP address to access the data. Don't store the data.

## Notes
Users should not use this service directly, they should not bookmark it. But they will...
Users should only use it to discover their device and then bookmark it.

The device should not publish his IP address all the time, best would be only after a startup or after user requests it.

## License
[MIT](https://tldrlegal.com/license/mit-license)
