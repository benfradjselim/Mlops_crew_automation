import psutil
import prometheus_client
import logging
import time
import requests

# Set up logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

# Define metrics
cpu_usage = prometheus_client.Gauge('cpu_usage', 'CPU usage')
memory_usage = prometheus_client.Gauge('memory_usage', 'Memory usage')
latency = prometheus_client.Gauge('latency', 'Latency')

def collect_metrics():
    """
    Collect CPU, memory, and latency metrics.

    Raises:
        Exception: If an error occurs while collecting metrics.
    """
    try:
        # Collect CPU usage
        cpu_usage.set(psutil.cpu_percent())

        # Collect memory usage
        memory_usage.set(psutil.virtual_memory().percent)

        # Collect latency (for example, using a simple HTTP request)
        start_time = time.time()
        response = requests.get('http://example.com', timeout=5)
        response.raise_for_status()  # Raise an exception for HTTP errors
        latency.set((time.time() - start_time) * 1000)
    except Exception as e:
        logging.error(f"Error collecting metrics: {e}")

def start_metrics_server(port: int = 8000) -> None:
    """
    Start the Prometheus metrics server.

    Args:
        port (int): The port to listen on. Defaults to 8000.
    """
    try:
        prometheus_client.start_http_server(port)
    except Exception as e:
        logging.error(f"Error starting metrics server: {e}")

def main() -> None:
    """
    Main function.

    Collects metrics every 1 second and logs any errors.
    """
    start_metrics_server()

    while True:
        collect_metrics()
        time.sleep(1)

if __name__ == "__main__":
    main()