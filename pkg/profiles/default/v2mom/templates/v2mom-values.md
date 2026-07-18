# Values

<!--
V2MOM Values Template
Values are ranked beliefs and principles that guide trade-off decisions.
When two priorities conflict, the higher-ranked value wins.
-->

## Overview

**V2MOM Level**: {{ .Level }}
**Fiscal Period**: {{ .FiscalPeriod }}
**Owner**: {{ .Owner }}

## Ranked Values

<!--
List values in priority order (1 = highest priority).
Each value should include:
- A clear, action-oriented statement
- What it means in practice
- Example trade-offs it guides
-->

### 1. {{ .Value1Name }}

**Statement**: {{ .Value1Statement }}

**In Practice**:
<!-- How does this value show up in daily decisions? -->

**Trade-off Example**:
<!-- When this conflicts with lower values, what do we choose? -->
> "We will {{ .Value1Choice }} even if it means {{ .Value1Tradeoff }}"

---

### 2. {{ .Value2Name }}

**Statement**: {{ .Value2Statement }}

**In Practice**:

**Trade-off Example**:
> "We will {{ .Value2Choice }} even if it means {{ .Value2Tradeoff }}"

---

### 3. {{ .Value3Name }}

**Statement**: {{ .Value3Statement }}

**In Practice**:

**Trade-off Example**:
> "We will {{ .Value3Choice }} even if it means {{ .Value3Tradeoff }}"

---

### 4. {{ .Value4Name }}

**Statement**: {{ .Value4Statement }}

**In Practice**:

**Trade-off Example**:
> "We will {{ .Value4Choice }} even if it means {{ .Value4Tradeoff }}"

---

### 5. {{ .Value5Name }}

**Statement**: {{ .Value5Statement }}

**In Practice**:

**Trade-off Example**:
> "We will {{ .Value5Choice }} even if it means {{ .Value5Tradeoff }}"

## Values Alignment

### Organizational Values
<!-- How do these values align with broader organizational values? -->

### Parent V2MOM Values
<!-- For non-company levels: How do these values support parent values? -->

## Decision Framework

When facing a decision where values conflict:

1. Identify which values are in tension
2. Refer to the ranking above
3. Choose the action that honors the higher-ranked value
4. Document the trade-off decision

---

## Metadata

```yaml
v2mom:
  level: {{ .Level }}
  fiscal_period: {{ .FiscalPeriod }}
  values:
    - rank: 1
      name: "{{ .Value1Name }}"
    - rank: 2
      name: "{{ .Value2Name }}"
    - rank: 3
      name: "{{ .Value3Name }}"
    - rank: 4
      name: "{{ .Value4Name }}"
    - rank: 5
      name: "{{ .Value5Name }}"
```
