# WORKSHOP

## Building Simple MCP server for EHQ

In this workshop we will be building a simple MCP server for EHQ which can be integrated into Github Copilot (can be any LLM application) and can be accessed through Copilot Chat. For this current session we will focus on integrating just one API of EHQ and invoke it using Copilot chat through MCP server. Below development setup we will be using for building MCP server.

### Dev Setup
- Golang (1.22), That is it

#### Note
- There are MCP servers with different communication protocols (stdio, HTTP, WebSocket). For this session we will just focus on building one with supporting HTTP protocol
- We will be mostly working in Agent mode and using Claude Sonnet 4 (because I built the MCP server using it)

## Tasks

### 1. Understand what MCP is - [15 min]
#### Task description
  - In this task we will spend 15 mins trying to understand what MCP and MCP server is using Copilot Chat and Prompts
#### Use case:
  - We just want some straight forward information and want to know about a concept. "Zero-shot prompting" makes sense over here. We can also try "Persona" prompt and see how LLM responds in both the cases.
#### Prompt example
#### Zero-shot pattern
```
Explain the Model Context Protocol and Server in detail. Your explanation should cover:
- The core concept and relevant technical details
- A visual diagram, if applicable
- Real-world use cases — from both development and business perspectives
- Key advantages in the context of modern AI applications
```
#### Persona pattern
```
You are a seasoned Machine Learning and AI Infrastructure Engineer, with deep expertise in Large Language Models (LLMs) and distributed systems.
- Please provide a comprehensive explanation of the Model Context Protocol and Server, covering the following aspects:
  - The core concept, including relevant technical architecture and protocol-level details
  - A conceptual or technical diagram, if possible
  - Practical use cases from both a development and business/product perspective
  - The advantages of using this protocol in the modern AI ecosystem, especially in the era of scalable LLM deployment and context management
```

### 2. Install Golang 1.22 - [5 min]
#### Task description
  - In this task we will be using Github Copilot Chat in Agent mode to set up Golang 1.22 in our local.
#### Use case:
  - Again a very simple request. Suitable candidate for "zero-shot prompting"

#### Prompt example
#### Zero-shot pattern
Example 1
```
Install golang version 1.22 in my system
```

Example 2
```
Install Go (Golang) version 1.22 on my system. Ensure that any issues or errors encountered during the installation process are identified and resolved. After installation, verify that Go is correctly installed by checking its availability from the command line and confirming that it runs as expected.
```
### 3. Build and run the current setup - [5 min]
#### Task description
- In this task we will try to build and run the current hello world main.go file under cmd/server/.
#### Use case:
- This could be good use case for "Iteration and Task Breakdown prompt" since we want to accomplish something but we don't know whether it will happen in a single stretch, there could be some errors in between. So we go step by step
#### Prompt example
#### Iteration and Task Breakdown pattern
```
Build main.go under cmd/server using Makefile and output the build under cmd/
```
Wait for it to happen, if doesn't happen successfully, let the LLM handle it if not then use the below prompt
```
Fix the #terminalSelection error and build this project again
```
If the installation goes fine, then,
```
Run the above build using the ehq_mcp_server build created above
```

### 4. Build and run a simple MCP Server using HTTP protocol - [10 min]
#### Task description
- We will have to build a very simple MCP server with HTTP protocol having two tools
  - `get_time`: A MCP tool to get the current time from the MCP server
#### Use case:
- This would be a good use case for a "Persona pattern" again. We can always prompt it straightaway but if you can set the guideline and set the tone to think from ML/AI Engineers perspective there could be a chance of getting better result
#### Prompt example
#### Persona pattern
```
You are an experienced Machine Learning and AI Infrastructure Engineer, specializing in the design of scalable systems and tooling for LLM-based environments.

Your task is to build a minimal MCP (Model Context Protocol) server using the HTTP protocol. The server should expose the following tool:
- `get_time`: A tool that returns the current server time when invoked.

Design the server with simplicity and clarity, while adhering to best practices for building reliable, modular MCP-compatible endpoints.
Provide:
- Modify existing main.go file located under the cmd/server directory to implement the MCP server.
- Make sure it accepts the PORT as a parameter while running the server
- The complete code for the MCP server
- Documentation on how it works
- A sample interaction using the `get_time` tool (e.g., via HTTP request)
```
### 5. Integrate it inside Copilot Chat - [10 min]
#### Task description
- Till here if you are able to run a simple MCP server. We can start integrating it in VS Code and start testing the `get_time` tool
### 6. Add a simple integration test and a test runner - [5 min]
#### Task description
#### Use case:
- By this time in the same session it will have context on what are the things that have been implemented. Maybe we can use "zero-shot prompt" pattern to get our work done.
#### Prompt example
#### Zero shot pattern
```
Implement end-to-end integration tests for the main.go file located in the cmd/server directory. Focus on covering the straightforward happy path scenarios for now.
- Additionally, create a test runner using a Makefile.
- If a test runner already exists, verify that it works correctly.
- Use the test runner to execute the integration tests.
- Identify and fix any failing tests.
- Repeat the process until all tests pass successfully.
```
### 7. Extend the MCP to integrate project API  - [15 min]
#### Task description
- In this task we will be adding one more MCP tool called `get_projects`. On invoking that tool from Copilot it should hit the `GET api/v2/projects?filterable=true` API and get the list of projects and display the response in the Copilot chat window.
- You can use the root URL as `https://naveen.vikings.bangthetable.in/` for now.
- You also need to implement necessary token generation logic before hitting the above project API. Which is
  - Hit the `POST api/v2/tokens` API
  - With the request body
    - ```json
      {
        "data": {
          "attributes": {
            "login": "siteadmin",
            "password": "Password123"
          }
        }
      }
      ```
    - Use the token from the tokens API as `Authorization: Bearer <token>` for the above projects API
- Response of both the APIs are in JSON-API format
  - For projects API it is
    - ```json
      {
        "data":[
          {
            "type": "projects",
            "id": "<id-1>",
            "attributes": {}
          },
          {
            "type": "projects",
            "id": "<id-2>",
            "attributes": {}
          },
          ...
        ]
      }
      ```
  - For tokens API it is
    - ```json
      {
        "data":{
          "attributes": {
            "token": "<jwt-token>"
          }
        }
      }
      ```
#### Use Case
- This is the classic use case of providing so many information in a single prompt and prompt can get really complex. So this would be a perfect example of using "Cognitive verifier pattern"
#### Prompt example
#### Cognitive verifier pattern
```
Implement one more tool called `get_projects` which should call another API and return the response of that API response. Ask more questions to get more clarifications about, API endpoint to hit, root URL, Authentication mechanism, response format, request format etc before implementing. Wait for me to say "Go ahead" to start implementing the changes till then keep asking more clariying questions
```
The response to the above prompt should start asking more questions regarding the implementation as below, For example
```
Great idea! I'd like to implement the ``get_projects`` tool that calls an external API. Before I start coding, I need some clarifications to make sure I implement it correctly:

## 🤔 **Questions for Implementation:**

### 1. **API Details:**
- **What API should we call?**
  - Do you have a specific EngagementHQ API endpoint in mind?
  - Or should I use a mock/example API for demonstration?
  - What's the base URL and endpoint path?

### 2. **Authentication:**
- **Does the API require authentication?**
  - API key? Bearer token? Basic auth?
  - Should credentials be configurable (environment variables, config file)?
  - Or should we start with a public API for simplicity?

### 3. **Request Parameters:**
- **What parameters should the tool accept?**
  - Project filters (status, owner, date range)?
  - Pagination (limit, offset)?
  - Search terms?
  - Or should it be a simple "get all projects" call?

### 4. **Response Format:**
- **How should we handle the API response?**
  - Return raw JSON?
  - Format/transform the data?
  - Extract specific fields?
  - Handle errors gracefully?

### 5. **Configuration:**
- **Where should API settings be stored?**
  - Environment variables?
  - Configuration file?
  - Command-line arguments?
  - Hardcoded for now?
```
you can answer these questions and once satisfied you can enter "Go ahead" and it should start making the changes.

### 8. Extend the projects tool to take 'search' as an input parameter - [5 min]
#### Task description
- Extend the above implementation of `get_projects` MCP tool to accept 'name' parameter from the Copilot Chat and pass it to the earlier projects API `api/v2/projects?filterable=true&filters[search]=<name>` which should return the projects with name attribute having `<name>` in it.
#### Use Case
- "Iteration and breakdown"

### 9. Build and test the integration with Copilot chat
- Take reference from the earlier build and test task to accomplish this one
