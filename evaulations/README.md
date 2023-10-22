# Evaluation setups and results

## Setups

- QUICRQ
- LiveKit

1. (baseline) Baseline (localhost)
2. (local-relay-eu) Clients locally, relay in EU
3. (local-relay-us) Clients locally, relay in USA
4. (local-relay-eu-us-split) Clients locally, one relay in EU and one relay in USA, one client connects to EU and another to USA
5. (eu-usa-relay-split) One client in EU, one in USA, one relay in EU, one relay in USA, USA client connect to USA relay, EU client connect to EU client
6. (scale-relay) Scaling one relay. Restrict the CPU and generate traffic to scale the deployment.