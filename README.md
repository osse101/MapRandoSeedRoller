# Seed Roller

A bot and API for rolling randomizer seeds, with support for presets, overrides, and race mode.  The bot holds the spoiler token and can unlock the spoiler log upon request.

---

## Table of Contents

- [User Commands](#user-commands)
  - [Presets](#presets)
  - [Value Overrides](#value-overrides)
  - [Override Values Reference](#override-values-reference)
- [API Reference](#api-reference)
  - [POST /roll](#post-roll)
  - [POST /unlock](#post-unlock)

---

## User Commands

```
!roll <args>
```

| Example | Result |
|---|---|
| `!roll` | Rolls a seed using the current season's preset |
| `!roll s5` | Rolls a seed using the `s5` preset |
| `!roll s5 RDS:MORPH` | Rolls a seed using `s5` with the specified field overrides |

---

### Presets

| Preset | Difficulty | Logic | Notes |
|---|---|---|---|
| `s2` | Hard | Tricky | |
| `s3` | Hard | Tricky | |
| `s3a` | Hard | Tricky | Save the animals |
| `s4` | Hard | Tricky | |
| `s5` | Hard | Tricky | Current season default |
| `default` | Normal | Basic | This is Map Rando's introductory preset |
| `expert` | Expert | Challenge | Based on s5 |
| `mentor` | Medium | Basic | |
| `objectives` | Hard | Tricky | Same as s4; objectives set via args |
| `suits` | Hard | Tricky | Starting Gravity + Varia; rest is s4 |
| `g91` | Hard | Tricky | Starting Gravity + 9 E-Tanks + 1 Reserve; rest is s4 |
| `draft` | Hard | Tricky | Same as s4; starting items set via args |
| `metroids` | Hard | Tricky | Metroid objectives, no Mother Brain 2; rest is s4 |
| `noobjectives` | Hard | Tricky | No objectives; rest is s4 |
| `vmode` | Hard | Tricky | Starting Varia, Grapple, HiJump, Ice, Springball; rest is s4 |

---

### Value Overrides

Overrides are passed as flags after the preset name. Letter case controls the override behaviour:

| Case | Meaning |
|---|---|
| `UPPERCASE` | Enable |
| `lowercase` | Disable |
| `miXEDcasE` | Maybe |

**Available flags:**

| Flag | Description |
|---|---|
| `R` | Race mode |
| `D` | Use dev site |
| `X:<value>` | Escape timer multiplier |
| `S:<items>` | Starting items |
| `O:<objectives>` | Objectives |
| `L:<layout>` | Map layout |

---

### Override Values Reference

**Objectives**

```
kraid, phan, dray, ridley, spore, croc, bot, gt, bt, bowling, acid, pit, babyk, plasma, metal, m1, m2, m3, m4
```

**Starting Items**

```
missile, etank, rtank, super, pb, charge, ice, wave, spazer, plasma, xray, morph, bomb, grapple, hjb, speed, spring, space, screw, varia, gravity
```

**Map Layouts**

```
vanilla, small, standard, wild
```

---

## API Reference

### POST /roll

Rolls a new seed and returns its URL.

**Request**

| Field | Type | Description |
|---|---|---|
| `args` | string | Seed rolling parameters such as Preset and field override arguments |
| `source` | string | The source identifier for the request |
| `source_info` | string | (TBD) Data adding context to the request such as current title and event name |

**Response**

| Field | Type | Description |
|---|---|---|
| `seed_url` | string | The URL of the generated seed |
| `seed_hash` | string | The ingame hash code for the seed |
| `info` | string | A title for this 

**Example**

```json
// Request
{
  "args": "s5 RDS:MORPH",
  "source": "racetime",
  "source_info": "s5 preset"
}

// Response
{
  "seed_url": "https://maprando.com/seed/tc2pHBSZc/",
  "seed_hash": "YARD YARD YARD YARD",
  "info": "s5 preset | https://maprando.com/seed/tc2pHBSZc/ | YARD YARD YARD YARD",
  "message": "Your seed: https://maprando.com/seed/tc2pHBSZc/"
}
```

---

### POST /unlock

Unlocks a previously generated seed (e.g. after a race concludes).

**Request**

| Field | Type | Description |
|---|---|---|
| `seed_url` | string | The URL of the seed to unlock |

**Response**

Returns `200 OK` with an empty body on success.
