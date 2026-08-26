# Seed Roller

A bot and API for rolling randomizer seeds, with support for presets, overrides, and race mode. The bot holds the spoiler token and can unlock the spoiler log upon request.

---

## Table of Contents

- [User Commands](#user-commands)
  - [Presets](#presets)
  - [Value Overrides](#value-overrides)
  - [Override Values Reference](#override-values-reference)
- [API Reference](#api-reference)
  - [POST /api/roll](#post-apiroll)
  - [Actions](#actions)

---

## User Commands

```
!roll <args>
```

| Example              | Result                                                     |
| -------------------- | ---------------------------------------------------------- |
| `!roll`              | Rolls a seed using the current season's preset             |
| `!roll s5`           | Rolls a seed using the `s5` preset                         |
| `!roll s5 RDS:MORPH` | Rolls a seed using `s5` with the specified field overrides |

---

### Presets

| Preset         | Difficulty | Logic     | Notes                                                        |
| -------------- | ---------- | --------- | ------------------------------------------------------------ |
| `s2`           | Hard       | Tricky    |                                                              |
| `s3`           | Hard       | Tricky    |                                                              |
| `s3a`          | Hard       | Tricky    | Save the animals                                             |
| `s4`           | Hard       | Tricky    |                                                              |
| `s5`           | Hard       | Tricky    | Current season default                                       |
| `default`      | Normal     | Basic     | This is Map Rando's introductory preset                      |
| `expert`       | Expert     | Challenge | Based on s5                                                  |
| `mentor`       | Medium     | Basic     |                                                              |
| `objectives`   | Hard       | Tricky    | Same as s4; objectives set via args                          |
| `suits`        | Hard       | Tricky    | Starting Gravity + Varia; rest is s4                         |
| `g91`          | Hard       | Tricky    | Starting Gravity + 9 E-Tanks + 1 Reserve; rest is s4         |
| `draft`        | Hard       | Tricky    | Same as s4; starting items set via args                      |
| `metroids`     | Hard       | Tricky    | Metroid objectives, no Mother Brain 2; rest is s4            |
| `noobjectives` | Hard       | Tricky    | No objectives; rest is s4                                    |
| `vmove`        | Hard       | Tricky    | Starting Varia, Grapple, HiJump, Ice, Springball; rest is s4 |
| `ammo-balance` | Hard       | Tricky    | Ammo-balanced tournament settings                            |
| `nis`          | Very Hard  | Random    | NIS Very Hard, randomly selected difficulty and sprite       |

---

### Value Overrides

Overrides are passed as flags after the preset name. Letter case controls the override behaviour:

| Case        | Meaning |
| ----------- | ------- |
| `UPPERCASE` | Enable  |
| `lowercase` | Disable |
| `miXEDcasE` | Maybe   |

**Available flags:**

| Flag             | Description             |
| ---------------- | ----------------------- |
| `R`              | Race mode               |
| `D`              | Use dev site            |
| `X:<value>`      | Escape timer multiplier |
| `S:<items>`      | Starting items          |
| `O:<objectives>` | Objectives              |
| `L:<layout>`     | Map layout              |

---

### Override Values Reference

**Objectives**

```
kraid, phan, dray, ridley, spore, croc, bot, gt, bt, bowling, acid, pit, babyk, plasma, metal, m1, m2, m3, m4
```

**Starting Items**

```
missile, etank, rtank, super, pb, charge, ice, wave, spazer, plasma, xray, morph, bomb, grapple, hjb, speed, spring, space, screw, varia, gravity, wj, blue, spark
```

**Map Layouts**

```
vanilla, small, standard, wild
```

---

## API Reference

The Seed Roller API accepts requests as JSON payloads.

### POST /api/roll

Main endpoint for interacting with the API. Also available at `/api/secure_roll` for Svix-verified webhooks.

**Request**

| Field    | Type   | Description                                        |
| -------- | ------ | -------------------------------------------------- |
| `action` | string | Action to perform (e.g., `roll`, `unlock`, `help`) |
| `event`  | string | Acknowledge an event (e.g., `seed.finished`)       |
| `source` | string | The source identifier for the request              |
| `data`   | any    | Payload specific to the action or event            |

_Note: Exactly one of `action` or `event` should be provided._

**Response**

| Field     | Type   | Description                                                      |
| --------- | ------ | ---------------------------------------------------------------- |
| `status`  | string | `"success"` or `"error"`                                         |
| `message` | string | Optional status or error message                                 |
| `data`    | object | Optional action-specific response data (e.g., seed URL and hash) |

---

### Actions

#### `roll`

Rolls a new seed using the provided preset and flags.

**Request `data`**
A string specifying seed rolling parameters, such as the preset and field override flags. (e.g., `"s5 RDS:MORPH"`)

**Response `data`**

| Field       | Type   | Description                                                        |
| ----------- | ------ | ------------------------------------------------------------------ |
| `seed_url`  | string | The URL of the generated seed                                      |
| `seed_hash` | string | The in-game hash code for the seed                                 |
| `extra`     | object | Optional, preset-specific extra data (see below); omitted if none  |

Some presets run custom logic on top of the normal roll flow and attach an `extra`
payload whose shape depends on the preset. Callers should only interpret `extra` when
they recognize the preset that was rolled.

| Preset | `extra` shape               | Notes                                                                          |
| ------ | --------------------------- | ------------------------------------------------------------------------------ |
| `nis`  | `{ "sprite_name": string }` | Also randomizes the item progression difficulty (Technical/Challenge/Desolate) |

**Example**

```json
// Request
{
  "action": "roll",
  "source": "racetime",
  "data": "s5 RDS:MORPH"
}

// Response
{
  "status": "success",
  "data": {
    "seed_url": "https://maprando.com/seed/tc2pHBSZc/",
    "seed_hash": "YARD YARD YARD YARD"
  }
}
```

**Example (preset with `extra` data)**

```json
// Request
{
  "action": "roll",
  "source": "racetime",
  "data": "nis"
}

// Response
{
  "status": "success",
  "data": {
    "seed_url": "https://maprando.com/seed/6jZq6Jxcg/",
    "seed_hash": "OWTCH ZEBBO TATORI EVIR",
    "extra": {
      "sprite_name": "Dread Samus"
    }
  }
}
```

---

#### `unlock`

Unlocks a previously generated seed (e.g., after a race concludes).

**Request `data`**
A string formatting the URL of the seed to unlock. (e.g., `"https://maprando.com/seed/tc2pHBSZc/"`)

**Response**
Returns a `"success"` status with no `data`.

**Example**

```json
// Request
{
  "action": "unlock",
  "source": "racetime",
  "data": "https://maprando.com/seed/tc2pHBSZc/"
}

// Response
{
  "status": "success"
}
```
