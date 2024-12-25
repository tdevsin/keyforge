# Keyforge Integration Test

These tests depend on the application code and generate a binary during runtime. This binary will be executed, and the Keyforge service will be started. During the tests, actual calls will be made to the gRPC endpoints, which should return the expected responses. Ensure all dependencies are installed on the system before running the integration tests.

> The integration tests should run sequentially rather than in parallel. The reason is that they may start on different ports in each test, causing connectivity issues.
