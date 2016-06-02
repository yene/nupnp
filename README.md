# nupnp

Discovery broker for IoT devices. 🤖

## API
Register device with
`http://localhost:3000/register?id=2323&address=192.168.100.151`

## Inspiration
* http://www.meethue.com/api/nupnp
* http://find.z-wave.me

## TODO
- [ ] Check and validate parameters.
- [ ] Add a secret to limit access.

## Notes
Users should not use this service directly, they should not bookmark it. But they will...
Users should only use it to discover their device and then bookmark it.

The device should not publish his IP address all the time, best would be only after a startup.

## License
MIT License
