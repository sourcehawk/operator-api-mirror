# Operator API Mirror

This is a small utility that extracts the **public API types** from upstream Kubernetes operators 
(e.g., OpenTelemetry Operator, ECK) and republishes them as **version-pinned, standalone Go modules**.

This allows downstream consumers to depend on operator APIs **in a reproducible, conflict-free way**, without 
inheriting the upstream operator's dependency tree or Go module drift.

---

## ✨ Why this project exists

Many Kubernetes operators ship their CRD Go types inside complex monolithic repositories. Importing those API types 
directly often pulls in:

* Huge dependency graphs
* Unwanted Kubernetes version constraints
* Internal packages that are not meant to be consumed
* Frequent breaking changes as upstream moves fast
* Non-reproducible builds when upstream bumps its own dependencies

**`operator-api-mirror` fixes this.**

It produces stable, minimal, isolated Go modules containing *only*:

* the selected API directories (e.g. `apis/*`, `pkg/apis/elasticsearch/*`)
* the internal packages required to build them
* rewritten imports so the API code is self-contained
* a minimal `go.mod` with just the required dependencies
* optional overrides for Kubernetes versions or other modules

This guarantees reproducible builds while allowing consumers to use *exactly the API types they need*.

---

## 🧩 How it works

For each operator defined in `operators.yaml`, the tool:

1. **Clones the upstream repository** at the specified version.
2. **Copies the requested API directories** and filters out all non-Go files and test files.
3. **Parses API imports** to discover and copy required internal packages.
4. **Rewrites import paths** to point to the mirror module.
5. **Generates a fresh `go.mod`** based on upstream requirements.
6. **Applies dependency overrides** from `operators.yaml` using `replace`.
7. **Runs `go mod tidy`** to prune unused dependencies.
8. Writes the final result to:

```
mirrors/<slug>/<version>/
```

Each mirrored API version is a fully standalone Go module.

---

## 📦 Example: operators.yaml

```yaml
operators:
  - slug: otel-operator
    repo: github.com/open-telemetry/opentelemetry-operator
    currentVersion: v0.138.0
    goModPath: go.mod
    apiPaths:
      - "apis/*"
    overwriteDependencies:
      - name: k8s.io/api
        version: v0.32.10
      - name: k8s.io/apimachinery
        version: v0.32.10
      - name: sigs.k8s.io/controller-runtime
        version: v0.20.4

  - slug: eck-operator
    repo: github.com/elastic/cloud-on-k8s
    currentVersion: v3.2.0
    goModPath: go.mod
    apiPaths:
      - "pkg/apis/elasticsearch/*"
    overwriteDependencies:
      - name: k8s.io/api
        version: v0.32.10
      - name: k8s.io/apimachinery
        version: v0.32.10
```

---

## 🚀 Creating your own mirrors

This repository contains **ready-made mirrors** that do not overwrite or pin dependencies differently from the upstream operators.

If you need **custom dependency overrides** (for example, to keep your Kubernetes libraries at specific versions different from what the upstream operator uses), you can create your **own mirror repository** using this Go tool as a library and CLI.

### Build / install the tool

From *another* repo (your own mirror repo) you can install the CLI globally:

```bash
go install github.com/sourcehawk/operator-api-mirror/cmd/api-mirrorer@latest
```

This will put an `api-mirrorer` binary in your `$GOBIN` (usually `~/go/bin`).

Or, if you’re working directly inside a clone of this repo:

```bash
go build -o api-mirrorer ./cmd/api-mirrorer
```

### Run the mirror process

At minimum you must tell the tool what your **root module name** is – that’s the prefix used when generating modules like:

```go
module github.com/my-org/operator-api-mirror/otel-operator/v0.138.0
```

Example:

```bash
./api-mirrorer \
  -rootModuleName github.com/my-org/operator-api-mirror
```

The mirrored modules will appear under:

```text
mirrors/<operator>/<version>/
```

For example:

```text
mirrors/
└── otel-operator/
    └── v0.138.0/
        ├── apis/
        ├── internal/
        ├── pkg/
        └── go.mod
```

You can then import from your own repo:

```go
import "github.com/my-org/operator-api-mirror/otel-operator/v0.138.0/apis/v1beta1"
```

### CLI flags

The tool accepts several optional flags that let you customize where configuration is read from and where mirror
modules are written.

| Flag              | Default           | Description                                                                                | Example                                 |
|-------------------|-------------------|--------------------------------------------------------------------------------------------|-----------------------------------------|
| `-config`         | `operators.yaml`  | Path to the operators configuration file defining which operators and API paths to mirror. | `operators.yaml`                        |
| `-mirrorsPath`    | `./mirrors`       | Directory where all mirrored operator modules will be generated.                           | `./generated/mirrors`                   |
| `-rootModuleName` | *none* (required) | Root Go module path used when creating `module <path>` in each mirrored module.            | `github.com/my-org/operator-api-mirror` |

Typical full invocation:

```bash
./api-mirrorer \
  -config ./operators.yaml \
  -mirrorsPath ./mirrors \
  -rootModuleName github.com/my-org/operator-api-mirror
```

---

## ⚙️ How dependency overrides work

`overwriteDependencies` applies **replace-only** overrides in `go.mod`:

```go
replace k8s.io/apimachinery => k8s.io/apimachinery v0.32.10
replace k8s.io/api => k8s.io/api v0.32.10
```

These do not modify upstream requirements; they simply redirect them.

This ensures:

* consistent Kubernetes versions across all mirrored modules
* maximum compatibility for downstream controllers
* no accidental dependency bumps caused by upstream operators

---

## 🛠 Development

Update `operators.yaml` to add new operators or bump versions.

You can also automate version bumps using Renovate or GitHub Actions.
A typical workflow:

* Renovate detects new upstream tags
* It bumps `currentVersion` in `operators.yaml`
* GitHub Actions runs the mirror job
* The mirrored modules are committed and tagged automatically

If you'd like an example GitHub workflow, just ask!

---

## ❓ FAQ

### Why not just import the upstream operator?

Because:

* Operators have *huge transitive dependency trees*
* Many require pinned Kubernetes versions that break your project
* They often pull in unrelated internal code
* They are not intended to be library-friendly

Using a mirror gives you:

* minimal deps
* reproducible API versions
* clean import paths
* stable builds

### Why copy internal packages?

Some operator API code references internal utilities (e.g. `ptr`, `hash`, `version`).
We copy only what is needed and rewrite the import paths so everything is self-contained.

### Does this violate licensing?

All upstream APIs are open-source (Apache 2.0); the mirror preserves licenses verbatim.

---

## 🤝 Contributing

Contributions welcome!

* Support for more operators
  * Create a fork
  * Create a feature branch on your fork 
  * Add configuration to the operators.yaml file
  * Create a PR on this repository

---

## 📝 License

Apache 2.0 — same as upstream operator sources.
