# TOON: Our Universal Cross-Language Data Model

> **TL;DR:** TOON _is already_ the DealScale-wide schema language. It replaces disparate SQL, Protobuf, or ad-hoc models with a single, AI-native definition that generates strongly typed code for every runtime we care about (Python, Go, TypeScript, C, SQL, even automation agents).

---

## Why TOON Beats “SQL + Protobuf + X”

TOON is:

- **Language-agnostic & deterministic** – one schema, identical results across Python, Go, TypeScript, C/C++, Rust, SQL.
- **AST-safe & composable** – schemas are parsed into an AST that LLMs can reason about, diff, and autofix.
- **Validation-first** – supports enums, nested models, metadata, DAG dependencies, semantic validators, and autofix strategies.
- **Human-readable** – TOON definitions are text, so agents and humans can inspect or patch them without bespoke tooling.

Where SQL stops at persistence fields and Protobuf stops at binary transport, TOON adds:

| Capability                    | SQL | Protobuf | TOON |
| ----------------------------- | --- | -------- | ---- |
| Nested models / enums         | ⚠️ limited | ✅ | ✅ |
| Multi-agent safety (two-person rule, cross-agent validation) | ❌ | ❌ | ✅ |
| AST embeddings + semantic diffs | ❌ | ❌ | ✅ |
| Autofix / self-healing        | ❌ | ❌ | ✅ |
| Human + LLM readability       | ⚪ schemas only | ❌ | ✅ |
| Binary serialization          | ❌ | ✅ | ⚪ via plugin |

If we ever need binary-optimized transport, we can still **generate Protobuf from TOON**—but TOON stays the source of truth.

---

## What TOON Gives DealScale Today

### 1. **Single Schema → Every Language**

TOON definitions drive codegen for:

- Python dataclasses / pydantic models (FastAPI, embeddings pipeline)
- Go structs (scheduler, poster, automation workers)
- TypeScript interfaces + Zod schemas (Next.js UI + API validators)
- C/C++ structs (audio/WebRTC utilities)
- Rust types (new `rust/toon-rust` submodule)
- SQLModel / Prisma (persistence), JSON Schema, Pulsar schemas

### 2. **Agent-Ready Semantics**

- Deterministic DAGs for dependency ordering
- Safety Agent + Rollback Agent hooks (two-person approvals)
- Semantic diffs + AST embeddings for “explainable” changes
- Autofix mode so orchestration agents can repair mismatched payloads

### 3. **Human-Friendly, Reversible Format**

- Easier to edit and review than C headers or binary IDLs
- Perfect for LLM workflows—model definitions stay inspectable and diffable
- Works as the “intermediate language” between SQL, JSON, gRPC, REST, workflows, etc.

---

## Recommended Workflow

```
TOON schema (canonical source)
        ↓ codegen
+----------------------+------------------------+
| Runtime Targets      | Examples               |
+----------------------+------------------------+
| Python               | FastAPI models, pydantic, dataclasses |
| Go                   | microservice structs, validation      |
| TypeScript           | Interfaces, Zod validators, Next.js   |
| C / C++              | Header structs for media/WebRTC       |
| Rust                 | `rust/toon-rust` types & builders     |
| SQL / Prisma         | Table definitions, migrations         |
| Events / Messaging   | Pulsar, JSON schema, gRPC stubs       |
+----------------------+------------------------+
```

1. **Author / update TOON definitions**.
2. **Run codegen** (scripts forthcoming) to update language targets.
3. **Commit both the source schema + generated code** so agents and humans stay in sync.

---

## When To Add Protobuf (Optional)

- **Only** when we need compact binary serialization or existing gRPC toolchains.
- Wrap TOON → Protobuf via generator to keep TOON as the primary schema.

---

## Next Steps

- Identify which models to formalize first in TOON (e.g., `Lead`, `User`, `Workflow`, `AgentTask`, `Campaign`).
- Hook CI to validate:
  - Schema parses & round-trips
  - Generated targets compile (Go/Python/TS/C/Rust)
  - Agents run “two-person rule” checks on mutations

Ready to bootstrap? Ping with the model names (Lead, User, Workflow, etc.) and we’ll scaffold the TOON schema + generators.

