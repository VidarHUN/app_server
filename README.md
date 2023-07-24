# app_server

TODO:

- [ ] Add an HTTP server to the client, so the app_server can send notifications
  - Maybe a better idea to use WebTransport because its bidirectional behaviour
  - Or use WebSockets. However, this will be no good for our use-case
- [ ] Add exec commands to start a quicrq server or relay
- [ ] Add exec commands to start a quicrq client
- [ ] Define how to test
- [ ] Make tests for the server
- [ ] Make tests for the client

## Generate certs

```shell
openssl genpkey -algorithm RSA -out localhost.key
openssl req -new -key localhost.key -out localhost.csr -subj "/CN=localhost" -reqexts SAN -config <(cat /etc/ssl/openssl.cnf <(printf "[SAN]\nsubjectAltName=DNS:localhost"))
openssl x509 -req -in localhost.csr -signkey localhost.key -out localhost.crt
```