# Installation

## Kubernetes

The recommended production deployment uses the bundled manifests or Helm chart.

### Using `kubectl`

```bash
git clone https://github.com/benfradjselim/ruptura.git
cd ruptura

# Build the image (or use a pre-built tag)
docker build -t ruptura:6.2.1 .

# Create namespace + deploy
kubectl apply -f deploy/

# Verify pods are running
kubectl get pods -n ruptura-system

# Port-forward to test locally
kubectl port-forward svc/ruptura 8080:8080 -n ruptura-system
curl http://localhost:8080/api/v2/health
```

### Using Helm

```bash
helm install ruptura ./helm \
  --namespace ruptura-system \
  --create-namespace \
  --set apiKey=$(openssl rand -hex 32) \
  --set storage.size=20Gi
```

To upgrade:

```bash
helm upgrade ruptura ./helm --namespace ruptura-system
```

### Using the RupturaInstance CRD (Operator)

If you have the Ruptura Operator installed, deploy a full instance declaratively:

```yaml
apiVersion: ruptura.io/v1alpha1
kind: RupturaInstance
metadata:
  name: production
  namespace: ruptura-system
spec:
  image: ruptura:6.2.1
  port: 8080
  storageSize: 20Gi
  apiKey:
    secretRef: ruptura-api-key
  replicas: 1
```

```bash
kubectl apply -f ruptura-instance.yaml
```

See [Operator →](../architecture/operator.md) for full CRD reference.

---

## Docker

```bash
docker run -d \
  --name ruptura \
  -p 8080:8080 \
  -p 4317:4317 \
  -v ruptura-data:/var/lib/ruptura \
  -e RUPTURA_API_KEY=$(openssl rand -hex 32) \
  ruptura:6.2.1
```

| Port | Protocol | Purpose |
|------|----------|---------|
| 8080 | HTTP | REST API v2, Prometheus metrics |
| 4317 | HTTP | OTLP ingest (metrics, logs, traces) |

Verify:

```bash
curl http://localhost:8080/api/v2/health
# {"status":"ok","rupture_detection":"active"}
```

---

## Binary

Download the latest release binary and run it directly:

```bash
# Linux amd64
curl -fsSL https://github.com/benfradjselim/ruptura/releases/latest/download/ruptura-linux-amd64 \
  -o /usr/local/bin/ruptura
chmod +x /usr/local/bin/ruptura

ruptura --config=/etc/ruptura/ruptura.yaml
```

Ruptura ships as a **single static binary** — no runtime dependencies, no external database.

---

## Build from Source

Requires Go 1.18+:

```bash
git clone https://github.com/benfradjselim/ruptura.git
cd ruptura/workdir
go build -o ruptura ./cmd/ruptura
./ruptura --config=configs/ruptura.yaml
```

Run tests:

```bash
go test -race -timeout=120s ./...
go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out | grep total
```
