# Contributing to Probe

Thank you for your interest in contributing to Probe! This document provides guidelines and instructions for contributing.

## Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md).

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR-USERNAME/probe.git
   cd probe
   ```
3. **Add the upstream repository**:
   ```bash
   git remote add upstream https://github.com/droqsic/probe.git
   ```
4. **Create a branch** for your work:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Workflow

1. **Make your changes** in your feature branch
2. **Write or update tests** as needed
3. **Run tests** to ensure they pass:
   ```bash
   go test ./...
   ```
4. **Run benchmarks** if you're making performance-related changes:
   ```bash
   go test -bench=. ./...
   ```
5. **Format your code**:
   ```bash
   go fmt ./...
   ```
6. **Verify with go vet**:
   ```bash
   go vet ./...
   ```
7. **Commit your changes** with a clear commit message:
   ```bash
   git commit -m "Add feature: your feature description"
   ```
8. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```
9. **Create a Pull Request** from your fork to the main repository

## Pull Request Guidelines

When submitting a pull request:

1. **Include tests** for any new functionality
2. **Update documentation** as needed
3. **Follow Go coding conventions**
4. **Keep PRs focused** - submit separate PRs for unrelated changes
5. **Be responsive to feedback** - be willing to update your PR based on reviews

## Reporting Issues

When reporting issues:

1. **Use the issue template** provided
2. **Include reproduction steps** - how can we reproduce the issue?
3. **Include environment details** - Go version, OS, etc.
4. **Include logs or error messages** if applicable

## Adding New Platform Support

If you're adding support for a new platform:

1. Create a new file in the `platform` directory named `terminal_yourplatform.go`
2. Use the appropriate build tags
3. Implement the `isTerminal` and `isCygwin` functions
4. Add tests for the new platform
5. Update the README.md to include the new platform in the supported list

## License

By contributing to Probe, you agree that your contributions will be licensed under the project's [MIT License](LICENSE).
