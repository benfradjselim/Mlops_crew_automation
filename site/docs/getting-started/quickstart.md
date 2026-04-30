# Quickstart

Get Ruptura running and see the Rupture Index in under 5 minutes.

## Step 1 — Start Ruptura

```bash
docker run -d \
  --name ruptura \
  -p 8080:8080 \
  -p 4317:4317 \
  -v ruptura-data:/var/lib/ruptura \
  -e RUPTURA_API_KEY=dev-secret-change-in-prod \
  ruptura:6.2.1
```

## Step 2 — Verify health

```bash
curl http://localhost:8080/api/v2/health
```

Expected response:

```json
{"status":"ok","rupture_detection":"active","uptime_seconds":3}
```

## Step 3 — Authenticate

Pass the `RUPTURA_API_KEY` value as a Bearer token on all subsequent requests:

```bash
export API_KEY=dev-secret-change-in-prod
```

## Step 4 — Send metrics (Prometheus remote_write)

```bash
# Push a sample metric payload
curl -s -X POST http://localhost:8080/api/v2/write \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/x-protobuf" \
  --data-binary @sample.prw
```

Or configure your Prometheus to remote_write to Ruptura:

```yaml
# prometheus.yml
remote_write:
  - url: http://ruptura:8080/api/v2/write
    authorization:
      credentials: <your-api-key>
```

## Step 5 — Query the Rupture Index

```bash
# By host
curl -s -H "Authorization: Bearer $API_KEY" \
  http://localhost:8080/api/v2/rupture/web-01 | python3 -m json.tool

# By Kubernetes workload (namespace/name)
curl -s -H "Authorization: Bearer $API_KEY" \
  http://localhost:8080/api/v2/rupture/default/my-deployment | python3 -m json.tool
```

## Step 6 — Query composite KPI signals

```bash
# health_score (composite 0–100)
curl -s -H "Authorization: Bearer $API_KEY" \
  "http://localhost:8080/api/v2/kpi/health_score/web-01"

# All 10 signals
for sig in stress fatigue mood pressure humidity contagion resilience entropy velocity health_score; do
  echo -n "$sig: "
  curl -s -H "Authorization: Bearer $API_KEY" \
    "http://localhost:8080/api/v2/kpi/$sig/web-01" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('value','?'))"
done
```

## Step 7 — View anomaly events

```bash
# All hosts, last 15 min
curl -s -H "Authorization: Bearer $API_KEY" \
  http://localhost:8080/api/v2/anomalies | python3 -m json.tool

# Single host
curl -s -H "Authorization: Bearer $API_KEY" \
  "http://localhost:8080/api/v2/anomalies/web-01"

# Custom time window
curl -s -H "Authorization: Bearer $API_KEY" \
  "http://localhost:8080/api/v2/anomalies?since=2026-04-30T00:00:00Z"
```

## Step 8 — Explain a prediction

```bash
curl -s -H "Authorization: Bearer $API_KEY" \
  http://localhost:8080/api/v2/explain/<rupture_id> | python3 -m json.tool
```

This returns the formula, contributing signals, and recommended action.

---

## Next steps

- [Configuration reference →](configuration.md)
- [API reference →](../api/reference.md)
- [Adaptive ensemble →](../concepts/surge-profiles.md)
- [Action engine →](../concepts/action-engine.md)
