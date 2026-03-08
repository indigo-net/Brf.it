---
layout: default
title: Bash/Shell
parent: Language Guides
nav_order: 13
---

# Bash/Shell Support

[English](bash.md) | [한국어](../ko/languages/bash.md) | [日本語](../ja/languages/bash.md) | [हिन्दी](../hi/languages/bash.md) | [Deutsch](../de/languages/bash.md)

## Supported Extensions

- `.sh`
- `.bash`

## Extraction Targets

| Element | Kind | Example |
|---------|------|---------|
| Function | `function` | `function greet { ... }` |
| Function | `function` | `greet() { ... }` |
| Variable Assignment | `variable` | `NAME="value"` |
| Declaration | `variable` | `declare VERBOSE` |
| Local Variable | `variable` | `local count=0` |
| Readonly Variable | `variable` | `readonly VERSION="1.0"` |
| Comment | `doc` | `# Description` |
| Source Statement | `import` | `source /path/to/lib.sh` |
| Dot Statement | `import` | `. ./config.sh` |

## Example

### Input

```bash
#!/bin/bash

# Configuration settings
CONFIG_PATH="/etc/myapp"
VERSION="1.0.0"
declare VERBOSE=false

# Deploy application to server
function deploy {
    local app_name="$1"
    echo "Deploying $app_name"
}

# Build the project
build() {
    echo "Building..."
}

source ./utils.sh
. ./config.sh
```

### Output (XML)

```xml
<file path="deploy.sh" language="bash">
  <variable kind="variable" line="4">
    <name>CONFIG_PATH</name>
    <text>CONFIG_PATH="/etc/myapp"</text>
  </variable>
  <variable kind="variable" line="5">
    <name>VERSION</name>
    <text>VERSION="1.0.0"</text>
  </variable>
  <variable kind="variable" line="6">
    <name>VERBOSE</name>
    <text>declare VERBOSE=false</text>
  </variable>
  <function kind="function" line="9">
    <name>deploy</name>
    <text>function deploy</text>
  </function>
  <function kind="function" line="15">
    <name>build</name>
    <text>build()</text>
  </function>
</file>
```

## Notes

### Visibility

- All declarations are extracted (Bash has no access modifiers)
- `local` variables inside functions are extracted if declared at parse time

### Function Syntax

Bash supports two function declaration styles:

- `function name { ... }` - with `function` keyword
- `name() { ... }` - with parentheses

Both are extracted as `function` kind.

### Body Removal

When `--include-body` flag is not used:

- Functions: body removed after opening brace `{`
- Variables: only first line preserved (handles multi-line assignments)

### Import Extraction

- `source` and `.` commands are extracted with `--include-imports` flag
- Supports both quoted and unquoted paths

### Doc Comments

- Shell comments starting with `#` are extracted
- Shebang lines (`#!/bin/bash`) are not treated as comments

### Limitations

- Nested functions are supported
- Here-documents in function bodies are handled correctly
- Complex variable expansions are preserved in signatures
