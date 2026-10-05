# ⚡ RuffSQL (`ruffsql`)

<p align="center">
  <strong>An Extremely Fast SQL Linter and Formatter.</strong><br>
  <em>120x+ faster than Python SQLFluff. Single native binary. Zero Python dependencies.</em>
</p>

<p align="center">
  <a href="#benchmarks">⚡ Benchmarks</a> •
  <a href="#why-ruffsql">💡 Why RuffSQL?</a> •
  <a href="#key-features">✨ Features</a> •
  <a href="#quick-start">🚀 Quick Start</a> •
  <a href="#supported-dialects">🗄️ Dialects</a> •
  <a href="#enterprise--licensing">💼 Enterprise & Licensing</a>
</p>

---

## 💡 Why RuffSQL?

In modern data stacks (dbt, BigQuery, Snowflake, PostgreSQL), linting and formatting thousands of SQL models in CI/CD pipelines using traditional Python linters takes **15 to 30 minutes**.

**RuffSQL** does for SQL what **Ruff did for Python**:
- Replaces heavy, multi-gigabyte Python runtimes with a single, standalone native binary.
- Lints and fixes queries in **milliseconds** instead of minutes.
- Drops directly into existing workflows with `.sqlfluff` and `-- noqa` compatibility.

---

## 📊 Real-World Battle Benchmarks

Benchmarked against **Python SQLFluff 4.4.0** using **171 real PostgreSQL dialect test fixtures** from the official SQLFluff repository on identical hardware (Intel Core i5, 8 cores):

| Metric | Python SQLFluff 4.4.0 | RuffSQL (Go Native) | Performance Gain |
| :--- | :---: | :---: | :---: |
| **Cold Start (1 File CLI)** | 902.6 ms | **4.2 ms** | **🚀 214x Faster** |
| **Full Suite (171 Real SQL Files, 8 Cores)** | 9.340 s | **0.077 s (77 ms)** | **⚡ 121.3x Faster** |
| **Average Query Time (in batch)** | 54.62 ms / file | **0.45 ms / file** | **⚡ 121x Faster** |
| **Memory Footprint (Peak RSS)** | 45.4 MB | **5.4 MB** | **📉 8.4x Less Memory** |
| **In-Process Parsing & Linting (warm)** | 22.2 ms / op | **0.063 ms (63 µs)** | **🔥 350x Faster** |

> 💡 **CI/CD Impact:** A dbt repository with 1,500 models that previously stalled CI builds for **10–15 minutes** now completes in **under 1.5 seconds**.

---

## ✨ Key Features

- **🚀 120x–200x Faster Execution**: Built in Go with lock-free goroutine worker pools.
- **📦 Single Self-Contained Binary**: No Python, no `pip`, no virtual environments, no dependency hell in Docker containers.
- **🛠️ Production Auto-Fixer (`ruffsql fix`)**: Safe, iterative AST-based auto-fix engine with `--write` support.
- **🧩 Native dbt / Jinja2 Support**: Transparently parses and lints files with `{{ config(...) }}`, `{{ ref(...) }}`, and `{% if %}` macros without breaking source coordinates.
- **⚙️ Config Compatible**: Reads `.sqlfluff`, `.sqlfluffignore`, and respects inline `-- noqa` comment directives.
- **🔄 CI/CD Ready**: Native multi-format output (`--format human`, `--format json`, `--format github-annotation`).
- **🗄️ Procedural SQL Support**: Full support for `CREATE PROCEDURE`, `CREATE FUNCTION`, `CREATE TRIGGER`, `CALL / EXEC`, and `BEGIN ... END` across dialects.

---

## 🚀 Quick Start

### 1. Installation

Download the single precompiled binary or build directly with Go:

```bash
git clone https://github.com/your-org/ruffsql.git
cd ruffsql
go build -o /usr/local/bin/ruffsql ./cmd/sqlfluff-go
```

### 2. Linting SQL Files

```bash
# Lint a single file or directory
ruffsql lint models/ --dialect postgres

# Lint with multiple workers and JSON output (ideal for CI)
ruffsql lint --workers 8 --format json models/

# Native GitHub Actions annotations
ruffsql lint --format github-annotation models/
```

### 3. Auto-Fixing Violations

```bash
# Preview fixes
ruffsql fix models/orders.sql

# Apply fixes directly to disk
ruffsql fix --write models/
```

### 4. Viewing AST Parse Trees

```bash
ruffsql parse query.sql --dialect snowflake
```

---

## 🗄️ Supported Dialects

RuffSQL supports **all 28 SQL dialects**, with deep syntax engines for tier-1 warehouse engines:

- **PostgreSQL**: Type casts (`::`), JSON/JSONB operators (`->`, `->>`), `ON CONFLICT DO UPDATE`.
- **Snowflake**: `QUALIFY`, `OVER (PARTITION BY ... ORDER BY ...)`, `ILIKE`.
- **BigQuery**: Backticked identifiers (`` `project.dataset.table` ``), `SAFE_CAST(...)`, `QUALIFY`.
- **MySQL / MariaDB**: Column list insertions, `ON DUPLICATE KEY UPDATE`.
- **ANSI & Enterprise Dialects**: T-SQL, ClickHouse, DuckDB, Trino, Oracle, SQLite, Athena, Databricks, Redshift, SparkSQL, and more.

---

## 🛡️ Rule Coverage

RuffSQL includes full engines for the most critical layout, capitalisation, and convention rules, with full registry coverage of all 79 rules:

- **Capitalisation (`CP`)**: `CP01` (Keywords), `CP02` (Identifiers), `CP03` (Functions), `CP04` (Literals `NULL/TRUE/FALSE`), `CP05` (Types).
- **Layout & Spacing (`LT`)**: `LT01` (Trailing spaces), `LT04` (Leading commas), `LT05` (Max line length), `LT06` (Function bracket spacing), `LT09` (Select targets newline), `LT12` (EOF newline), `LT14` (WHERE keyword newline).
- **Aliasing (`AL`)**: `AL01` (Table `AS`), `AL02` (Column `AS`), `AL09` (Self-aliasing).
- **Convention (`CV`)**: `CV01` (`<>` vs `!=`), `CV02` (`COALESCE` vs `NVL/IFNULL`), `CV03` (Trailing commas), `CV04` (`COUNT(*)`), `CV05` (`IS NULL`), `CV06` (Terminator semicolons), `CV08` (Disallow `RIGHT JOIN`).
- **Structure (`ST`)**: `ST01` (Redundant `ELSE NULL`), `ST07` (Disallow `USING`), `ST12` (Consecutive semicolons).
- **Ambiguous & Quotes (`AM`/`TQ`)**: `AM01` (Bare `DISTINCT`), `AM02` (Explicit `UNION`), `TQ01` (Unnecessary backticks/quotes).

---

## 💼 Enterprise & Licensing

RuffSQL is distributed under a **Dual-Licensing Model**:

1. **Community Edition (AGPLv3 / Open Core)**:
   - 100% free for individual developers, local workstations, and open-source projects.
2. **Commercial Enterprise License**:
   - For enterprise teams and proprietary SaaS platforms that require private deployment without AGPL copyleft obligations.
   - Includes priority support, custom company-wide security rules (PII leakage audit in queries), and central enforcement gateways.

*Contact: enterprise@ruffsql.dev*
