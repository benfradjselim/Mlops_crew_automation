# Design Patterns Utilisés

## 1. Factory Pattern - Création des Widgets

type WidgetFactory interface {
    CreateWidget(config WidgetConfig) (Widget, error)
}

type TimeseriesWidgetFactory struct{}
func (f *TimeseriesWidgetFactory) CreateWidget(config WidgetConfig) (Widget, error) {
    return &TimeseriesWidget{Metric: config.Metric, Aggregation: config.Aggregation}, nil
}

type KpiWidgetFactory struct{}
func (f *KpiWidgetFactory) CreateWidget(config WidgetConfig) (Widget, error) {
    return &KpiWidget{Metric: config.Metric, Threshold: config.Threshold}, nil
}

## 2. Strategy Pattern - Algorithmes de Prédiction

type PredictionStrategy interface {
    Predict(values []float64, horizon int) ([]float64, error)
}

type LinearRegressionStrategy struct{}
func (s *LinearRegressionStrategy) Predict(values []float64, horizon int) ([]float64, error) {
    // ILR implementation
    return result, nil
}

type ARIMAStrategy struct{}
func (s *ARIMAStrategy) Predict(values []float64, horizon int) ([]float64, error) {
    // ARIMA implementation
    return result, nil
}

type ExponentialSmoothingStrategy struct{}
func (s *ExponentialSmoothingStrategy) Predict(values []float64, horizon int) ([]float64, error) {
    // Holt-Winters
    return result, nil
}

## 3. Observer Pattern - Streaming Métriques

type Subscriber interface {
    OnMetric(metric Metric)
}

type MetricStream struct {
    subscribers []Subscriber
}

func (s *MetricStream) Subscribe(sub Subscriber) {
    s.subscribers = append(s.subscribers, sub)
}

func (s *MetricStream) Publish(metric Metric) {
    for _, sub := range s.subscribers {
        sub.OnMetric(metric)
    }
}

## 4. Repository Pattern - Accès Données

type MetricRepository interface {
    Save(metric Metric) error
    FindByTimeRange(metricName string, from, to time.Time) ([]Metric, error)
    Aggregate(metricName string, aggregation string, from, to time.Time) (float64, error)
}

type BadgerMetricRepository struct {
    db *badger.DB
}

func (r *BadgerMetricRepository) Save(metric Metric) error {
    // Badger implementation
    return nil
}

## 5. Builder Pattern - Construction Dashboard

type DashboardBuilder struct {
    dashboard *Dashboard
}

func NewDashboardBuilder(name string) *DashboardBuilder {
    return &DashboardBuilder{
        dashboard: &Dashboard{Name: name},
    }
}

func (b *DashboardBuilder) AddWidget(widget Widget) *DashboardBuilder {
    b.dashboard.Widgets = append(b.dashboard.Widgets, widget)
    return b
}

func (b *DashboardBuilder) SetRefresh(seconds int) *DashboardBuilder {
    b.dashboard.Refresh = seconds
    return b
}

func (b *DashboardBuilder) Build() *Dashboard {
    return b.dashboard
}

## 6. Singleton Pattern - Configuration

type Config struct {
    Server   ServerConfig   `yaml:"server"`
    Storage  StorageConfig  `yaml:"storage"`
    Metrics  MetricsConfig  `yaml:"metrics"`
}

var (
    instance *Config
    once     sync.Once
)

func GetConfig() *Config {
    once.Do(func() {
        instance = loadConfig()
    })
    return instance
}

## 7. Circuit Breaker - Sources Externes

type CircuitBreaker struct {
    failures    int
    threshold   int
    state       string
    timeout     time.Duration
    lastFailure time.Time
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    if cb.state == "open" {
        if time.Since(cb.lastFailure) > cb.timeout {
            cb.state = "half-open"
        } else {
            return ErrCircuitOpen
        }
    }
    
    err := fn()
    if err != nil {
        cb.failures++
        if cb.failures >= cb.threshold {
            cb.state = "open"
            cb.lastFailure = time.Now()
        }
        return err
    }
    
    cb.failures = 0
    cb.state = "closed"
    return nil
}

## 8. Worker Pool - Collecte Concurrente

type WorkerPool struct {
    workers int
    jobs    chan Job
    results chan Result
}

func NewWorkerPool(workers int) *WorkerPool {
    return &WorkerPool{
        workers: workers,
        jobs:    make(chan Job, 100),
        results: make(chan Result, 100),
    }
}

func (p *WorkerPool) Start() {
    for i := 0; i < p.workers; i++ {
        go p.worker()
    }
}

func (p *WorkerPool) worker() {
    for job := range p.jobs {
        result := job.Execute()
        p.results <- result
    }
}
