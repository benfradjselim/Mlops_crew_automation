# SDKs & Client Libraries

## Go client library

The official Go client is embedded in the main module under `pkg/client`:

```go
import "github.com/benfradjselim/ruptura/pkg/client"

c := client.New("http://ruptura:8080", client.WithAPIKey("your-key"))

// Fused Rupture Index for a workload
rupture, err := c.RuptureIndex(ctx, "default", "payment-api")

// All 10 KPI signals
kpi, err := c.KPISignal(ctx, "fatigue", "default", "payment-api")

// Narrative explain
narrative, err := c.Narrative(ctx, ruptureID)
```

See `pkg/client/` in the source tree for the full interface.

## REST API (language-agnostic)

All functionality is available via the REST API v2 — use it directly from any language with a standard HTTP client:

```bash
# Health check
curl http://ruptura:8080/api/v2/health

# Rupture index for a workload
curl -H "Authorization: Bearer $API_KEY" \
  http://ruptura:8080/api/v2/rupture/default/payment-api

# KPI signal
curl -H "Authorization: Bearer $API_KEY" \
  http://ruptura:8080/api/v2/kpi/fatigue/default/payment-api

# Narrative explain
curl -H "Authorization: Bearer $API_KEY" \
  http://ruptura:8080/api/v2/explain/{rupture_id}/narrative
```

All endpoints, request/response schemas, and error codes are documented in the [API Reference →](../api/reference.md).

## Authentication

All API requests require a Bearer token matching `RUPTURA_API_KEY`:

```
Authorization: Bearer <your-api-key>
```

If `RUPTURA_API_KEY` is empty (dev/test mode), authentication is disabled and all requests are allowed.
