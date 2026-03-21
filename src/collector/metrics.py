import psutil
import prometheus_client

# Define metrics
cpu_usage = prometheus_client.Gauge('cpu_usage', 'CPU usage')
memory_usage = prometheus_client.Gauge('memory_usage', 'Memory usage')
latency = prometheus_client.Gauge('latency', 'Latency')

# Function to collect metrics
def collect_metrics():
    # Collect CPU usage
    cpu_usage.set(psutil.cpu_percent())

    # Collect memory usage
    memory_usage.set(psutil.virtual_memory().percent)

    # Collect latency (for example, using a simple HTTP request)
    import requests
    start_time = time.time()
    requests.get('http://example.com')
    latency.set((time.time() - start_time) * 1000)

# Start the metrics server
prometheus_client.start_http_server(8000)

# Collect metrics every 1 second
while True:
    collect_metrics()
    time.sleep(1)