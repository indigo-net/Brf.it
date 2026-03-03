---
layout: default
title: Getting Started
nav_order: 2
---

# Getting Started

## Installation

### macOS (Homebrew)

```bash
brew install indigo-net/tap/brfit
```

### Linux / macOS (Script)

```bash
curl -fsSL https://raw.githubusercontent.com/indigo-net/Brf.it/main/install.sh | sh
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/indigo-net/Brf.it/main/install.ps1 | iex
```

### From Source

```bash
git clone https://github.com/indigo-net/Brf.it.git
cd Brf.it
go build -o brfit ./cmd/brfit
```

## Verify Installation

```bash
brfit --version
```

## First Run

```bash
# Analyze current directory
brfit .

# Output in Markdown format
brfit . -f md

# Save to file
brfit . -o briefing.xml
```

## Using with AI Assistants

The most common use case is copying the output to your clipboard and pasting it into an AI assistant:

**macOS**
{: .code-label }

```bash
brfit . | pbcopy
```

**Linux (with xclip)**
{: .code-label }

```bash
brfit . | xclip -selection clipboard
```

**Windows**
{: .code-label }

```bash
brfit . | clip
```

## Next Steps

- [CLI Reference](cli-reference) — Full command-line options
- [Language Guides](languages/) — Language-specific documentation
