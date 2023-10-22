# Baseline

## QUICRQ

### Used commands

Run each command in a different terminal!

```bash
./quicrq_app -p 3000 -c certs/cert.pem -k certs/key.pem server
./quicrq_app client 127.0.0.1 d 3000 get:videotest:lttetest3.bin > get.csv 2>&1
./quicrq_app client 127.0.0.1 d 3000 post:videotest:tests/new_video1_source.bin > post.csv 2>&1
```

