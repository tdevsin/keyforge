# KeyForge

Keyforge is a high-performance, highly scalable, and fault-tolerant key-value store that can be used in:

1. **Caching**:

   - Keeping frequently used data ready for fast access, like storing a copy of your favorite website for quicker loading.

2. **Session Management**:

   - Saving user sessions, like when a website remembers you're logged in, even if you refresh the page.

3. **Configuration Storage**:

   - Storing settings for applications that need to change based on the environment, like a game adapting its controls for different players.

4. **Event Tracking**:

   - Keeping track of events in real-time systems, like logging actions in a multiplayer game.

5. **Metadata Storage**:
   - Saving extra information about files or data, like descriptions of photos in a photo-sharing app.

And in a lot more usecases.

> This project is part of the **Golang (Go) for Production Systems** course by **tdevs.in**. You can check out the project details at [https://tdevs.in/golang_in_production](https://tdevs.in/golang_in_production).

## Important Links

- [Architecture Design](https://tdevs.in/golang_in_production/contents/key-forge/architecture-design)
- [API Specification](https://tdevs.in/golang_in_production/contents/key-forge/api-specification)
- [CLI Specification](https://tdevs.in/golang_in_production/contents/key-forge/cli-specification)

Though this project is part of a learning initiative, it is highly performant and of excellent quality. It can still be used in production if it fits your use case.

## Keyforge API Benchmark Results

### Benchmark Configuration

- **CPU**: Intel(R) Core(TM) i5-7300U CPU @ 2.60GHz
- **Go Version**: `go1.23.4`
- **Operating System**: Linux (amd64)
- **Test Duration**: 10 seconds per operation
- **Concurrency**: 512 parallelism for both `GetKey` and `SetKey`

### GetKey Operation

| Metric                 | Value              |
| ---------------------- | ------------------ |
| **Iterations**         | 13,513,574         |
| **Time per Operation** | 855.8 ns           |
| **Throughput**         | ~1,168,252 req/sec |
| **Memory Usage**       | 192 B/op           |
| **Allocations**        | 3 allocs/op        |

### SetKey Operation

| Metric                 | Value           |
| ---------------------- | --------------- |
| **Iterations**         | 1,583,232       |
| **Time per Operation** | 10,467 ns       |
| **Throughput**         | ~95,508 req/sec |
| **Memory Usage**       | 391 B/op        |
| **Allocations**        | 8 allocs/op     |
