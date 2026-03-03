---
layout: default
title: Documentation
nav_order: 1
has_children: true
---

# Documentation

Welcome to the Brf.it documentation. Brf.it is a CLI tool that packages your codebase for AI comprehension by extracting function signatures and documentation.

## Getting Started

New to Brf.it? Start here:

1. [**Installation**](getting-started) — Install brfit on your system
2. [**CLI Reference**](cli-reference) — Complete command-line options
3. [**Language Guides**](languages/) — Language-specific configuration

## What Brf.it Does

Brf.it analyzes your codebase and extracts:

- **Function signatures** — Name, parameters, return types
- **Documentation comments** — JSDoc, GoDoc, docstrings, etc.
- **Import statements** — Module dependencies
- **Directory structure** — Project overview

Output formats:
- **XML** — Structured, machine-readable
- **Markdown** — Human-readable, great for documentation

## Quick Example

```bash
# Analyze current directory
brfit .

# Output in Markdown
brfit . -f md

# Save to file
brfit . -o briefing.xml

# Copy to clipboard for AI
brfit . | pbcopy    # macOS
```
