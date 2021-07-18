## XA
The `xa` example use http protocol for all rpc call.

`tm` is a simaple transaction manager. There has 3 services: `api` and `customer` and `merchant`.

### How to start
```
GET http://api:8080/order
```
`/order` api will call `customer` and `merchant` service.

## Task
- [ ] xa example docker-compose
