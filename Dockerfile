# Import the required modules
import logging
import os

# Set up logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

def build_docker_image():
    """
    Builds a Docker image for the application.

    This function uses a Node.js 16 image as the base, sets the working directory,
    copies package files, installs dependencies, copies the rest of the application,
    exposes the port the app runs on, and sets the command to run the application.

    Returns:
        None
    """

    try:
        # Use an official Node.js 16 image as the base
        logging.info("Using Node.js 16 image as the base")
        image = "node:16-alpine"

        # Set the working directory in the container
        logging.info("Setting the working directory in the container")
        workdir = "/app"

        # Copy package files
        logging.info("Copying package files")
        package_files = ["package*.json"]
        os.chdir(workdir)
        for file in package_files:
            if not os.path.exists(file):
                logging.error(f"Package file {file} not found")
                raise FileNotFoundError(f"Package file {file} not found")
            os.system(f"cp {file} .")

        # Install dependencies
        logging.info("Installing dependencies")
        os.system("npm install")

        # Copy the rest of the application
        logging.info("Copying the rest of the application")
        os.system("cp -r . .")

        # Expose the port the app runs on
        logging.info("Exposing the port the app runs on")
        port = 3000
        os.system(f"EXPOSE {port}")

        # Command to run the application
        logging.info("Setting the command to run the application")
        command = ["npm", "start"]
        os.system(f"CMD {command}")

    except FileNotFoundError as e:
        logging.error(f"Error: {e}")
    except Exception as e:
        logging.error(f"An error occurred: {e}")

# Call the function to build the Docker image
build_docker_image()