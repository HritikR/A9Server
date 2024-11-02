# A9 WiFi Mini-Camera HTTP Streaming Server

This project is a lightweight HTTP streaming server for the A9 WiFi mini-camera, implemented in Go. Based on [JavaScript implementation in the A9_PPPP repository](https://github.com/datenstau/A9_PPPP.git), this Go version provides a minimal, dependency-free solution to stream MJPG video directly from the camera. With this project, there's no need to set up additional libraries, making it ideal for straightforward streaming setups.

## Overview
This implementation retrieves only the MJPG video stream from the A9 camera.
- A single binary, ideal for systems where lightweight and efficient code is essential.

## Limitations
This Go implementation is limited to MJPG video streaming. If you require additional camera commands (e.g., controlling the camera or retrieving audio), consider using the [original A9_PPPP repository](https://github.com/datenstau/A9_PPPP.git) or extending this project.

## Getting Started

### Requirements
- Go 1.16 or higher
- A9 WiFi mini-camera connected to the same network

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/HritikR/A9Server
   cd A9Server
   ```

2. Build the server:
   ```bash
   go build -o a9server main.go
   ```

3. Run the server:
   ```bash
   ./a9server
   ```
   
## Usage
Once the server is running, you can connect to the MJPG stream `http://localhost:8080/stream` to view the live video feed from the A9 WiFi mini-camera.

## Contributing
Contributions are welcome! If you’d like to add features or improve functionality, please submit a pull request.

## Special Thanks
- [A9_PPPP](https://github.com/datenstau/A9_PPPP.git) for the JavaScript implementation.
- [Home Assistant Community](https://community.home-assistant.io/t/popular-a9-mini-wi-fi-camera-the-ha-challenge/230108)