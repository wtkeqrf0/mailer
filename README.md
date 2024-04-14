Config example:
```yaml
server:
  name: EMAIL_SENDER

email:
  host: ""
  port: ""
  username: ""
  password: ""
  returnPath: ""
  name: ""
  errorsTo: ""
  privateKeyPath: ""

rabbit:
  email:
    url: amqp://user:password@rabbit:5672
    queueName: ""
  clog:
    url: amqp://user:password@rabbit:5672
    queue: ""

mongo:
  url: mongodb://user:password@mongo:27017
  dbName: ""
```