# Browser Sandbox for Viewing Potentially Malicious Content

## Project Overview
A Browserling alternative using containers as a way to investigate websites. This will be a web-based user interface that will utilize an iframe connected to a sandboxed browser to isolate the content and enhance provided security. The use of containers will allow for easy cleanup without directly restarting the host. There will be an opt-in IP anonymization trigger that routes the IP through Tor and is compatible with Chrome and Firefox.

**Project Purpose:**
- Serve as a safe and agile alternative to Browserling
- Provide a sandboxed environment where malicious sites can be analyzed without risk
- Ensure usability and consistency with standard browsing experiences

**Key Features:**
- Web-based design UI with sandboxed iframe rendering
- Containerized environment for safe isolation
- Automated cleanup
- Optional Tor routing for IP anonymization
- Support for multiple browsers (Firefox, Chromium)

## Tool Requirements
**Go:** v1.24.7  
**SDK:** https://github.com/docker/go-sdk (specific versions can be found in go.mod)  
**Firefox Docker Container:** v25.09.1 via https://github.com/jlesage/docker-firefox

### Dependency Notes
**Go Modules:** Dependancies managed with the go.mod file  
**go mod tidy:** Pulls in libraries  
**Docker:** github.com/docker/docker

## Installation Instructions

## Usage Instructions

## Deliverables
* Single installer for setup
  * One script or executable installs and activates the sandbox automatically on all supported platforms.
* WebUI
  * Locally accessible live sandbox window, a Tor toggle switch to enable and disable Tor Browser's anonymity features, and a clean-up button that terminates and restores the sandbox to a clean snapshot. 
* CI/CD deployment
  * New container, or environment changes pushes to the repository automatically triggering builds published to GitHub Container Registry (GHCR) or equivalent.
* Progress and evolution reports
  * Team's progress report and individual reports are clear, concise, correct, and delivered on time.
* VM/Sandbox guide and tutorial
  * Documentation for installing and launching the sandbox will cover all aspects of software setup and usage.
* Class Presentation
  * 
* Final Report
  *

