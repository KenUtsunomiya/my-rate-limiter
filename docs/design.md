# Rate Limiter

## Requirements

- restricts an excessive amount of requests appropriately
- low latency
- the least memory usage
- distributed rate limiter, i.e., multiple servers or processes can share it
- returns an appropriate exception to the user when requests are restricted
- high fault tolerance, that is, the rate limiter failure never affects the entire system

## High level Architecture

![](./images/architecture.drawio.svg)
