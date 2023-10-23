# Local-Relay-EU-US-split

## QUICRQ

### Used commands

Run each command in a different terminal!

```bash
kubectl config use-context quicrq-usa
./deploy -a 5
./quicrq_app client <relay-ip-eu> d <relay-port-eu> post:videotest:tests/new_video1_source.bin > post.csv 2>&1
./quicrq_app client <relay-ip-us> d <relay-port-us> get:videotest:lttetest3.bin > get.csv 2>&1
```

