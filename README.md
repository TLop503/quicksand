# Browser Sandbox for Viewing Potentially Malicious Content

## Project Overview
A Browserling alternative using containers as a way to investigate websites. This will be a web-based user interface that will utilize an iframe connected to a sandboxed browser to isolate the content and enhance overall security. The use of containers will allow for easy cleanup without directly restarting the host. There will be an opt-in IP anonymization trigger that routes the IP through Tor and is compatible with Tor and Firefox.

**Project Purpose:**
- Serve as a safe and agile alternative to Browserling
- Provides a sandboxed environment where malicious sites can be analyzed without risk
- Ensure usability and consistency with standard browsing experiences

**Key Features:**
- Web-based design UI with sandboxed iframe rendering
- Containerized environment for safe isolation
- Automated cleanup and reset
- Optional Tor routing for IP anonymization
- Support for multiple browsers 

## Tool Requirements
**Go:** v1.24.7  
**SDK:** https://github.com/docker/go-sdk (specific versions can be found in go.mod)  
**Firefox Docker Container:** v25.09.1 via https://github.com/jlesage/docker-firefox

### Dependency Notes
**Go Modules:** Dependencies managed with the go.mod file  
**go mod tidy:** Pulls in libraries  
**Docker:** github.com/docker/docker

## Deliverables
* Single installer for setup
  * One script or executable installs and activates the sandbox automatically on all supported platforms.
* WebUI
  * Locally accessible live sandbox window, a Tor toggle switch to enable and disable Tor Browser anonymity features, and a clean-up button that terminates and restores the sandbox to a clean snapshot. 
* CI/CD deployment
  * New container, or environment changes pushes to the repository automatically triggering builds published to GitHub Container Registry (GHCR) or equivalent.
* VM/Sandbox guide and tutorial
  * Documentation for installing and launching the sandbox will cover all aspects of software setup and usage.


## Installation Instructions
The following instructions work on Linux. Cross-platform installation is not yet implemented :).
1) **Clone Repository**
   * On your Linux system, clone the Quicksand repository in your terminal. 

2) **Install Docker**
   * **Warning:** These instructions follow **Ubuntu Linux** installation! Refer to Docker’s official documentation to find your specific system installation instructions: https://docs.docker.com/engine/install/
   * **Uninstall Conflicting Packages**
     * `sudo apt remove $(dpkg --get-selections docker.io docker-compose docker-compose-v2 docker-doc podman-docker containerd runc | cut -f1)`
   * **Set up Docker's `apt` repository**
     * Add Docker's GPG key:
       * `sudo apt update`
       * `sudo apt install ca-certificates curl`
       * `sudo install -m 0755 -d /etc/apt/keyrings`
       * `sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc`
       * `sudo chmod a+r /etc/apt/keyrings/docker.asc`
      * Add repository to Apt sources:
        * `sudo tee /etc/apt/sources.list.d/docker.sources <<EOF
  Types: deb
  URIs: https://download.docker.com/linux/ubuntu
  Suites: $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}")
  Components: stable
  Signed-By: /etc/apt/keyrings/docker.asc
  EOF`
         * `sudo apt update`
    *  Install Docker Packages:
       * `sudo apt install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin`
    *  Docker should start automatically, but to check run:
       * `sudo systemctl status docker`
  3) **Install Go Programming Language**
     * Delete any old Go installation:
       * `sudo rm -rf /usr/local/go`
     * Check:
       * `go version`
     * Install the program:
       * `wget https://go.dev/dl/(latest version of Go)`
       * e.g. `wget https://go.dev/dl/go1.24.7.linux-amd64.tar.gz`
     *  Extract Go into /usr/local:
       * `sudo tar -C /usr/local -xzf (latest version of Go)`
     * Add Go to PATH:
       * `echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc`
     * Reload:
       * `source ~/.bashrc`
     * Confirm Installation:
       * `go version`
         * Output should give you something like this: `go version go1.24.7 linux/amd64`
   
## Usage Instructions
### 1. Verify Docker is Running 
* Before starting Quicksand, ensure the Docker engine is active:
  * `sudo systemctl status docker`
* You should see:
  * Active: active (running)
* Press **Ctrl + C** to exit this status view.

### 2. Start the Quicksand application
* Chmod the ./run_quicksand script file to make it an executable:
  * `chmod +x run_quicksand.sh`
* Run the main startup script:
  * `./run_quicksand.sh`
* If the script runs successfully, you will see an output similar to the following:
    * 🔧 Quicksand - Building and Starting Server...
🧹 Cleaning previous build...
🏗️  Building Go application...
✅ Build successful!
🚀 Starting Quicksand server on http://localhost:8080
📋 Available endpoints:
   • http://localhost:8080/ (Web Interface)
   • http://localhost:8080/api/health (Health Check)
   • http://localhost:8080/api/start (Start Container)
   • http://localhost:8080/api/stop (Stop Container)
   • http://localhost:8080/api/swap (Swap between Tor and Firefox)
   • http://localhost:8080/api/restart (Restart Container) Press Ctrl+C to stop the server
==========================================
2025/12/08 22:30:45 Containers downloaded!
2025/12/08 22:30:45 Listening on :8080

### 3. Open the Web Interface
* Navigate to:
  * `http://localhost:8080`
* You should see the Quicksand web interface.
* Once it loads, click **Start Container** to launch the sandbox environment.
  * In the terminal, you should see the terminal pull the required Docker images.
  * Once these automatic downloads stop generating on your terminal, you will see a running sandbox environment!
* To shut down Quicksand, return to the terminal where it is running and press:
  * **Ctrl + C** 
 
### Troubleshooting
* If the Firefox or Tor container does not start after clicking its button in the WebUI and the interface remains on the “Reconnecting” screen, simply refresh the page. The container may have started successfully in the background, and refreshing ensures the WebUI reconnects properly.
* Make sure you have the latest version of Go

### Quicksand Web Interface Overview
The Quicksand WebUI provides a simple control panel for managing the browser sandbox containers. Each button in the interface performs a specific action through the backend API.

🟢**Start Container**

Launches a new browser container (either Firefox or Tor, depending on your current mode).

Use this when:
* You are starting the sandbox for the first time.
* The container is not running.
* You have recently stopped or restarted a container.

🔴**Stop Container**

Stops the currently running browser container.

Use this when:
* You want to cleanly shut down the sandbox.
* You want to switch modes after a failed load.
Stopping the container will disconnect the WebUI until another container is started.

🔁**Restart Container**

Stops the current container and launches a fresh one.

Use this when:
* The environment needs to be reset.
* You want to clear temporary data inside the container.
This performs a clean reset without needing a full shutdown.

🕵️**Swap to Tor**

Stops any active Firefox container and launches a Tor Browser container instead.

Use this when:
* You want to switch to a privacy-focused, Tor-routed browsing session.
* You need to test or analyze traffic through Tor.

🦊**Swap to Firefox**

Stops any active Tor container and launches a Firefox container.

Use this when:
* You want a standard browsing environment.
* You are done using Tor and returning to normal analysis.

## Security Best Practices 
The use of Quicksand significantly reduces untrusted potential harm when visiting a malicious website, but safe and responsible operation still requires proper cybersecurity practices. This section outlines important legal and social considerations to ensure that users interact with the sandbox, safely and ethically.

### Handling Potential Malicious Content
The sandboxed browser may allow files to be downloaded during analysis. These files should always be treated as malicious:
* DO NOT open downloaded files on your host system.
* Avoid transferring files out of the container unless absolutely necessary.
### Container Isolation
Quicksand relies on containerized isolation to protect the host:
* Use the **Restart Container** button after visiting high-risk website to clear browser data.
* Keep container sessions short and reset frequently.
### Keep Tool Updated
Outdated software can introduce security vulnerabilities:
* Regularly update Docker and Go to latest version.
* Periodically rebuild containers to ensure you have the latest security patches.
### Tor-Based IP Anonymity
Tor mode provides optional IP anonymization, but anonymity is not guaranteed.
* Do not perform actions that could deanonymize your traffic.
* Avoid logging into personal accounts when using Tor.
* Use Tor only for analysis scenarios that require privacy or obscuring network traffic.
### Ethical and Legal Considerations
When handling the sandbox, be aware of your ethical responsibilities:
* Ensure your actions comply with organizational policies, laws, and responsible-use guidelines.
* Avoid visiting websites with content you are not authorized to handle.
* DO NOT use the sandbox to distribute or access harmful material.

