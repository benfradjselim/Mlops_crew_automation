# API Specification

## Endpoints

### Metrics

| Endpoint | Method | Description |
|----------|--------|-------------|
| /api/v1/metrics | GET | Liste des métriques disponibles |
| /api/v1/metrics/{name} | GET | Récupérer une métrique |
| /api/v1/metrics/{name}/aggregate | GET | Métrique agrégée (avg, min, max, p95, p99) |
| /api/v1/query | POST | Requête QQL |

### KPIs

| Endpoint | Method | Description |
|----------|--------|-------------|
| /api/v1/kpis | GET | Liste des KPIs calculés |
| /api/v1/kpis/{name} | GET | Récupérer un KPI |
| /api/v1/kpis/{name}/predict | GET | Prédiction du KPI |

### Dashboards

| Endpoint | Method | Description |
|----------|--------|-------------|
| /api/v1/dashboards | GET | Liste des dashboards |
| /api/v1/dashboards | POST | Créer un dashboard |
| /api/v1/dashboards/{id} | GET | Récupérer un dashboard |
| /api/v1/dashboards/{id} | PUT | Mettre à jour |
| /api/v1/dashboards/{id} | DELETE | Supprimer |
| /api/v1/dashboards/{id}/export | GET | Exporter en JSON |
| /api/v1/dashboards/import | POST | Importer un dashboard |

### Templates

| Endpoint | Method | Description |
|----------|--------|-------------|
| /api/v1/templates | GET | Liste des templates |
| /api/v1/templates/{id} | GET | Récupérer un template |
| /api/v1/templates/{id}/apply | POST | Appliquer un template |

### Data Sources

| Endpoint | Method | Description |
|----------|--------|-------------|
| /api/v1/datasources | GET | Liste des sources |
| /api/v1/datasources | POST | Ajouter une source |
| /api/v1/datasources/{id} | GET | Récupérer une source |
| /api/v1/datasources/{id} | PUT | Mettre à jour |
| /api/v1/datasources/{id} | DELETE | Supprimer |
| /api/v1/datasources/{id}/test | POST | Tester la connexion |

### Alerts

| Endpoint | Method | Description |
|----------|--------|-------------|
| /api/v1/alerts | GET | Liste des alertes |
| /api/v1/alerts | POST | Créer une alerte |
| /api/v1/alerts/{id} | GET | Récupérer une alerte |
| /api/v1/alerts/{id} | PUT | Mettre à jour |
| /api/v1/alerts/{id} | DELETE | Supprimer |
| /api/v1/alerts/{id}/acknowledge | POST | Acquitter |
| /api/v1/alerts/{id}/silence | POST | Silencer |

### Auth

| Endpoint | Method | Description |
|----------|--------|-------------|
| /api/v1/auth/login | POST | Authentification |
| /api/v1/auth/logout | POST | Déconnexion |
| /api/v1/auth/refresh | POST | Rafraîchir token |
| /api/v1/auth/users | GET | Liste des utilisateurs |
| /api/v1/auth/users | POST | Créer un utilisateur |
| /api/v1/auth/users/{id} | GET | Récupérer un utilisateur |
| /api/v1/auth/users/{id} | PUT | Mettre à jour |
| /api/v1/auth/users/{id} | DELETE | Supprimer |

### System

| Endpoint | Method | Description |
|----------|--------|-------------|
| /api/v1/health | GET | Health check |
| /api/v1/metrics | GET | Métriques Prometheus |
| /api/v1/config | GET | Configuration |
| /api/v1/config | PUT | Mettre à jour config |
| /api/v1/reload | POST | Recharger configuration |

## Formats de réponse

### Success (200)
{
  "success": true,
  "data": {},
  "timestamp": "2026-04-13T10:00:00Z"
}

### Error (4xx/5xx)
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Description de l'erreur"
  },
  "timestamp": "2026-04-13T10:00:00Z"
}

## Authentification

Bearer token dans le header: Authorization: Bearer <token>
