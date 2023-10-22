# Local-Relay-USA

## QUICRQ

### Used commands

Run each command in a different terminal!

```bash
kubectl config use-context quicrq-usa
./deploy -a 2
./quicrq_app client <relay-ip> d <relay-port> get:videotest:lttetest3.bin > get.csv 2>&1
./quicrq_app client <relay-ip> d <relay-port> post:videotest:tests/new_video1_source.bin > post.csv 2>&1
```

