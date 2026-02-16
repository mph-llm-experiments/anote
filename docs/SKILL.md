---
name: anote
description: Agent-first idea management using the anote CLI. Use when capturing ideas, tracking idea maturity, managing relationships between ideas, or reviewing the idea pipeline. NOT for tasks — use denote-tasks for actionable work items.
---

# anote — Agent-First Idea Management

Manage ideas using the anote CLI tool. Ideas are things that **mature** rather than things that get **done**. This is a sibling tool to denote-tasks (atask) — use anote for thinking, atask for doing.

## When to Use anote vs atask

| Use anote when... | Use atask when... |
|---|---|
| Capturing a thought or aspiration | Creating actionable work |
| Recording a belief or conviction | Assigning tasks with deadlines |
| Exploring "what if" scenarios | Tracking task completion |
| Tracking idea maturity (crawl/walk/run) | Managing project deliverables |
| Connecting related concepts | Managing project deliverables |

## Core Concepts

### Three Orthogonal Dimensions

**Kind** — what type of idea:
- `aspiration` (default): Something to build/do. States display as: active, iterating, implemented
- `belief`: Something held as true. States display as: considering, reconsidering, accepted

**State** — where in the journey:
- `seed` → `draft` → `active` ↔ `iterating` → terminal state
- Terminal states: `implemented`, `archived`, `rejected`, `dropped`
- For beliefs, display labels differ (see table below)

**Maturity** — how baked (orthogonal to state and kind):
- `crawl` → `walk` → `run`
- Only meaningful once an idea is `active`/`considering`

### Display Label Mapping

The state machine is identical for both kinds. Only the display labels change:

| Position | Aspiration | Belief |
|----------|-----------|--------|
| seed | seed | seed |
| draft | draft | draft |
| engaged | **active** | **considering** |
| rethinking | **iterating** | **reconsidering** |
| arrived | **implemented** | **accepted** |
| shelved | archived | archived |
| no | rejected | rejected |
| fizzled | dropped | dropped |

CLI accepts display labels as input (e.g., `--state considering` works for beliefs).

### Rules to Enforce
1. **Rejected requires a reason.** Always ask the human why before rejecting.
2. **Active encourages a project link.** When an aspiration goes active, suggest linking to an atask project.
3. **State transitions are validated.** You cannot skip states (e.g., seed → active is invalid, must go seed → draft → active).
4. **Beliefs inform context.** When the user discusses a topic, pull accepted beliefs as context that shapes advice.

## Command Reference

### Create an Idea
```bash
anote new "My idea title"
anote new --kind belief "Chaos-to-system translation is unsustainable"
anote new --tag coaching --tag leadership "Coaching practice"
anote new --kind belief --tag work "Remote work requires trust"
```
State defaults to `seed`. Kind defaults to `aspiration`. The `idea` tag is always added to the filename.

### List Ideas
```bash
anote list                          # Non-terminal ideas
anote list --state active           # Filter by state
anote list --state considering      # Beliefs in the "active" position
anote list --kind belief            # Filter by kind
anote list --kind aspiration        # Only aspirations
anote list --maturity crawl         # Filter by maturity
anote list --tag coaching           # Filter by tag
anote list -a                       # All ideas including terminal
anote --json list                   # JSON for agent processing
anote --json list --kind belief     # JSON beliefs for agent context
```

### Show Idea Details
```bash
anote show 5                        # By index_id
anote show 20260216T103045          # By Denote ID
anote --json show 5                 # JSON output
```

### Update State or Maturity
```bash
anote update 5 --state draft        # Progress state
anote update 5 --state considering  # Use belief display labels
anote update 5 --state accepted     # Belief reaches terminal
anote update 5 --maturity walk      # Advance maturity
anote update 5 --state active --maturity crawl  # Both at once
```

### Reject an Idea (Reason Required)
```bash
anote reject 5 "Too expensive for current budget"
```
The reason is stored in the `rejected_reason` field. **Never reject without a reason.**

### Add/Remove Tags
```bash
anote tag 5 coaching                # Add tag
anote tag 5 coaching --remove       # Remove tag
```
Tags update both the filename and frontmatter.

### Link Related Ideas
```bash
anote link 5 8                      # Bidirectional link
```
Both ideas get each other's Denote ID in their `related` array.

### Link to atask Project
```bash
anote project 5 20260215T140000     # Link to project by Denote ID
```

## Agent Workflow Patterns

### Idea Capture Session
When the user wants to brainstorm or capture ideas:
```bash
# Capture aspirations
anote new --tag brainstorm "Idea title"

# Capture beliefs
anote new --kind belief --tag brainstorm "Belief statement"

# Review what was captured
anote list --tag brainstorm
```

### Idea Review / Pipeline Check
When the user wants to review their ideas:
```bash
# What aspirations are active?
anote --json list --state active --kind aspiration

# What beliefs are accepted? (use as context)
anote --json list --state accepted

# What beliefs are being reconsidered?
anote --json list --state reconsidering

# What seeds need attention?
anote --json list --state seed

# What's iterating?
anote --json list --state iterating
```

### Contextual Conversations
When the user says "let's talk about [topic]":
1. Pull accepted beliefs tagged with relevant topics — these inform your advice
2. Pull considering beliefs — probe these: "You've been thinking about X — how does this land now?"
3. Pull active aspirations — these are in-flight work

```bash
# Pull beliefs about work
anote --json list --kind belief --tag work

# Pull all active things
anote --json list --state active
anote --json list --state considering
```

### Maturing an Aspiration
When an aspiration progresses:
1. `seed → draft`: User has written prose fleshing it out
2. `draft → active`: Work is beginning — suggest linking an atask project
3. `active → iterating`: Coming back to rework — may advance maturity
4. `iterating → active`: After rework, consider advancing maturity (crawl → walk → run)

```bash
anote update 5 --state draft
# ... user adds content ...
anote update 5 --state active
# "Consider linking an atask project: anote project 5 <project-id>"
anote update 5 --maturity crawl
```

### Maturing a Belief
When a belief progresses:
1. `seed → draft`: User has articulated the belief clearly
2. `draft → considering`: User is actively evaluating whether they accept it
3. `considering → reconsidering`: Something challenged the belief
4. `reconsidering → considering`: Re-evaluated, still holding
5. `considering → accepted`: User commits to this as true

```bash
anote update 5 --state draft
# ... user refines the belief ...
anote update 5 --state considering
# Agent can probe: "You've been considering X — does it still ring true?"
anote update 5 --state accepted
```

When a user says "I'm not sure about that anymore":
```bash
anote update 5 --state reconsidering
# Probe: "What's changed? What challenged this belief?"
```

When a user says "I don't believe that anymore":
```bash
anote reject 5 "Reason for no longer believing"
```

### Closing Ideas
```bash
# Aspiration shipped
anote update 5 --state implemented

# Belief accepted
anote update 5 --state accepted

# Shelving for later
anote update 5 --state archived

# Deliberate no (reason required)
anote reject 5 "Decided this doesn't align with current goals"

# Fizzled out
anote update 5 --state dropped
```

## File Format

Ideas are Markdown files with Denote naming:
```
YYYYMMDDTHHMMSS--title-slug__idea_tag1_tag2.md
```

YAML frontmatter managed by the agent:
```yaml
---
title: Coaching Practice
index_id: 5
type: idea
kind: aspiration
state: active
maturity: crawl
tags: [coaching, leadership]
related: ["20260301T091500"]
project: ["20260215T140000"]
created: "2026-02-16T10:30:45-08:00"
modified: "2026-02-16T11:15:22-08:00"
---
```

Content below the frontmatter is free-form Markdown written by the human.

## Configuration

Config file: `~/.config/anote/config.toml`
```toml
ideas_directory = "~/ideas"
editor = "vim"
```

Override with `--dir` flag: `anote --dir /path/to/ideas list`

## Valid State Transitions

```
seed → draft
draft → active
active → iterating | implemented | archived | rejected | dropped
iterating → active | implemented | archived | rejected | dropped
archived → active
```

Terminal states (implemented, rejected, dropped) have no outbound transitions.
Display labels (considering, reconsidering, accepted) map to canonical states (active, iterating, implemented).
