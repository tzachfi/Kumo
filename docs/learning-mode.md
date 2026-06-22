# Learning Mode — How to Work With the AI on Kumo

Use this doc when you want to **learn by coding**, not watch the AI implement everything. Attach it to a chat (`@docs/learning-mode.md`) or paste the summary at the start of a session.

---

## Goal

Build real skills (Go, architecture, testing) by doing the **middle** of each task yourself. The AI owns structure, specs, and review — not every line of production code.

**You learn from:** typing, breaking, fixing, running tests.  
**You don't learn from:** reading a finished implementation you didn't write.

---

## Roles

| Phase | AI | You |
|-------|----|-----|
| **Start** | Spec, file layout, signatures, `// TODO(you):` stubs | Read the goal; ask clarifying questions |
| **Build** | Hints when stuck (escalating: nudge → snippet → full solution only on request) | Implement the logic |
| **Verify** | Suggest `go test` / `go build` commands | Run them locally |
| **End** | Code review (bugs, idioms, plan alignment) | Fix and re-submit |
| **Done** | Remove learning scaffolding from production code | Say *"we are done"* (or *"C2 done — clean up"*) |

---

## When you're done — cleanup pass

When you say **"we are done"** (for a slice, task, or exercise), the AI should **remove and clean all learning-related artifacts** from the codebase. The repo should read like normal production code, not a tutorial.

### You say

```text
We are done with C2 — clean up learning comments.
```

(or `C2 done`, `learning cleanup`, etc.)

### AI removes

| Artifact | Example |
|----------|---------|
| `// TODO(you):` comments | `// TODO(you): C2b — client.Do` |
| Exercise references in code | `// see docs/c2-exercise.md`, `C2 — your turn` |
| Stub / placeholder returns | `return "", fmt.Errorf("... not implemented")` once real code exists |
| Commented-out compile checks | Uncomment `var _ Provider = (*OpenAICompatProvider)(nil)` if appropriate, delete the comment wrapper |
| `panic("not implemented")` | Replace with real logic or delete dead code |
| Scratch `_ = ctx` / `_ = prompt` | Remove once variables are used |

### AI keeps

- **Package doc comments** that explain *why* (architecture, boundaries) — not *that you're learning*
- **Exercise docs** in `docs/` (e.g. `c2-exercise.md`) — optional to archive or delete later; they are not imported by Go code
- **Tests** you wrote — they stay; only remove comments like `// learning: ...` inside them

### Verify after cleanup

```bash
go build ./...
go test ./...
go vet ./...
```

No learning markers should remain in `internal/` or `cmd/` unless you intentionally keep a short, professional doc comment.

---

## Default rule for the AI

Unless you explicitly say **"implement this for me"** or switch to full Agent delivery:

1. **Do not** write complete implementations of the task you're learning.
2. **Do** provide scaffolds: packages, types, function signatures, and [numbered `TODO(you)` comments](#in-code-todoyou-scaffolds) at each step you must implement.
3. **Do** list acceptance criteria and the verify command (exercise doc in `docs/` when the task is non-trivial).
4. **Do** review your code when you ask ("check my X", "review C2").
5. **Prefer** one function or one test per session over whole tracks at once.
6. **When you say "we are done"** — run the [cleanup pass](#when-youre-done--cleanup-pass): strip `TODO(you)`, exercise references, stubs, and learning comments from `internal/` and `cmd/`.

---

## Prompts you can use

| Intent | Example |
|--------|---------|
| Scaffold only | *"C2 scaffold only — signatures and TODO(you) comments, no Complete() body"* |
| Add TODO markers | *"Add // TODO(you) where I need to complete the scaffold"* |
| Stuck, want a hint | *"I'm stuck on step 3 of Complete() — hint only, no full solution"* |
| Review my work | *"Review my assemble_test.go"* / *"Check my C2"* |
| Full implementation | *"Implement C2 for me"* (use when blocked or out of time) |
| Explain after you coded | *"Walk through my Complete() — what did I get right/wrong?"* |
| Finish + tidy codebase | *"We are done with C2 — clean up learning comments"* |

---

## Task sizing

| Size | Time | Good for |
|------|------|----------|
| One function | 30–60 min | `Complete()`, `Assemble()`, `Validate()` |
| One test file / case | 15–30 min | table test, one subtest |
| Wire-up in `main` | 15 min | import + 3-line orchestration |

**One slice per session** beats "finish Track C today."

---

## Standard exercise shape

Every learning task should have:

1. **Outcome** — what works when done (e.g. `go test ./internal/prompthub/provider/...` passes).
2. **Files** — which files you edit vs which are read-only scaffolds.
3. **Steps** — numbered checklist (you implement each).
4. **Verify** — exact shell command.
5. **Out of scope** — what to defer (stops scope creep).
6. **In-code `TODO(you)` markers** — see [below](#in-code-todoyou-scaffolds).

---

## In-code `TODO(you)` scaffolds

When starting a learning exercise, the AI should leave **numbered `// TODO(you)` comments** in the code at every place **you** must write logic — not a blank file, not a full solution.

### What the AI provides

| Layer | AI delivers |
|-------|-------------|
| **Docs** | Exercise spec in `docs/` (e.g. `c2-exercise.md`) — API shape, slices, verify commands |
| **Types** | Structs, interfaces, constructors already done (e.g. `OpenAICompatConfig`, `chatRequest`) |
| **Function shell** | Signature + section comments + one `TODO(you)` per step |
| **Compile safety** | Code still `go build`s (stub return until you wire through) |

### `TODO(you)` conventions

- Prefix: **`TODO(you)`** — distinguishes learning tasks from normal tech debt `TODO`.
- Number by slice: **`C2a-1`**, **`C2a-2`**, … **`C2b-1`**, … so order is obvious.
- One line per action — what to do, not how (hints come when you ask).
- **Imports**: `// TODO(you) C2: add "bytes", "encoding/json", "io"` in the import block.
- **End of task**: `// TODO(you) C2 done: uncomment var _ Provider = (*OpenAICompatProvider)(nil)`
- **Section headers** in the function body: `// --- C2a: build HTTP request ---`

### Example (inside `Complete()`)

```go
func (p *OpenAICompatProvider) Complete(ctx context.Context, prompt string) (string, error) {
	// --- C2a: build HTTP request ---

	// TODO(you) C2a-1: chatRequest{Model: p.model, Messages: ...}

	// TODO(you) C2a-2: body, err := json.Marshal(...)

	// TODO(you) C2a-3: url := p.baseURL + "/chat/completions"

	// TODO(you) C2a-4: req, err := http.NewRequestWithContext(ctx, ...)

	// TODO(you) C2a-5: set Content-Type and Authorization headers

	// --- C2b: execute request, read response ---

	// TODO(you) C2b-1: resp, err := p.client.Do(req)
	// ...

	// --- C2c: parse JSON, return assistant text ---

	// TODO(you) C2c-1: json.Unmarshal into chatResponse
	// ...

	return "", fmt.Errorf("provider: Complete not implemented (C2 — remove when done)")
}

// TODO(you) C2 done: uncomment the line below
// var _ Provider = (*OpenAICompatProvider)(nil)
```

Work **top to bottom**; delete each `TODO(you)` as you complete that step.

### Prompt to request this layout

```text
Add // TODO(you) comments where I need to complete the scaffold for C2.
```

(or include in session opener: *"Deliver: scaffold + numbered TODO(you) comments"*)

### AI must not

- Fill in the body between TODOs (that's your job).
- Use unnumbered vague TODOs like `// TODO: implement` without slice IDs.
- Leave the file non-compiling unless you're mid-edit (stub return is OK).

### Removed on cleanup

All `TODO(you)` lines, section headers used only for learning, stub returns, and commented compile checks — see [cleanup pass](#when-youre-done--cleanup-pass).

---

## Example: breaking C2 into slices

Instead of "implement HTTP provider":

```text
C2a  Build request JSON + http.NewRequestWithContext     (you)
C2b  client.Do, status check, io.ReadAll                  (you)
C2c  Unmarshal response, return choices[0].message.content (you)
     → AI reviews when all three are done
```

Same pattern works for any track: **A3** was three tests; **C4** can be mock branch first, then llm branch.

---

## Review checklist (what the AI should look for)

When you submit work for review, expect feedback on:

- **Compile** — missing imports, package name, shadowed variables (e.g. `t` vs `*testing.T`)
- **Logic** — empty loops, wrong comparisons, off-by-one
- **Go idioms** — `range` copies vs index/pointer, `(T, error)`, `%w` wrapping
- **Plan alignment** — Option B boundaries (Hub vs `journey` vs `main`)
- **Tests** — vacuous passes (loop over empty slice), table-driven cases

---

## When full AI implementation is OK

- You're **blocked** after hints and a time box (e.g. 30–45 min on one step).
- **Boilerplate** with no learning value (`.gitignore`, `go.mod`, repetitive config).
- You explicitly want speed: *"ship C4, I'll read it later."*
- **Unblocking** the repo (empty file on disk, broken compile) so you can continue learning the next slice.

After a full implementation, consider **re-implementing from memory** in a branch or scratch file, then compare.

---

## Verify habits (always you)

```bash
go build ./...
go test ./path/to/package/... -v
go vet ./...
```

Run after every small change. The compiler and tests are your feedback loop.

---

## Session opener (copy into chat)

```text
Learning mode (@docs/learning-mode.md).

Task: [e.g. C2 — OpenAICompatProvider.Complete]
Deliver: scaffold + numbered TODO(you) comments + docs/c*-exercise.md; no solution.
I'll implement top-to-bottom and ask for review when done.
```

---

## Related

- Phase plans: `.cursor/plans/` (Kumo Phase 1, Phase 2)
- Architecture: `server/README.md`
- Env/secrets (not in git): `server/.env.example`, `internal/prompthub/secrets/`
- Example exercise + TODO layout: `docs/c2-exercise.md`, `server/internal/prompthub/provider/openai_compat.go`
