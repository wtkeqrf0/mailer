Config example:
```yaml
server:
  name: EMAIL_SENDER

email:
  host: ""
  port: 465
  username: ""
  password: ""
  returnPath: ""
  name: Kain Technologies
  errorsTo: ""
  privateKeyPath: ""

rabbit:
  email:
    url: amqp://user:password@rabbit:5678
    queueName: email
  clog:
    url: amqp://admin:admin@rabbit:5672
    queue: logs

mongo:
  url: mongodb://user:password@mongo:27027
  dbName: mailer
```