# Go DNS Server with Redis Caching

This project is a DNS forwarding server written in Go from scratch, designed to provide efficient DNS query resolution. It utilizes Redis as a caching layer to improve performance by reducing redundant upstream DNS queries. 

## Features

- **DNS Query Forwarding**: Forwards DNS queries to upstream DNS servers (e.g., Google DNS, Cloudflare DNS).
- **Redis Cache Integration**: Caches DNS query results in Redis to minimize repeated requests and reduce latency.
- **Compression Support**: Efficiently handles DNS message compression as per the DNS protocol.
- **Configurable Upstream Servers**: Allows customization of upstream DNS servers.
- **Minimal Resource Consumption**: Written in Go for optimal performance and efficiency.

## Prerequisites

- **Go (>=1.19)**: Make sure you have Go installed. You can download it from [here](https://golang.org/dl/).
- **Redis**: A Redis instance running on `localhost:6379` or configured as per your needs.
- **Docker (optional)**: For containerized deployment.

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/safiachraf/DNServer-with-Redis.git

cd DNServer-with-Redis
```

### 2. Install project dependencies

```bash
go mod tidy
```

### 3. BUILD The Dockerfile

<image_name>: The name you want to give the Docker image

```bash
docker build -t <image_name> .
```

### 3. RUN The Docker image

you are exposing port 53 (DNS) to only udp traffic and port 6379 (Redis) to your host machine.

```bash
docker run -p 53:53/udp -p 6379:6379 <image_name>
```

## Usage 
Once the DNS server is running, you can test it using the dig command, which allows you to query DNS servers. Here's an example:
![example command](https://i.ibb.co/6Hm97Q5/Screenshot-2024-09-05-191901.png)

