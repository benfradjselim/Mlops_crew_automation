import logging
import os
import subprocess

# Set up logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

def build_docker_image(image_name: str = "node:16-alpine", workdir: str = "/app", port: int = 3000) -> None:
    """
    Builds a Docker image for the application.

    Args:
        image_name (str): The name of the Docker image to use as the base. Defaults to "node:16-alpine".
        workdir (str): The working directory in the container. Defaults to "/app".
        port (int): The port the application runs on. Defaults to 3000.

    Returns:
        None
    """

    try:
        # Use an official Node.js 16 image as the base
        logging.info("Using Node.js 16 image as the base")
        if not image_name:
            logging.error("Image name is required")
            raise ValueError("Image name is required")

        # Set the working directory in the container
        logging.info("Setting the working directory in the container")
        if not workdir:
            logging.error("Working directory is required")
            raise ValueError("Working directory is required")

        # Copy package files
        logging.info("Copying package files")
        package_files = ["package*.json"]
        os.chdir(workdir)
        for file in package_files:
            if not os.path.exists(file):
                logging.error(f"Package file {file} not found")
                raise FileNotFoundError(f"Package file {file} not found")
            subprocess.run(f"cp {file} .", shell=True, check=True)

        # Install dependencies
        logging.info("Installing dependencies")
        subprocess.run("npm install", shell=True, check=True)

        # Copy the rest of the application
        logging.info("Copying the rest of the application")
        subprocess.run("cp -r . .", shell=True, check=True)

        # Expose the port the app runs on
        logging.info("Exposing the port the app runs on")
        if port <= 0:
            logging.error("Port must be a positive integer")
            raise ValueError("Port must be a positive integer")
        subprocess.run(f"EXPOSE {port}", shell=True, check=True)

        # Command to run the application
        logging.info("Setting the command to run the application")
        command = ["npm", "start"]
        subprocess.run(f"CMD {command}", shell=True, check=True)

    except FileNotFoundError as e:
        logging.error(f"Error: {e}")
    except subprocess.CalledProcessError as e:
        logging.error(f"Error: {e}")
    except Exception as e:
        logging.error(f"An error occurred: {e}")

# Call the function to build
build_docker_image()