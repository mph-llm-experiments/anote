# anote — Agent-First Idea Management

## What This Is
anote is a CLI tool for managing ideas using the Denote file naming convention. It is a sibling to [atask](../atask/), sharing file format conventions but with semantics for **thinking** rather than **doing**.

## Project Status
**Phase: v0.1.0 — Core CLI Complete**
- Spec: `docs/ANOTE_SPEC.md` (v0.1.0)
- Language: Go (to match atask)
- All core commands implemented: new, list, show, update, reject, tag, link, project
- Agent skill: `~/.claude/skills/anote/SKILL.md`
- Test suite: 40+ unit tests across denote, idea, and config packages

## Key Design Decisions
1. **Two orthogonal dimensions**: State (lifecycle) and Maturity (how baked) — see spec
2. **Denote file naming**: `YYYYMMDDTHHMMSS--title-slug__idea_tag1_tag2.md`
3. **Required `idea` tag** in filename (like `task`/`project` in atask)
4. **Rejected state requires a reason** — agent must enforce this
5. **Active state encourages atask project link** — agent should prompt
6. **Human writes prose, agent manages YAML** — never auto-modify content below frontmatter

## Architecture
```
anote/
├── main.go                          # Entry point, version, config, CLI dispatch
├── internal/
│   ├── cli/
│   │   ├── cli.go                   # Run() entry point, global flags
│   │   ├── commands.go              # Command struct, flag parsing
│   │   └── idea_commands.go         # All idea CLI commands
│   ├── config/
│   │   └── config.go                # TOML config, XDG lookup
│   ├── denote/
│   │   ├── types.go                 # File, IdeaMetadata, Idea, state/maturity
│   │   ├── parser.go                # Denote filename regex, frontmatter parsing
│   │   ├── scanner.go               # FindIdeas via glob
│   │   ├── id_counter.go            # Sequential ID counter
│   │   ├── create.go                # Slug, filename builder, file rename
│   │   └── frontmatter.go           # YAML serialization, write, update
│   └── idea/
│       ├── idea.go                  # Lookup helpers
│       └── create.go                # CreateIdea
├── docs/
│   └── ANOTE_SPEC.md
└── CLAUDE.md
```

## State Values
`seed` → `draft` → `active` ↔ `iterating` → `implemented` | `archived` | `rejected` | `dropped`

## Maturity Values (Orthogonal)
`crawl` → `walk` → `run`

## Development Protocols

### Context Protection
- All design decisions are captured in `docs/ANOTE_SPEC.md`
- This file (`CLAUDE.md`) serves as the quick-reference for any agent session
- If a design decision changes, update BOTH the spec and this file

### Before Making Changes
- Read this file and the spec before starting work
- Check git log for recent changes if resuming after a break
- Do not deviate from the spec without discussing with the human first

### Commit Conventions
- Commit messages should be clear and descriptive
- Reference spec version when implementing spec-defined behavior

### Relationship to atask
- atask lives at `../atask/`
- Shared conventions: Denote file naming, YAML frontmatter pattern, sequential ID system
- Separate counter files (`.anote-counter.json` vs `.denote-task-counter.json`)
- Cross-linking via Denote IDs in the `project` field
