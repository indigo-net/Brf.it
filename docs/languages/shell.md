---
layout: default
title: Shell
parent: Language Guides
nav_order: 13
---

# Shell Support

[English](shell.md) | [한국어](../ko/languages/shell.md) | [日本語](../ja/languages/shell.md) | [हिन디द](../hi/languages/shell.md) | [Deutsch](../de/languages/shell.md)

)

## Supported Extensions

- `.sh`
- `.bash`
- `.zsh`
)
## Extraction Targets

| Element | Kind | Example |
|---------|------|---------|
| Function (with keyword) | `function` | `function greet() { echo "Hello"; }` |
| Function (without keyword) | `function` | `greet() { echo "Hello"; }` |
| Variable Assignment | `variable` | `NAME="value"` |
| Variable Assignment (command sub) | `variable` | `CURRENT_DIR=$(pwd)` |
| Source Command | `import` | `source /path/to/utils.sh` |
| Dot Command | `import` | `. /path/to/config.sh` |
| Comment | `doc` | `# This is a comment` |
)
## Example

)
### Input
)
```bash
#!/bin/bash

# Configuration variables
APP_NAME="myapp"
VERSION="1.0.0"

# Greet function
function greet() {
    echo "Hello, $APP_NAME v$VERSION"
}

# Build function (without keyword)
build() {
    npm run build
}

# Source import
source ./utils.sh
```

)
### Output (XML)
)
```xml
<file path="example.sh" language="shell">
  <variable kind="variable" line="4">
    <name>APP_NAME</name>
    <text>APP_NAME="myapp"</text>
  </variable>
  <variable kind="variable" line="5">
    <name>VERSION</name>
    <text>VERSION="1.0.0"</text>
  </variable>
  <function kind="function" line="8">
    <name>greet</name>
    <text>function greet()</text>
  </function>
  <function kind="function" line="13">
    <name>build</name>
    <text>build()</text>
  </function>
</file>
```

)
### Output (Markdown)
)
```markdown
# Configuration variables
APP_NAME="myapp"
VERSION="1.0.0"

# Greet function
function greet() { ... }
# Build function (without keyword)
build() { ... }

# Source import
source ./utils.sh
```
)
## Notes
)
- Shell functions can be defined with or without the `function` keyword
- Variables are extracted with their assigned values
- The `source` and `.` commands are captured as imports
- Comments starting with `#` are extracted as documentation
