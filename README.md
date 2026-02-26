# AgentAuth-Core 🛡️

**The Identity and Access Management (IAM) layer for the Autonomous Economy.**

Human OAuth flows (redirects, browser sessions, magic links) inherently fail for autonomous AI agents. As we transition to an agentic economy, we require a purely Machine-to-Machine (M2M), highly secure, low-latency authentication protocol.

`AgentAuth-Core` is a lightweight, high-performance Go server that acts as the Identity Provider (IdP) for AI agents, leveraging Ed25519 cryptography to issue short-lived JSON Web Tokens (JWTs).

## ⚡ Core Philosophy

1. **Zero-Human Intervention:** Designed strictly for M2M interactions.
2. **Cryptographic Identity:** Agents are identified by Ed25519 keypairs, not fragile API keys.
3. **Ephemeral Access:** By default, agent JWTs expire in 5 minutes, significantly reducing the blast radius of a compromised agent.

## 🏗️ Flow

```mermaid
sequenceDiagram
    participant Agent as AI Agent (Python SDK)
    participant Core as Auth Server (Go Core)
    participant API as API Gateway (TS SDK)

    Note over Agent, Core: Phase 1: Identity Minting
    Agent->>Core: POST /register (Sends Ed25519 Public Key)
    Core-->>Agent: 201 Created (Agent Registered)

    Agent->>Core: POST /token (Header: X-Agent-ID)
    Core-->>Agent: Returns Signed JWT (5-min expiry)

    Note over Agent: SDK Caches JWT locally<br/>to prevent spamming Auth Server

    Note over Agent, API: Phase 2: Infrastructure Access
    Agent->>API: GET /protected-data (Header: Bearer <JWT>)

    Note over API: TS SDK mathematically verifies<br/>EdDSA signature offline using<br/>Core's Public Key.

    alt Token Valid & Not Expired
        API-->>Agent: 200 OK (Access Granted, Returns Data)
    else Token Forged or Expired
        API-->>Agent: 401 Unauthorized (Access Denied)
    end
```

## 🚀 Quick Start

Ensure you have Go 1.21+ installed.

```bash
# Clone the repository
git clone https://github.com/agent-auth-protocol/agentauth-core.git
cd agentauth-core

# Install dependencies
go mod tidy

# Run the server
go run main.go
```

## 📖 Protocol Flow (V1)

1. **Agent Registration:** An agent generates a local Ed25519 keypair and registers its public key with the Auth Server via `POST /register`.
2. **Token Issuance:** The agent requests a session token via `POST /token`.
3. **Infrastructure Access:** The server issues a short-lived JWT. The agent uses this JWT as a Bearer token to securely interact with Model Context Protocol (MCP) servers, databases, or cloud infrastructure.

---

_Built for the Agentic Era. Part of the AgentAuth Protocol Suite._
