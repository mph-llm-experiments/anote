---
name: anote
description: Agent-first idea management using the anote CLI. Use when capturing ideas, tracking idea maturity, managing relationships between ideas, or reviewing the idea pipeline. NOT for tasks (use atask) or contacts (use apeople).
---

# anote — Agent-First Idea Management

Manage ideas using the anote CLI tool. Ideas are things that **mature** rather than things that get **done**. This is a sibling tool to atask (tasks) and apeople (contacts) — use anote for thinking, atask for doing, apeople for relationships.

**IMPORTANT: Everything is an idea.** Aspirations and beliefs are both ideas — they are managed with the same `anote` commands and stored in the same ideas directory. The `kind` field is simply a property on an idea, like `state` or `maturity`. Never create separate files or use different tools for different kinds of ideas. Always use `anote`.

## When to Use anote vs atask

| Use anote when... | Use atask when... |
|---|---|
| Capturing any kind of idea | Creating actionable work |
| Recording something the user believes to be true | Assigning tasks with deadlines |
| Exploring "what if" scenarios | Tracking task completion |
| Tracking idea maturity (crawl/walk/run) | Managing project deliverables |
| Connecting related concepts | Managing project deliverables |

## Core Concepts

### Everything is an Idea

All items managed by anote are **ideas**. Every idea has a `kind` property:

- `aspiration` (default): An idea about something to build, ship, or do
- `belief`: An idea that represents a conviction, principle, or mental model the user holds

Both kinds are created with `anote new`, listed with `anote list`, updated with `anote update`, etc. The `kind` only affects which **display labels** the CLI uses for certain states.

### Three Orthogonal Dimensions

Every idea has three properties:

**Kind** — what sort of idea (`aspiration` or `belief`)

**State** — where in its journey:
- `seed` → `draft` → `active` ↔ `iterating` → terminal state
- Terminal states: `implemented`, `archived`, `rejected`, `dropped`

**Maturity** — how baked (orthogonal to state and kind):
- `crawl` → `walk` → `run`
- Only meaningful once an idea reaches the "engaged" state

### Display Labels by Kind

The state machine is identical for all ideas. The `kind` property only changes the **display labels** for three states:

| Position | kind: aspiration | kind: belief |
|----------|-----------|--------|
| seed | seed | seed |
| draft | draft | draft |
| engaged | **active** | **considering** |
| rethinking | **iterating** | **reconsidering** |
| arrived (terminal) | **implemented** | **accepted** |
| shelved (terminal) | archived | archived |
| no (terminal) | rejected | rejected |
| fizzled (terminal) | dropped | dropped |

The CLI accepts display labels as input (e.g., `--state considering` works) and shows them in output.

### Rules to Enforce
1. **Rejected requires a reason.** Always ask the human why before rejecting.
2. **Active encourages a project link.** When an aspiration-kind idea goes active, suggest linking to an atask project.
3. **No enforced state transitions.** Any state can move to any other state. The typical flow is seed → draft → active → implemented, but the user is free to skip or jump as they see fit.
4. **Accepted ideas inform context.** When the user discusses a topic, pull ideas with `kind: belief` and state `accepted` as context that shapes advice.

## Command Reference

### Create an Idea
```bash
anote new "My idea title"                                    # kind defaults to aspiration
anote new --kind belief "Chaos-to-system is unsustainable"   # explicitly set kind
anote new --tag coaching --tag leadership "Coaching practice" # with tags
anote new --kind belief --tag work "Remote work needs trust"  # kind + tags
```
State defaults to `seed`. Kind defaults to `aspiration`. The `idea` tag is always added to the filename.

### List Ideas
```bash
anote list                          # Non-terminal ideas (all kinds)
anote list --state active           # Filter by state
anote list --state considering      # Ideas with kind=belief in the "active" position
anote list --kind belief            # Filter by kind
anote list --kind aspiration        # Only aspiration-kind ideas
anote list --maturity crawl         # Filter by maturity
anote list --tag coaching           # Filter by tag
anote list -a                       # All ideas including terminal
anote --json list                   # JSON for agent processing
```

### Show Idea Details
```bash
anote show 5                        # By index_id
anote show 20260216T103045          # By Denote ID
anote --json show 5                 # JSON output
```

### Update State, Maturity, or Kind
```bash
anote update 5 --state draft        # Progress state
anote update 5 --state considering  # Use display labels (belief kind)
anote update 5 --state accepted     # Terminal state for belief kind
anote update 5 --maturity walk      # Advance maturity
anote update 5 --kind belief        # Reclassify an idea's kind
anote update 5 --state active --maturity crawl  # Multiple at once
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
Both ideas get each other's Denote ID in their `related_ideas` array. Ideas of any kind can be linked to each other.

### Link to atask Project
```bash
anote project 5 20260215T140000     # Link to project by Denote ID
```
Adds the project's Denote ID to the idea's `related_tasks` array.

### Cross-App Relationships

Link ideas to contacts (apeople), tasks (atask), or other ideas using Denote IDs:

```bash
# Add relationships
anote update 5 --add-person 20250612T080000
anote update 5 --add-task 20250610T141230
anote update 5 --add-idea 20250607T093045

# Remove relationships
anote update 5 --remove-person 20250612T080000
anote update 5 --remove-task 20250610T141230
anote update 5 --remove-idea 20250607T093045
```

Relationship fields in frontmatter: `related_people`, `related_tasks`, `related_ideas` (arrays of Denote IDs). These are NOT automatically bidirectional — update both entities to link both directions.

## Agent Workflow Patterns

### Idea Capture Session
When the user wants to brainstorm or capture ideas:
```bash
# Capture ideas — use --kind belief when it's a conviction, otherwise default is fine
anote new --tag brainstorm "Build a mentoring platform"
anote new --kind belief --tag brainstorm "Chaos-to-system translation is unsustainable"

# Review what was captured
anote list --tag brainstorm
```

### Idea Review / Pipeline Check
When the user wants to review their ideas:
```bash
# What ideas are actively being worked on or considered?
anote --json list --state active
anote --json list --state considering

# What ideas need attention?
anote --json list --state seed

# What's being reworked?
anote --json list --state iterating
anote --json list --state reconsidering
```

### Contextual Conversations
When the user says "let's talk about [topic]":
1. Pull accepted belief-kind ideas tagged with relevant topics — these inform your advice
2. Pull considering belief-kind ideas — probe these: "You've been thinking about X — how does this land now?"
3. Pull active aspiration-kind ideas — these are in-flight work

```bash
# Pull belief-kind ideas about work for context
anote --json list --kind belief --tag work

# Pull all engaged ideas
anote --json list --state active
anote --json list --state considering
```

### Progressing an Idea
The journey is the same regardless of kind — only the labels differ:

1. `seed → draft`: User has written prose fleshing it out / articulating it clearly
2. `draft → active/considering`: Work is beginning (aspiration) or user is evaluating it (belief)
3. `active/considering → iterating/reconsidering`: Coming back to rework or re-evaluate
4. `iterating/reconsidering → active/considering`: After rework, consider advancing maturity

```bash
anote update 5 --state draft
# ... user adds content ...
anote update 5 --state active          # aspiration: work begins
anote update 5 --state considering     # belief: user is evaluating
```

For aspiration-kind ideas going active, suggest linking an atask project:
```bash
# "Consider linking an atask project: anote project 5 <project-id>"
```

For belief-kind ideas in `considering`, probe the user:
```bash
# "You've been considering X — does it still ring true?"
```

When a user says "I'm not sure about that anymore" about a belief:
```bash
anote update 5 --state reconsidering
# Probe: "What's changed? What challenged this?"
```

### Closing Ideas
```bash
# Aspiration-kind: it shipped
anote update 5 --state implemented

# Belief-kind: user commits to it as true
anote update 5 --state accepted

# Shelving for later (any kind)
anote update 5 --state archived

# Deliberate no (reason required, any kind)
anote reject 5 "Decided this doesn't align with current goals"

# Fizzled out (any kind)
anote update 5 --state dropped
```

## File Format

All ideas are Markdown files with Denote naming:
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
related_ideas: ["20260301T091500"]
related_tasks: ["20260215T140000"]
related_people: []
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

Any state can transition to any other state — there are no enforced flow rules.
The typical progression is shown above, but skipping states is allowed.
Display labels (considering, reconsidering, accepted) map to canonical states (active, iterating, implemented) and are used when an idea has `kind: belief`.
