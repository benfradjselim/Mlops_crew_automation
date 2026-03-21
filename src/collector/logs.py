import logging
import os
import logging.config

# Configure logging
logging.config.dictConfig({
    'version': 1,
    'formatters': {
        'default': {
            'format': '[%(asctime)s] %(levelname)s in %(module)s: %(message)s',
        }
    },
    'handlers': {
        'console': {
            'class': 'logging.StreamHandler',
            'level': 'DEBUG',
            'formatter': 'default',
        },
        'file': {
            'class': 'logging.FileHandler',
            'level': 'DEBUG',
            'filename': 'app.log',
            'formatter': 'default',
        },
    },
    'root': {
        'level': 'DEBUG',
        'handlers': ['console', 'file']
    }
})

def get_error_logs(log_path: str, log_level: str = 'ERROR') -> list:
    """
    Collect error logs from a specified log path.

    Args:
        log_path (str): Path to the log file.
        log_level (str, optional): Log level to filter by. Defaults to 'ERROR'.

    Returns:
        list: List of error logs.

    Raises:
        FileNotFoundError: If the log file does not exist.
        ValueError: If the log level is not valid.
    """
    if not isinstance(log_path, str):
        raise TypeError("log_path must be a string")
    if not isinstance(log_level, str):
        raise TypeError("log_level must be a string")

    # Validate log level
    valid_log_levels = ['DEBUG', 'INFO', 'WARNING', 'ERROR', 'CRITICAL']
    if log_level.upper() not in valid_log_levels:
        raise ValueError(f"Invalid log level: {log_level}. Must be one of {', '.join(valid_log_levels)}")

    try:
        error_logs = []
        if os.path.exists(log_path):
            with open(log_path, 'r') as f:
                for line in f:
                    if log_level.upper() in line:
                        error_logs.append(line.strip())
        return error_logs
    except FileNotFoundError:
        logging.error(f"Log file not found: {log_path}")
        return []
    except Exception as e:
        logging.error(f"An error occurred: {e}")
        return []

# Example usage:
if __name__ == "__main__":
    log_path = 'path_to_your_log_file.log'
    error_logs = get_error_logs(log_path)
    for log in error_logs:
        print(log)