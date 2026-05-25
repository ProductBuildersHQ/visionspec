# profiles

Manage configuration profiles that bundle spec requirements, templates, and rubrics.

## Synopsis

```bash
multispec profiles <subcommand>
```

## Description

Configuration profiles allow you to customize which specs are required and how they're evaluated. MultiSpec includes built-in profiles for different product stages:

| Profile | Description | Required Specs |
|---------|-------------|----------------|
| `0-1` | Minimal for idea validation | hypothesis |
| `startup` | Lightweight for pre-PMF | prd |
| `growth` | Metrics-driven for 1-N scaling | prd, uxd, faq |
| `enterprise` | Comprehensive for post-PMF | prd, mrd, uxd, trd, press, faq |

## Subcommands

### list

List all available profiles.

```bash
multispec profiles list
```

**Output:**

```
Available profiles:

  0-1          Minimal configuration for idea validation [default]
  startup      Lightweight configuration for pre-PMF startups [default]
  growth       Metrics-driven configuration for 1-N scaling [default]
  enterprise   Comprehensive configuration for post-PMF enterprises [default]

Use with: multispec init <project> --profile <name>
```

### show

Show detailed information about a profile.

```bash
multispec profiles show <profile-name>
```

**Example:**

```bash
multispec profiles show enterprise
```

**Output:**

```
Profile: enterprise
Description: Comprehensive configuration for post-PMF enterprises. Full specs with security and compliance.

Required specs:
  - prd (source)
  - mrd (source)
  - uxd (source)
  - press (gtm)
  - faq (gtm)
  - trd (technical)

Custom templates:
  - prd
  - mrd
  - uxd
  - trd
  - press
  - faq

Custom rubrics:
  - prd
  - mrd
  - uxd
  - trd
  - press
  - faq
```

### export

Export a profile to a directory for customization.

```bash
multispec profiles export <profile-name> <output-dir>
```

This creates a complete profile directory with:

- `profile.yaml` - Configuration file
- `templates/` - Template files (`.md`)
- `rubrics/` - Rubric files (`.rubric.yaml`)

**Example:**

```bash
multispec profiles export enterprise ./my-profile
```

**Output:**

```
Created ./my-profile/profile.yaml
Created ./my-profile/templates/prd.md
Created ./my-profile/templates/mrd.md
Created ./my-profile/templates/uxd.md
Created ./my-profile/templates/trd.md
Created ./my-profile/templates/press.md
Created ./my-profile/templates/faq.md
Created ./my-profile/rubrics/prd.rubric.yaml
Created ./my-profile/rubrics/mrd.rubric.yaml
Created ./my-profile/rubrics/uxd.rubric.yaml
Created ./my-profile/rubrics/trd.rubric.yaml
Created ./my-profile/rubrics/press.rubric.yaml
Created ./my-profile/rubrics/faq.rubric.yaml

Profile exported to ./my-profile

To use this profile:
  multispec init my-project --profile-dir ./my-profile
```

## Using Profiles

### With init command

```bash
# Use a built-in profile
multispec init my-project --profile startup

# Use a custom profile directory
multispec init my-project --profile-dir ./my-profile
```

### Profile inheritance

Profiles can extend other profiles using the `extends` field:

```yaml
# my-profile/profile.yaml
name: my-enterprise
description: "Custom enterprise profile with additional requirements"
extends: enterprise

spec_config:
  # Add a custom spec type
  security-review:
    required: true
    category: technical
```

## See Also

- [Custom Profiles Guide](../guides/custom-profiles.md) - Complete guide to creating custom profiles
- [init](init.md) - Initialize projects with profiles
