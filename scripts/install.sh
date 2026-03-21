import subprocess

def install():
    # Check tools
    print("Checking tools...")
    subprocess.run(["which", "docker"], check=True)
    subprocess.run(["which", "helm"], check=True)

    # Build image
    print("Building image...")
    subprocess.run(["docker", "build", "-t", "my-image", "."], check=True)

    # Install Helm
    print("Installing Helm...")
    subprocess.run(["curl", "-fsSL", "https://raw.githubusercontent.com/helm/helm/master/install.sh"], check=True)
    subprocess.run(["chmod", "+x", "helm"], check=True)
    subprocess.run(["./helm", "init"], check=True)

if __name__ == "__main__":
    install()