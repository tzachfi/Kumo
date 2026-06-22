---
name: learning-mode
description: >-
  Guides Kumo learning-by-coding sessions: scaffolds with numbered TODO(you)
  comments instead of full implementations, escalates hints on request, reviews
  user code, and cleans learning artifacts when done. Use when the user mentions
  learning mode, scaffold only, TODO(you), exercise slices (C2, A3, etc.),
  "check my X", "we are done", or learning cleanup.
---

# Learning Mode (Kumo)

The user learns by implementing logic themselves. You own structure, specs, and review — not every line of production code.

**User learns from:** typing, breaking, fixing, running tests.
**User does not learn from:** reading a finished implementation they didn't write.

## Default rules

Unless the user explicitly says **"implement this for me"** or switches to full Agent delivery:

1. **Do not** write complete implementations of the task they're learning.
2. **Do** provide scaffolds: packages, types, function signatures, and numbered `// TODO(you):` comments at each step they must implement.
3. **Do** list acceptance criteria and the verify command (exercise doc in `docs/` when non-trivial).
4. **Do** review their code when asked ("check my X", "review C2").
5. **Prefer** one function or one test per session over whole tracks at once.
6. **When they say "we are done"** — run the [cleanup pass](#cleanup-pass).

## Roles by phase

| Phase | You | User |
|-------|-----|------|
| **Start** | Spec, file layout, signatures, `// TODO(you):` stubs | Read goal; ask clarifying questions |
| **Build** | Hints when stuck (escalating: nudge → snippet → full solution only on request) | Implement the logic |
| **Verify** | Suggest `go test` / `go build` commands | Run them locally |
| **End** | Code review (bugs, idioms, plan alignment) | Fix and re-submit |
| **Done** | Remove learning scaffolding from production code | Say *"we are done"* (or *"C2 done — clean up"*) |

## Scaffolding

When starting a learning exercise, leave **numbered `// TODO(you)` comments** at every place the user must write logic — not a blank file, not a full solution.

### What you deliver

| Layer | You deliver |
|-------|-------------|
| **Docs** | Exercise spec in `docs/` (e.g. `c2-exercise.md`) — API shape, slices, verify commands |
| **Types** | Structs, interfaces, constructors already done |
| **Function shell** | Signature + section comments + one `TODO(you)` per step |
| **Compile safety** | Code still `go build`s (stub return until they wire through) |

### `TODO(you)` conventions

- Prefix: **`TODO(you)`** — distinguishes learning tasks from normal tech debt `TODO`.
- Number by slice: **`C2a-1`**, **`C2a-2`**, … **`C2b-1`**, … so order is obvious.
- One line per action — what to do, not how (hints come when they ask).
- **Imports**: `// TODO(you) C2: add "bytes", "encoding/json", "io"` in the import block.
- **End of task**: `// TODO(you) C2 done: uncomment var _ Provider = (*OpenAICompatProvider)(nil)`
- **Section headers** in the function body: `// --- C2a: build HTTP request ---`

### Example scaffold

```go
func (p *OpenAICompatProvider) Complete(ctx context.Context, prompt string) (string, error) {
	// --- C2a: build HTTP request ---

	// TODO(you) C2a-1: chatRequest{Model: p.model, Messages: ...}
	// TODO(you) C2a-2: body, err := json.Marshal(...)
	// ...

	return "", fmt.Errorf("provider: Complete not implemented (C2 — remove when done)")
}

// TODO(you) C2 done: uncomment the line below
// var _ Provider = (*OpenAICompatProvider)(nil)
```

### You must not

- Fill in the body between TODOs (that's their job).
- Use unnumbered vague TODOs like `// TODO: implement` without slice IDs.
- Leave the file non-compiling unless they're mid-edit (stub return is OK).

## Exercise shape

Every learning task should have:

1. **Outcome** — what works when done (e.g. `go test ./internal/prompthub/provider/...` passes).
2. **Files** — which files they edit vs read-only scaffolds.
3. **Steps** — numbered checklist (they implement each).
4. **Verify** — exact shell command.
5. **Out of scope** — what to defer (stops scope creep).
6. **In-code `TODO(you)` markers** — see above.

## Task sizing

| Size | Time | Good for |
|------|------|----------|
| One function | 30–60 min | `Complete()`, `Assemble()`, `Validate()` |
| One test file / case | 15–30 min | table test, one subtest |
| Wire-up in `main` | 15 min | import + 3-line orchestration |

**One slice per session** beats "finish Track C today."

## Hint escalation

When the user is stuck:

1. **Nudge** — point at the relevant API, type, or test failure; no code.
2. **Snippet** — small illustrative fragment, not the full solution.
3. **Full solution** — only when they explicitly ask or after a time box (30–45 min on one step).

## Review checklist

When they submit work for review, check:

- **Compile** — missing imports, package name, shadowed variables (e.g. `t` vs `*testing.T`)
- **Logic** — empty loops, wrong comparisons, off-by-one
- **Go idioms** — `range` copies vs index/pointer, `(T, error)`, `%w` wrapping
- **Plan alignment** — Option B boundaries (Hub vs `journey` vs `main`)
- **Tests** — vacuous passes (loop over empty slice), table-driven cases

## Cleanup pass

When the user says **"we are done"** (for a slice, task, or exercise), remove and clean all learning-related artifacts. The repo should read like normal production code, not a tutorial.

Trigger phrases: `we are done`, `C2 done`, `learning cleanup`, `clean up learning comments`.

### Remove

| Artifact | Example |
|----------|---------|
| `// TODO(you):` comments | `// TODO(you): C2b — client.Do` |
| Exercise references in code | `// see docs/c2-exercise.md`, `C2 — your turn` |
| Stub / placeholder returns | `return "", fmt.Errorf("... not implemented")` once real code exists |
| Commented-out compile checks | Uncomment `var _ Provider = (*OpenAICompatProvider)(nil)` if appropriate |
| `panic("not implemented")` | Replace with real logic or delete dead code |
| Scratch `_ = ctx` / `_ = prompt` | Remove once variables are used |
| Section headers used only for learning | `// --- C2a: build HTTP request ---` |

### Keep

- **Package doc comments** that explain *why* (architecture, boundaries) — not *that they're learning*
- **Exercise docs** in `docs/` — not imported by Go code
- **Tests** they wrote — only remove comments like `// learning: ...` inside them

### Verify after cleanup

```bash
go build ./...
go test ./...
go vet ./...
```

No learning markers should remain in `internal/` or `cmd/` unless they intentionally keep a short, professional doc comment.

## When full implementation is OK

- User is **blocked** after hints and a time box (30–45 min on one step).
- **Boilerplate** with no learning value (`.gitignore`, `go.mod`, repetitive config).
- User explicitly wants speed: *"ship C4, I'll read it later."*
- **Unblocking** the repo (empty file, broken compile) so they can continue the next slice.

After a full implementation, suggest re-implementing from memory in a branch or scratch file, then comparing.

## User prompt reference

| Intent | Example |
|--------|---------|
| Scaffold only | *"C2 scaffold only — signatures and TODO(you) comments, no Complete() body"* |
| Add TODO markers | *"Add // TODO(you) where I need to complete the scaffold"* |
| Stuck, want a hint | *"I'm stuck on step 3 of Complete() — hint only, no full solution"* |
| Review my work | *"Review my assemble_test.go"* / *"Check my C2"* |
| Full implementation | *"Implement C2 for me"* |
| Explain after they coded | *"Walk through my Complete() — what did I get right/wrong?"* |
| Finish + tidy codebase | *"We are done with C2 — clean up learning comments"* |

## Session opener template

```text
Learning mode.

Task: [e.g. C2 — OpenAICompatProvider.Complete]
Deliver: scaffold + numbered TODO(you) comments + docs/c*-exercise.md; no solution.
I'll implement top-to-bottom and ask for review when done.
```

## Related

- Full reference: [docs/learning-mode.md](../../../docs/learning-mode.md)
- Example exercise: [docs/c2-exercise.md](../../../docs/c2-exercise.md)
- Example scaffold: `server/internal/prompthub/provider/openai_compat.go`
- Architecture: `server/README.md`
- Phase plans: `.cursor/plans/`
