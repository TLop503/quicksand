# Browser Sandbox for Viewing Potentially Malicious Content

## Project Overview
A Browserling alternative using containers as a way to investigate websites. This will be a web-based user interface that will utilize an iframe with the sandbox attribute to isolate the content and enhance provided security. The use of containers will allow for easy cleanup via Ansible without directly restarting the host. There will be an opt-in IP anonymization trigger that routes the IP through TOR and is compatible with Chrome and Firefox.

**Key features include:**
- Web-based design UI with sandboxed iframe rendering
- Containerized environment for safe isolation
- Automated cleanup via Ansible
- Optional TOR routing for IP anonymization
- Support for multiple browsers (Firefox, Chromium)

## Tool Requirements
**Go:** v1.24.7  
**Docker:** v25.0.0 via https://github.com/moby/moby  
**Firefox Docker Container:** v25.09.1 via https://github.com/jlesage/docker-firefox

### Dependency Notes
**Go Modules:** Dependancies managed with the go.mod file  
**Docker:** github.com/docker/docker is replaced with github.com/moby/moby@v25.0.0+incompatible
