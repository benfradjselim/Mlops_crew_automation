import logging
import os
import time
from kubernetes import client, config

# Set up logging
logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)

# Set up Kubernetes config
config.load_kube_config()
v1 = client.CoreV1Api()

def collect_logs(namespace, pod_name, container_name):
    """
    Collect logs from a Kubernetes pod.

    Args:
        namespace (str): The namespace of the pod.
        pod_name (str): The name of the pod.
        container_name (str): The name of the container.

    Returns:
        str: The collected logs.
    """
    try:
        # Get the pod logs
        response = v1.read_namespaced_pod_log(
            name=pod_name,
            namespace=namespace,
            container=container_name,
            follow=False
        )
        return response
    except client.ApiException as e:
        logger.error(f"Error collecting logs: {e}")
        return None

def collect_all_logs(namespace):
    """
    Collect logs from all pods in a namespace.

    Args:
        namespace (str): The namespace of the pods.

    Returns:
        dict: A dictionary with pod names as keys and logs as values.
    """
    logs = {}
    pods = v1.list_namespaced_pod(namespace=namespace)
    for pod in pods.items:
        for container in pod.spec.containers:
            logs[f"{pod.metadata.name}-{container.name}"] = collect_logs(namespace, pod.metadata.name, container.name)
    return logs

def main():
    namespace = os.environ.get("NAMESPACE", "default")
    pod_name = os.environ.get("POD_NAME", "my-pod")
    container_name = os.environ.get("CONTAINER_NAME", "my-container")
    logs = collect_logs(namespace, pod_name, container_name)
    if logs:
        print(logs)
    else:
        print("No logs found.")

if __name__ == "__main__":
    main()