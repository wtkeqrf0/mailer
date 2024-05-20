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

A service for sending emails. The application follows the basic steps below:
1. Consume json messages from RabbitMQ
2. Try to get a sample email from MongoDB by ids in json above
3. Send email message