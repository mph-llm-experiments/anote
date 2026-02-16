# denote-ideas - Agent-First Idea Management

#idea #tooling #agentic-workflows #human-machine-partnership

## Concept
A skill for managing ideas/abstractions/aspirations using the same conventions as denote-tasks (file naming, YAML/Markdown pattern), but with semantics appropriate for **thinking** rather than **doing**. Agent-first interface, with optional vibecode UI layer later.

## Why This Pattern Works
- **Single responsibility:** Ideas aren't tasks — different tool, different semantics
- **Proven conventions:** Same file naming and YAML frontmatter pattern as denote-tasks
- **Human/machine partnership:** Human writes prose, agent manages metadata
- **Agent-first:** CLI for capture/query, UI is optional future enhancement
- **Follows agentic workflow guidance:** Keep it simple and scoped in

## YAML Frontmatter Pattern
```yaml
---
id: 20260216T103045
created: 2026-02-16T10:30:45Z
modified: 2026-02-16T11:15:22Z
title: Coaching Practice - Non-Traditional Managers
tags: [coaching, retirement-planning, leadership, idea]
maturity: draft  # seed | draft | developed | archived
related: []  # IDs of connected ideas
---
```

## Key Design Decisions
1. **Maturity instead of status:** Ideas mature (seed → draft → developed → archived), not "done/not done"
2. **Related instead of dependencies:** Ideas connect to each other, but not in blocking ways
3. **Same file conventions:** `YYYYMMDDTHHMMSS-title-slug.md` — works with existing tooling
4. **Human writes prose, agent manages metadata:** I never touch YAML, agent handles it via CLI

## Core Commands (Sketch)
- `denote-ideas new "title"` → creates timestamped file with YAML frontmatter
- `denote-ideas list --tag coaching` → shows all coaching-related ideas
- `denote-ideas show <id>` → displays full idea with metadata
- `denote-ideas tag <id> tag-name` → adds tags
- `denote-ideas link <id1> <id2>` → creates relationships between ideas
- `denote-ideas mature <id> draft|developed|archived` → tracks idea evolution

## What This Gives Me
- Lightweight idea capture (no basic-memory MCP overhead)
- Machine-readable structure (agents can query/organize/surface patterns)
- Natural prose space (just write)
- Future UI possibilities (vibecode when I want)
- Clean separation from tasks (ideas aren't work items)

## Status
**Not ready to dig in yet.** Capturing the spec for future exploration.

---

Related: [[denote-tasks]], [[agentic workflows]], [[human-machine partnership]], [[vibecode]]
