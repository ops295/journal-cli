# Contributing to Journal CLI

Thank you for your interest in contributing to Journal CLI! We welcome contributions from the community and appreciate your help in making this project better.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Contribution Workflow](#contribution-workflow)
- [Coding Guidelines](#coding-guidelines)
- [Testing](#testing)
- [Submitting Pull Requests](#submitting-pull-requests)
- [Reporting Issues](#reporting-issues)
- [Community](#community)

## Code of Conduct

This project adheres to a simple [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please be respectful, collaborative, and kind in all interactions.

## Getting Started

Before you begin:

- Make sure you have Go 1.21 or higher installed
- Familiarize yourself with the project by reading the [README.md](README.md)
- Check out the [architecture.md](architecture.md) to understand the project structure
- Browse existing [issues](../../issues) and [pull requests](../../pulls) to see what's already being worked on

## Development Setup

1. **Fork the repository** on GitHub

2. **Clone your fork**:

   ```bash
   git clone https://github.com/YOUR_USERNAME/journal-cli.git
   cd journal-cli
   ```

3. **Add the upstream remote**:

   ```bash
   git remote add upstream https://github.com/ORIGINAL_OWNER/journal-cli.git
   ```

4. **Install dependencies**:

   ```bash
   make deps
   ```

5. **Build the project**:

   ```bash
   make build
   ```

6. **Run tests** to ensure everything is working:
   ```bash
   make test
   ```

## How to Contribute

There are many ways to contribute to Journal CLI:

### 🐛 Report Bugs

Found a bug? Please [open an issue](../../issues/new) with:

- A clear, descriptive title
- Steps to reproduce the problem
- Expected vs. actual behavior
- Your environment (OS, Go version)
- Any relevant logs or screenshots

### 💡 Suggest Features

Have an idea? We'd love to hear it! [Open an issue](../../issues/new) describing:

- The problem you're trying to solve
- Your proposed solution
- Any alternative solutions you've considered
- How this benefits other users

### 📝 Improve Documentation

Documentation improvements are always welcome:

- Fix typos or clarify existing docs
- Add examples or use cases
- Improve code comments
- Create tutorials or guides

### 🔧 Submit Code Changes

Ready to code? Great! See the [Contribution Workflow](#contribution-workflow) below.

## Contribution Workflow

1. **Create a new branch** for your work:

   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/your-bug-fix
   ```

2. **Make your changes**:

   - Write clean, readable code
   - Follow the [Coding Guidelines](#coding-guidelines)
   - Add tests for new functionality
   - Update documentation as needed

3. **Format your code**:

   ```bash
   make fmt
   ```

4. **Run the linter**:

   ```bash
   make lint
   ```

   > **Note**: If you don't have `golangci-lint` installed, get it from [golangci-lint.run](https://golangci-lint.run/usage/install/)

5. **Run tests**:

   ```bash
   make test
   ```

6. **Check test coverage** (optional):

   ```bash
   make coverage
   ```

7. **Commit your changes**:

   ```bash
   git add .
   git commit -m "Brief description of your changes"
   ```

   **Commit message guidelines**:

   - Use the present tense ("Add feature" not "Added feature")
   - Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
   - Limit the first line to 72 characters or less
   - Reference issues and pull requests when relevant

8. **Keep your branch up to date**:

   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

9. **Push to your fork**:

   ```bash
   git push origin feature/your-feature-name
   ```

10. **Open a Pull Request** on GitHub

## Coding Guidelines

### General Principles

- **Keep it simple**: Prefer clarity over cleverness
- **Write idiomatic Go**: Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- **DRY (Don't Repeat Yourself)**: Extract common functionality into reusable functions
- **Single Responsibility**: Each function/package should do one thing well

### Code Style

- Use `gofmt` for formatting (run `make fmt`)
- Follow standard Go naming conventions
- Write meaningful variable and function names
- Add comments for exported functions and complex logic
- Keep functions small and focused

### Project Structure

```
journal-cli/
├── cmd/           # Command-line entry points
├── internal/      # Private application code
├── templates/     # Default template files
├── Makefile       # Build automation
└── go.mod         # Go module definition
```

### Dependencies

- Minimize external dependencies
- Only add well-maintained, popular packages
- Discuss major dependency additions in an issue first

## Testing

We strive for high test coverage. When contributing:

### Writing Tests

- Add tests for all new functionality
- Update tests when modifying existing code
- Use table-driven tests where appropriate
- Test edge cases and error conditions

### Running Tests

```bash
# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Generate coverage report
make coverage
```

### Test Coverage

- Aim for at least 80% coverage for new code
- Check coverage with `make coverage`
- Don't sacrifice code quality for coverage metrics

## Submitting Pull Requests

### Before Submitting

- [ ] Code follows the project's style guidelines
- [ ] All tests pass (`make test`)
- [ ] Linter passes (`make lint`)
- [ ] Code is formatted (`make fmt`)
- [ ] Documentation is updated
- [ ] Commit messages are clear and descriptive

### PR Description

Your pull request should include:

- **What**: A clear description of what you changed
- **Why**: The motivation behind the change
- **How**: Any implementation details worth noting
- **Testing**: How you tested the changes
- **Screenshots**: If applicable (for UI changes)

### Review Process

- Maintainers will review your PR
- Address any feedback or requested changes
- Be patient and respectful during the review process
- Once approved, a maintainer will merge your PR

### Draft PRs

If you're working on something and want early feedback:

- Open a **Draft Pull Request**
- Clearly state what feedback you're looking for
- We're happy to guide you through the process!

## Reporting Issues

### Security Issues

If you discover a security vulnerability, please **DO NOT** open a public issue. Instead, email the maintainers directly (contact information in the repository).

### Bug Reports

When reporting bugs, include:

- **Environment**: OS, Go version, Journal CLI version
- **Steps to reproduce**: Clear, numbered steps
- **Expected behavior**: What should happen
- **Actual behavior**: What actually happens
- **Logs/Screenshots**: Any relevant output or images
- **Additional context**: Anything else that might help

### Feature Requests

When requesting features:

- Check if it's already been requested
- Explain the use case clearly
- Describe the desired behavior
- Consider how it fits with the project's goals

## Community

### Getting Help

- **Questions**: Open a [discussion](../../discussions) or issue
- **Chat**: Join our community chat (if available)
- **Documentation**: Check the [README.md](README.md) and [architecture.md](architecture.md)

### Recognition

We value all contributions! Contributors will be:

- Acknowledged in release notes
- Listed in the project's contributors
- Appreciated by the community 🎉

### Maintainer Responsibilities

Maintainers will:

- Review PRs in a timely manner
- Provide constructive feedback
- Help guide new contributors
- Maintain a welcoming environment

---

## Quick Reference

### Useful Commands

```bash
make help          # Show all available commands
make build         # Build the binary
make test          # Run tests
make lint          # Run linter
make fmt           # Format code
make coverage      # Generate coverage report
make clean         # Clean build artifacts
make run           # Build and run
```

### Getting Started Checklist

- [ ] Fork and clone the repository
- [ ] Set up development environment
- [ ] Run `make deps` to install dependencies
- [ ] Run `make test` to verify setup
- [ ] Create a feature branch
- [ ] Make your changes
- [ ] Run `make fmt && make lint && make test`
- [ ] Commit and push your changes
- [ ] Open a pull request

---

Thank you for contributing to Journal CLI! Your efforts help make this project better for everyone. 🙏

If you have any questions, don't hesitate to ask. We're here to help!
