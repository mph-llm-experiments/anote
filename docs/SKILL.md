---
name: anote
description: Agent-first idea management using the anote CLI. Use when capturing ideas, tracking idea maturity, managing relationships between ideas, or reviewing the idea pipeline. NOT for tasks (use atask) or contacts (use apeople).
---

# anote -- Idea Management

Manage ideas using the anote CLI. Ideas are things that **mature** rather than things that get **done**. Sibling tools: atask (tasks) and apeople (contacts).

Everything is an idea. Aspirations and beliefs are both ideas managed with the same commands and stored in the same directory. The `kind` field is a property on an idea, like `state` or `maturity`.

All data is stored as plain markdown files with YAML frontmatter. Filename format: `{ulid}--{slug}__idea.md`

## Core Concepts

### Kind

- `aspiration` (default): Something to build, ship, or do
- `belief`: A conviction, principle, or mental model

### State

Progression: `seed` -> `draft` -> `active` <-> `iterating` -> terminal

Terminal states: `implemented`, `archived`, `rejected`, `dropped`

No enforced transitions -- any state can move to any other.

### Display Labels by Kind

The `kind` only changes display labels for four states:

| Position | aspiration | belief |
|----------|-----------|--------|
| engaged | active | considering |
| rethinking | iterating | reconsidering |
| arrived | implemented | accepted |
| all others | same | same |

The CLI accepts display labels as input (e.g. `--state considering`) and shows them in output.

### Maturity (orthogonal)

`crawl` -> `walk` -> `run`

## Commands

### new -- Create an idea

```bash
anote new "Idea title"                                     # defaults to aspiration
anote new --kind belief "Remote work needs trust"          # belief kind
anote new --tag coaching --tag leadership "Coaching idea"  # with tags
```

State defaults to `seed`. The `idea` tag is always added to the filename.

### list -- List ideas

```bash
anote list --json                        # Non-terminal ideas (all kinds)
anote list --state active --json         # Filter by state
anote list --state considering --json    # Display label works too
anote list --kind belief --json          # Filter by kind
anote list --maturity crawl --json       # Filter by maturity
anote list --tag coaching --json         # Filter by tag
anote list -a --json                     # All ideas including terminal
```

### show -- Show idea details

```bash
anote show <index_id_or_ulid> --json
```

Accepts index_id (numeric) or ULID.

### update -- Update idea metadata

```bash
anote update <id> --state draft
anote update <id> --maturity walk
anote update <id> --kind belief
anote update <id> --title "New title"
anote update <id> --state active --maturity crawl    # Multiple at once
```

Options: `--state`, `--maturity`, `--kind`, `--title`

Cross-app relationship flags (values are ULIDs):
- `--add-person <ulid>` / `--remove-person <ulid>`
- `--add-task <ulid>` / `--remove-task <ulid>`
- `--add-idea <ulid>` / `--remove-idea <ulid>`

Note: `anote update` uses manual flag parsing, so `--help` does not work. The flags listed above are confirmed from source code.

### reject -- Reject an idea (reason required)

```bash
anote reject <id> "Too expensive for current budget"
```

Sets state to `rejected` and stores the reason. Never reject without a reason.

### tag -- Add or remove tags

```bash
anote tag <id> coaching              # Add tag
anote tag <id> coaching --remove     # Remove tag
```

### link -- Link two related ideas (bidirectional)

```bash
anote link <id1> <id2>
```

Both ideas get each other's ULID in their `related_ideas` array.

### project -- Link idea to an atask project

```bash
anote project <id> <project-id>
```

Adds the project ID to the idea's `related_tasks` array.

## JSON Structure

```json
{
  "id": "01KJ1KJ9CWACTZX7ATSK3AZBSE",
  "title": "Remote work needs trust as a foundation",
  "index_id": 23,
  "type": "idea",
  "tags": ["work", "management"],
  "created": "2026-02-18T21:44:48-08:00",
  "modified": "2026-02-21T19:00:40-08:00",
  "related_people": ["01KJ1KHY4NDY5FR6S9FYTHFTAV"],
  "related_tasks": [],
  "related_ideas": [],
  "file_path": "/path/to/01KJ1KJ9CW...--remote-work-needs-trust__idea.md",
  "kind": "belief",
  "state": "seed"
}
```

`anote show <id> --json` also includes a `content` field with the markdown body.

Key fields:
- `id` -- ULID, the canonical identifier
- `index_id` -- stable numeric ID for CLI commands
- `kind` -- aspiration or belief
- `state` -- current lifecycle state (uses display labels for belief kind)
- `related_people`, `related_tasks`, `related_ideas` -- arrays of ULIDs (always `[]`, never null)

## Rules for Agents

1. **Rejected requires a reason.** Always ask the human why before rejecting.
2. **Active encourages a project link.** When an aspiration goes active, suggest linking to an atask project.
3. **Accepted beliefs inform context.** When the user discusses a topic, pull `kind: belief` ideas with state `accepted` as context.
4. **No enforced transitions.** Any state can move to any other.

## Agent Workflows

### Idea capture

```bash
anote new --tag brainstorm "Build a mentoring platform"
anote new --kind belief --tag brainstorm "Trust beats verification"
```

### Pipeline review

```bash
anote list --state seed --json          # Ideas needing attention
anote list --state active --json        # In-flight aspirations
anote list --state considering --json   # Beliefs being evaluated
anote list --state iterating --json     # Being reworked
```

### Contextual conversations

When the user says "let's talk about [topic]":
1. Pull accepted beliefs tagged with the topic -- these inform your advice
2. Pull considering beliefs -- probe: "You've been thinking about X -- how does this land now?"
3. Pull active aspirations -- these are in-flight work

```bash
anote list --kind belief --tag work --json
anote list --state active --json
```

### Progressing an idea

```bash
anote update 5 --state draft              # User fleshed it out
anote update 5 --state active             # Work begins (aspiration)
anote update 5 --state considering        # Evaluating (belief)
anote update 5 --state implemented        # Shipped (aspiration)
anote update 5 --state accepted           # Committed (belief)
anote update 5 --state archived           # Shelved for later
anote update 5 --state dropped            # Fizzled out
anote reject 5 "Doesn't align with goals" # Deliberate no
```

### Cross-app: link idea to task and person

```bash
anote update 5 --add-task <task-ulid>
anote update 5 --add-person <contact-ulid>
# Link the other direction too:
atask update <task-index-id> --add-idea <idea-ulid>
apeople update <contact-index-id> --add-idea <idea-ulid>
```

## Configuration

Config: `~/.config/acore/config.toml`

```toml
[directories]
anote = "/path/to/ideas"
```

Override with `--dir` flag. Also supports `--config` for alternate config file.

## Global Options

```
--json         JSON output (always use for programmatic access)
--dir PATH     Override ideas directory
--config PATH  Use specific config file
--quiet, -q    Minimal output
--no-color     Disable color output
```
