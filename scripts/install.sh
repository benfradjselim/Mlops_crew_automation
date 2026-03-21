#!/bin/bash

# Import required modules
source /etc/bashrc

# Set the installation directory
INSTALL_DIR="/usr/local/bin"

# Set the script name
SCRIPT_NAME="install.sh"

# Set the script path
SCRIPT_PATH=$(readlink -f "$0")

# Function to log messages
log_message() {
  local level=$1
  local message=$2
  local timestamp=$(date +"%Y-%m-%d %H:%M:%S")
  echo "$timestamp - $level - $message" >> /var/log/install.log
}

# Function to check if directory exists
check_directory() {
  local dir=$1
  if [ ! -d "$dir" ]; then
    mkdir -p "$dir"
    log_message "INFO" "Directory $dir created."
  fi
}

# Function to check if script is already installed
check_script_installed() {
  local script_dir=$1
  local script_name=$2
  if [ -f "$script_dir/$script_name" ]; then
    log_message "INFO" "Script $script_name already installed, skipping installation."
    return 0
  fi
  return 1
}

# Function to install script
install_script() {
  local script_path=$1
  local script_dir=$2
  local script_name=$3
  cp "$script_path" "$script_dir/$script_name"
  log_message "INFO" "Script $script_name installed."
  chmod +x "$script_dir/$script_name"
  log_message "INFO" "Script $script_name made executable."
}

# Main function
main() {
  # Check if installation directory exists
  check_directory "$INSTALL_DIR"

  # Check if script is already installed
  if check_script_installed "$INSTALL_DIR" "$SCRIPT_NAME"; then
    log_message "INFO" "Script installation successful."
    echo "Script installed successfully!"
    return 0
  fi

  # Install script
  install_script "$SCRIPT_PATH" "$INSTALL_DIR" "$SCRIPT_NAME"

  # Print success message
  echo "Script installed successfully!"
  return 0
}

# Call main function
main