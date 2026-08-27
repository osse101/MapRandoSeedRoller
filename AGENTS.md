# AGENTS.md — MapRandoSeedRoller

A Go serverless API deployed on Vercel. It waits for an API command, turns a chat-style
command string (`"s5 RDS:MORPH"`) into a customized MapRando settings JSON, POSTs it to
maprando.com, and returns the seed URL/hash in the documented response structure.

---

## Deployment model

Vercel's Go runtime turns **every `api/*.go` file with an exported
`func(http.ResponseWriter, *http.Request)` into its own function**. Adding a file to `api/`
adds an endpoint; there is no router and no `main`.

| File                  | Handler                | Notes                                            |
| --------------------- | ---------------------- | ------------------------------------------------ |
| `api/roll.go`         | `RandomizeHandler`     | Main entrypoint. GET returns a liveness string.  |
| `api/secure_roll.go`  | `InertiaWebhookHandler`| Same payload, gated by Svix signature verification. |
| `api/inertia_roll.go` | `InertiaHandler`       | GET only, gated on an exact `User-Agent` match.  |

Shared helpers (`DecodeRequest`, `WriteJSONResponse`) live in `lib/httpio`, not in
`api/roll.go` itself — `vercel dev`'s local Go builder compiles each `api/*.go` file in
isolation (plus a generated entrypoint) and does not see unexported declarations from
sibling files in the same package, so cross-handler helpers must be a real importable
package, not same-directory unexported functions.

---

## The inbound contract

**Request** — `models.RequestRaw` (`lib/models/requests.go:5`)

| Field    | Type              | Notes                                     |
| -------- | ----------------- | ----------------------------------------- |
| `action` | string            | `roll`, `unlock`, `help`                  |
| `event`  | string            | `seed.finished`, `system.alert`           |
| `source` | string            | e.g. `racetime`, `inertia`                |
| `data`   | `json.RawMessage` | shape depends on the action               |

Exactly one of `action` or `event` should be set. `action` wins if both are present
(`workflow.Process` checks it first). Neither → error.

Per-action `data`:
- `roll` — a JSON **string** of preset + flags, e.g. `"s5 RDS:MORPH"`.
- `unlock` — a JSON **string** holding the seed URL.
- `help` — ignored; `GetHelp(req.Source)` currently ignores `source` too.

Unknown events are answered `success` on purpose, to stop webhook retries
(`lib/workflow/manager.go:54`).

**Response** — `models.ResponseOut`

| Field     | Type          | Notes                                            |
| --------- | ------------- | ------------------------------------------------ |
| `status`  | string        | `"success"` or `"error"`                         |
| `message` | string        | omitempty                                        |
| `data`    | `interface{}` | omitempty; `models.SeedData` for `roll`          |

`models.SeedData` is `{seed_url, seed_hash}`. `models.InertiaResponseOut` is a separate,
flattened shape (`url`, `hash`, `message`) used only by `api/inertia_roll.go`.

`models.ResponseMapRando` is the *upstream* shape from maprando.com — its `seed_url` is a
path, not a full URL; `MakeRequest` prefixes the base site when it starts with `/`.

---

## Dataflow

```
handler → decode → workflow.Process (lib/workflow/manager.go:10)
  → handleAction "roll" → ExecuteRoll (lib/workflow/roll.go:15)
      → PrepareGameData: split "<preset> <flags>", default preset "s5"
          → lib.MergeAndSortAliases  (lib/helpers.go:22)      alias table, longest-first
          → parser.Lex               (lib/parser/lexer.go:10) → []models.Token
          → preset.LoadTemplate      (preset/presets.go:36)   embedded JSON → map
          → parser.Hydrate           (lib/parser/mapping.go:15) → ([]byte, isDev, error)
      → randomize.Randomize          (lib/randomize/service.go:8)
          → MakeRequest: multipart POST {base}/randomize (lib/randomize/networking.go:17)
      → models.SeedData → models.ResponseOut
```

`unlock` is a separate path: `ExecuteUnlock` (`lib/workflow/unlock.go:18`) derives the seed ID
from the URL and form-POSTs `spoiler_token` to `{seedURL}/unlock`.

### Package map

| Package             | Responsibility                                                       |
| ------------------- | -------------------------------------------------------------------- |
| `api`               | HTTP entrypoints, decoding, JSON response writing                     |
| `lib/workflow`      | Action/event dispatch and orchestration (`manager.go` is the switch)  |
| `lib/parser`        | `Lex` (flags → tokens) and `Hydrate` (tokens → settings JSON)         |
| `lib/models`        | Wire types, `PresetFields`, `TriState`, alias tables                  |
| `lib/randomize`     | Outbound HTTP to maprando.com (multipart construction)                |
| `lib`               | Env/site helpers, alias table merging                                 |
| `preset`            | Embedded preset templates and name → filename mapping                 |

---

## Preset system

- `preset/data/*.json` are embedded with `go:embed data/*.json`. `preset.TemplateMap`
  (`preset/presets.go:14`) maps short names (`s5`, `expert`, `ammo-balance`, …) to filenames;
  `GetPresetNames()` is what `PrepareGameData` validates against.
- Template top-level sections: `version`, `name`, `skill_assumption_settings`,
  `item_progression_settings`, `quality_of_life_settings`, `objective_settings`, `map_layout`,
  `doors_settings`, `start_location_settings`, `save_animals`, `other_settings`.
- `models.PresetFields` (`lib/models/preset_fields.go:3`) carries `path:"..."` dotted-path tags
  mirroring those sections. **This is the extension point for new overrides**: add the field
  with its path tag, route it in `tokensToPresetFields`, apply it in `applyPresetFields`.

`Hydrate` runs three steps:
1. `tokensToPresetFields` — routes single-value flags by `tok.ID` (`race_mode`, `version`,
   `escape_timer_multiplier`) and multi-value flags by the sticky `tok.Flag` (`o`, `s`, `l`).
2. `applyPresetFields` — scalars via `SetNestedValue` (dotted path, creates missing maps);
   slices **merged into** the template's existing arrays, matched on `item` / `objective`.
   Overriding a sub-setting also nils out that section's sibling `preset` key so MapRando
   honors the raw values instead of the named preset.
3. `postprocess` → `clampObjectiveCounts` sets `min_objectives`/`max_objectives` to the number
   of `"Yes"` objectives, clamped to 1..19.

> **Silent no-op rule:** `mergeStartingItems` / `mergeObjectiveOptions` do nothing when the key
> isn't found in the template array. A misspelled alias produces no error and no effect —
> which is exactly how the objective typos below went unnoticed.

`Hydrate` also sets `name` to `"Custom"` whenever any token was parsed.

### Flag mini-language (`lib/parser/lexer.go`)

- Alias table is merged and sorted **longest-first**, then matched as a case-insensitive prefix
  scan over the input.
- Casing determines `TriState`: `ALLCAPS` → `True`, `alllower` → `False`, `MiXeD` → `Maybe`.
- Trailing digits/dots after a match are captured into `RawValue` (`X:1.5`, `S:ETANK3`).
- A one-character match sets the **sticky `lastFlag`**, which routes every following value
  token until the next flag. This is why `S:MORPHVARIA` works.
- Unmatched characters are skipped one byte at a time — **typos never error**, they vanish.

---

## Commands

```
make test    # go test ./... -v
make lint    # golangci-lint via tools/go.mod (separate module)
make fmt     # go fmt ./...
make clean   # go clean -testcache
go build ./...
```

## Conventions

- Imports grouped by gci: standard → default → `prefix(maprandoseedroller)`.
- `lib/parser/mapping.go` is excluded from linting (`.golangci.yml`).
- Logging is `slog`; the JSON handler is installed by the blank import
  `_ "maprandoseedroller/lib/logger"` in `api/roll.go` / `api/inertia_roll.go`.
- Errors are wrapped with `%w`.
- Tests are table-driven. `lib/parser` compares against golden files in
  `testdata/mapping_test/` (`template.json` in, `simple_case.json` expected).

## Environment

| Variable              | Used by                                    |
| --------------------- | ------------------------------------------ |
| `SPOILER_TOKEN`       | `lib.BuildSpoilerToken` — roll and unlock  |
| `SVIX_INERTIA_SECRET` | `api/secure_roll.go` signature verification |

`lib.BuildSite(isDev)` switches between `https://maprando.com` and `https://dev.maprando.com`.
`isDev` is set when the `D` flag pushes `version` to `preset.DevVersion` (variable).

---

## Known issues / active work areas

Recorded deliberately — these are **not** fixed. Verified against the tree at the time of writing.

**Response / dataflow (`api/roll` + `models.ResponseOut`)**
- `roll` never populates `ResponseOut.Message`, though the README documents one.

**Preset customization**
- `PresetFields` declares `SaveAnimals`, `WallJump`, `FreeShinesparks`, `SplitSpeed`,
  `MinObjectives`, `MaxObjectives` and `applyPresetFields` honors most of them, but no flag
  alias populates them — they are reachable only by extending the alias tables.
- `preset/data/Community_Race_Season_2.json` (v115) is the only template still using
  `"map_layout": "Tame"`; every other template uses `"Standard"` and the current MapRando
  value is `"Small"`. Left as-is: the value is correct for that template's version.
- Alias tables now cover every `objective` and `starting_items` name present in
  `preset/data/*.json` in both directions. If a template is added or MapRando renames a
  value, re-check with a merge of both sets before assuming a flag works — unmatched keys
  are silent no-ops.

**Coverage** (also in the gitignored `todo.txt`)
- `lib/randomize` has no tests — neither multipart construction nor response decoding.
  `randomize.HTTPClient` (`lib/randomize/networking.go`) is now an exported package var
  that can be swapped for a client with a fake `http.RoundTripper`, so this no longer
  needs a new seam, just the tests themselves — see `api/api_test.go` for the pattern.
- No tests for Svix signature verification in `api/secure_roll.go` (valid vs spoofed headers).
