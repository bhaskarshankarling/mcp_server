# WORKSHOP

## Building a Simple MCP Server for EHQ

In this workshop, we will be building a simple MCP server for EHQ, which can be integrated into GitHub Copilot (or any LLM application) and accessed through Copilot Chat. In this session, we will focus on integrating just one API of EHQ and invoking it using Copilot Chat through the MCP server. Below is the development setup we will use for building the MCP server.

### Development Setup
- Golang (version 1.22)

#### Note
- MCP servers use different communication protocols (stdio, HTTP, WebSocket). In this session, we will focus only on building an MCP server that supports the HTTP protocol.
- We will primarily work in Agent mode, using Claude Sonet 4 (as I built the MCP server using it).

## Tasks

### 1. Understand What MCP Is - [15 min]
#### Task Description
- Spend 15 minutes understanding MCP and MCP servers using Copilot Chat and prompts.

#### Use Case
- We want straightforward information about a concept. "Zero-shot prompting" makes sense here. We can also try the "Persona" prompt and see how the LLM responds in both cases.

#### Prompt Examples
##### Zero-Shot Pattern
```
Explain the Model Context Protocol and Server in detail. Your explanation should cover:
- The core concept and relevant technical details
- A visual diagram, if applicable
- Real-world use cases—from both development and business perspectives
- Key advantages in the context of modern AI applications
```

##### Persona Pattern
```
You are a seasoned Machine Learning and AI Infrastructure Engineer, with deep expertise in Large Language Models (LLMs) and distributed systems.

Please provide a comprehensive explanation of the Model Context Protocol and Server, covering:
- The core concept, including relevant technical architecture and protocol-level details
- A conceptual or technical diagram, if possible
- Practical use cases from both development and business/product perspectives
- Advantages of using this protocol in the modern AI ecosystem, especially in scalable LLM deployment and context management
```

### 2. Install Golang 1.22 - [5 min]
#### Task Description
- Use GitHub Copilot Chat in Agent mode to set up Golang 1.22 locally.

#### Use Case
- Simple request suitable for "Zero-shot prompting."

#### Prompt Examples
##### Zero-Shot Pattern
```
Install Golang version 1.22 on my system.
```

```
Install Go (Golang) version 1.22 on my system. Identify and resolve any issues or errors during installation. Verify the Go installation by checking its availability from the command line.
```

### 3. Build and Run the Current Setup - [5 min]
#### Task Description
- Build and run the current Hello World `main.go` file under `cmd/server/`.

#### Use Case
- Suitable for "Iteration and Task Breakdown" prompt, allowing step-by-step error resolution.

#### Prompt Example
##### Iteration and Task Breakdown Pattern
```
Build `main.go` under `cmd/server` using a Makefile and output the build under `cmd/`.
```
If unsuccessful,
```
Fix the #terminalSelection error and build this project again.
```
If successful,
```
Run the build created above (`ehq_mcp_server`).
```

### 4. Build and Run a Simple MCP Server Using HTTP Protocol - [10 min]
#### Task Description
- Build a simple MCP server using HTTP protocol with one tool:
  - `get_time`: MCP tool to get the current time from the MCP server.

#### Use Case
- Ideal for "Persona pattern." Setting the perspective of an ML/AI Engineer can yield better results.

#### Prompt Example
##### Persona Pattern
```
You are an experienced Machine Learning and AI Infrastructure Engineer, specializing in scalable systems and tooling for LLM-based environments.

Your task is to build a minimal MCP (Model Context Protocol) server using HTTP protocol, exposing the following tool:
- `get_time`: Returns the current server time when invoked.

Design the server with simplicity and clarity, adhering to best practices.
Provide:
- Modifications to the existing `main.go` file under the `cmd/server` directory.
- Ensure it accepts the PORT parameter at runtime.
- Complete code and documentation for the MCP server.
- Sample interaction using `get_time` tool (via HTTP request).
```

### 5. Integrate Inside Copilot Chat - [10 min]
#### Task Description
- Integrate and test the `get_time` tool within VS Code using Copilot Chat.

### 6. Add Integration Test and Test Runner - [5 min]
#### Task Description
- Use "Zero-shot prompt" to leverage existing session context.

#### Prompt Example
##### Zero-Shot Pattern
```
Implement end-to-end integration tests for `main.go` in `cmd/server`. Cover straightforward happy path scenarios.
- Create or verify an existing test runner using a Makefile.
- Execute tests using the test runner.
- Fix failing tests and repeat until all tests pass.
```

### 7. Extend MCP to Integrate Project API - [15 min]
#### Task Description
- Add another MCP tool called `get_projects` that hits the `GET api/v2/projects?filterable=true` API and displays results in Copilot chat.
- Use root URL: `https://naveen.vikings.bangthetable.in/`
- Implement token generation logic by hitting `POST api/v2/tokens` API.

#### Prompt Example
##### Cognitive Verifier Pattern
```
Implement `GetProjects` tool calling an external API. Ask questions for clarifications (API endpoint, root URL, authentication, request/response format) before implementation. Wait for "Go ahead" to proceed.
```

### 8. Extend Projects Tool to Take 'Search' Parameter - [5 min]
#### Task Description
- Extend `GetProjects` tool to accept a `name` parameter, querying projects matching the name (`api/v2/projects?filterable=true&filters[search]=<name>`).

#### Use Case
- "Iteration and Task Breakdown."

### 9. Build and Test Integration with Copilot Chat
- Refer to earlier build and test tasks to accomplish this.

