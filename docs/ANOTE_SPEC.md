# anote Specification

**Version:** 0.2.0
**Date:** 2026-02-16

## Overview

anote is an agent-first idea management CLI that treats ideas as things that **mature** rather than things that get **done**. It is a sibling to [atask](https://github.com/mikehall/atask), sharing the same Denote file naming conventions and YAML/Markdown pattern, but with semantics appropriate for **thinking** rather than **doing**.

## Principles

1. **Same file conventions as atask** — Denote file naming, YAML frontmatter, Markdown content
2. **Three orthogonal dimensions** — State (lifecycle), Maturity (how baked), Kind (aspiration vs belief)
3. **Ideas connect via relationships** — not blocking dependencies
4. **Human writes prose, agent manages metadata** — human never touches YAML directly (TUI planned for metadata editing)
5. **CLI-first** — designed as an agent skill first, human interaction tool second
6. **Lightweight** — just files on disk, no database, no MCP overhead

## File Naming Convention (Denote Format)

### Pattern
```
YYYYMMDDTHHMMSS--title-slug__idea_tag1_tag2.md
```

### Components
1. **Denote ID** (14 chars): `YYYYMMDDTHHMMSS` — timestamp of creation, canonical immutable identifier
2. **Separator**: `--` (double dash)
3. **Title Slug**: kebab-case title (spaces become hyphens, lowercase)
4. **Tag Separator**: `__` (double underscore)
5. **Tags**: underscore-separated; `idea` is the **required** tag (analogous to `task`/`project` in atask)
6. **Extension**: `.md`

### Examples
```
20260216T103045--coaching-practice-non-traditional-managers__idea_coaching_leadership.md
20260216T142200--home-lab-kubernetes-cluster__idea_homelab_infrastructure.md
20260301T091500--writing-a-book-on-engineering-management__idea_writing_career.md
```

## YAML Frontmatter

```yaml
---
title: Coaching Practice - Non-Traditional Managers
index_id: 5
type: idea
kind: aspiration
state: active
maturity: crawl
tags: [coaching, retirement-planning, leadership]
related: []           # Denote IDs of connected ideas
project: []           # Denote IDs of linked atask projects (encouraged when active)
rejected_reason:      # Required when state is "rejected"
created: 2026-02-16T10:30:45Z
modified: 2026-02-16T11:15:22Z
---
```

### Required Fields
- `title` (string): Human-readable title
- `index_id` (integer): Unique sequential ID for CLI convenience

### Core Fields
- `type` (string): Always `"idea"`
- `kind` (string): `"aspiration"` (default) or `"belief"` — what type of idea this is (see Kind below)
- `state` (string): Lifecycle state (see State below)
- `maturity` (string): How baked the idea is (see Maturity below)
- `tags` (array): Tags beyond those in the filename
- `related` (array of strings): Denote IDs of connected ideas
- `project` (array of strings): Denote IDs of linked atask projects
- `rejected_reason` (string): **Required** when state is `rejected`
- `created` (string): ISO 8601 timestamp
- `modified` (string): ISO 8601 timestamp

## State (Lifecycle)

State tracks **where an idea is in its journey**.

| State | Meaning |
|-------|---------|
| `seed` | Just captured the spark — minimal or no prose yet |
| `draft` | Human has added prose fleshing out what this could look like |
| `active` | Being worked on; linking to an atask project is strongly encouraged |
| `iterating` | Decided to come back and rework/improve; cycles back to active |
| `implemented` | The idea became real — it shipped or was realized |
| `archived` | Set aside for future reference — not dead, might revisit |
| `rejected` | Deliberate "no" — `rejected_reason` is **required** |
| `dropped` | Fizzled, lost energy, moved on — softer than rejected |

### State Transitions
```
seed → draft → active ↔ iterating → implemented
                                   → archived
                                   → rejected (reason required)
                                   → dropped
```

- `seed → draft`: When human adds substantive prose
- `draft → active`: When work begins; encourage atask project link
- `active → iterating`: When revisiting to rework/improve
- `iterating → active`: After rework, potentially at higher maturity
- `active/iterating → implemented`: Idea is realized
- `active/iterating → archived`: Set aside for later
- `active/iterating → rejected`: Deliberate no (agent must require reason from human)
- `active/iterating → dropped`: Abandoned without ceremony
- `archived → active`: Revisiting a shelved idea

## Maturity (Orthogonal to State)

Maturity tracks **how baked an idea is** — independent of lifecycle state.

| Level | Meaning |
|-------|---------|
| `crawl` | First pass, MVP thinking, getting the shape right |
| `walk` | Refined, gaining confidence, more detail and connections |
| `run` | Fully realized, polished, comprehensive |

### Maturity Rules
- Maturity is **optional** for `seed` and `draft` states (not yet meaningful)
- Maturity becomes relevant once an idea is `active`
- An idea can move through `active → iterating → active` and advance maturity each cycle
- Maturity is not strictly linear — an idea could be reassessed downward

## Kind (Orthogonal to State and Maturity)

Kind determines **what type of idea** this is. The same state machine applies to both kinds, but the display labels differ for three states.

| Kind | Description |
|------|-------------|
| `aspiration` | Something to build, ship, or do (default) |
| `belief` | Something held as true — a conviction, principle, or mental model |
| `plan` | A concrete plan with a timeframe |
| `note` | A piece of information — variable confidence, consult but verify |
| `fact` | A high-confidence ground truth — treat as authoritative |

### Simple Kinds (note, fact)

Notes and facts are **simple kinds** — they exist as reference material and skip the full lifecycle:

- **States:** Only `active` and `archived`. Default to `active` on creation (not `seed`).
- **Maturity:** Not applicable. `--maturity` is rejected for simple kinds.
- **Reject:** Not supported. Use `--state archived` instead.
- **Trust levels:** Facts are ground truth (assert directly). Notes are advisory (verify before asserting).

### Display Label Mapping

The underlying state machine is identical. Kind determines the vocabulary:

| Position | Aspiration | Belief |
|----------|-----------|--------|
| seed | seed | seed |
| draft | draft | draft |
| engaged | **active** | **considering** |
| rethinking | **iterating** | **reconsidering** |
| arrived (terminal) | **implemented** | **accepted** |
| shelved (terminal) | archived | archived |
| no (terminal) | rejected | rejected |
| fizzled (terminal) | dropped | dropped |

### Kind Rules
- Default kind is `aspiration` (backward compatible — existing files without `kind` are aspirations)
- CLI accepts display labels as input: `anote update 5 --state considering` works for beliefs
- CLI displays kind-specific labels in output
- JSON output uses kind-specific display labels for agent consumption
- Canonical (aspiration) state values are stored in YAML frontmatter regardless of kind

## Sequential ID System

### Counter File: `.anote-counter.json`
Located in the ideas directory:
```json
{
  "next_index_id": 12,
  "spec_version": "0.1.0"
}
```

- Ideas use their own counter, separate from atask
- Counter auto-created on first use
- If missing, scan all files to find highest existing ID

## Cross-Linking

### Between Ideas
The `related` field holds Denote IDs of connected ideas:
```yaml
related:
  - "20260216T103045"
  - "20260301T091500"
```

Relationships are **non-directional** and **non-blocking**. They represent conceptual connections, not dependencies.

### Between Ideas and Tasks/Projects
The `project` field holds Denote IDs of linked atask projects:
```yaml
project:
  - "20260215T140000"
```

This link is **strongly encouraged** when an idea reaches `active` state. The agent should prompt for it.

## Core Commands (Planned)

```bash
anote new "title"                        # Create new idea (state: seed, kind: aspiration)
anote new --kind belief "title"          # Create a belief
anote new --tag coaching "title"         # Create with tags
anote list                               # List ideas
anote list --state active                # Filter by state
anote list --state considering           # Filter by belief display label
anote list --kind belief                 # Filter by kind
anote list --maturity crawl              # Filter by maturity
anote list --tag coaching                # Filter by tag
anote show <id>                          # Display full idea
anote update <id> --state active         # Change state
anote update <id> --state considering    # Use display labels for beliefs
anote update <id> --maturity walk        # Change maturity
anote tag <id> tag-name                  # Add tag
anote link <id1> <id2>                   # Create relationship between ideas
anote reject <id> "reason"              # Reject with required reason
anote project <id> <project-denote-id>  # Link to atask project
```

## Configuration

Config file: `~/.config/anote/config.toml`

```toml
ideas_directory = "~/ideas"    # Required: where idea files live
editor = "vim"                 # External editor
```

## Content Format

After the YAML frontmatter, the file body is free-form Markdown written by the human. The agent does not modify content below the frontmatter unless explicitly asked.

```markdown
---
(frontmatter)
---

## The Idea

Human writes whatever they want here. Stream of consciousness,
structured arguments, links, whatever helps capture the thinking.

## Notes

More prose, refined over time.
```
