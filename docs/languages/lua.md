---
layout: default
title: Lua
parent: Language Guides
nav_order: 12
---

# Lua Support

[English](lua.md) | [한국어](../ko/languages/lua.md) | [日本語](../ja/languages/lua.md) | [हिन्दी](../hi/languages/lua.md) | [Deutsch](../de/languages/lua.md)

## Supported Extensions

- `.lua`

## Extraction Targets

| Element | Kind | Example |
|---------|------|---------|
| Global Function | `function` | `function globalFunc(a, b)` |
| Local Function | `local_function` | `local function helper()` |
| Module Function | `module_function` | `function M.create(name)` |
| Method | `method` | `function M:init(config)` |
| Table Assignment | `variable` | `local M = {}` |
| Function Assignment | `variable` | `local foo = function() end` |
| LuaDoc Comment | `doc` | `--- Description` |
| Line Comment | `doc` | `-- Description` |
| require() | `import` | `local json = require("json")` |

## Example

### Input

```lua
local M = {}

--- Greets a person by name.
-- @param name string The person's name
function M.greet(name)
    print("Hello, " .. name)
end

function M:init(config)
    self.config = config
end

local function helper()
    return 42
end

function globalFunc(a, b)
    return a + b
end

local json = require("json")
```

### Output (XML)

```xml
<file path="example.lua" language="lua">
  <variable kind="variable" line="1">
    <name>M</name>
    <text>local M = {}</text>
  </variable>
  <function kind="module_function" line="4">
    <name>greet</name>
    <text>function M.greet(name)</text>
  </function>
  <function kind="method" line="9">
    <name>init</name>
    <text>function M:init(config)</text>
  </function>
  <function kind="local_function" line="13">
    <name>helper</name>
    <text>local function helper()</text>
  </function>
  <function kind="function" line="17">
    <name>globalFunc</name>
    <text>function globalFunc(a, b)</text>
  </function>
</file>
```

## Notes

### Visibility

- All declarations are extracted (Lua has no access modifiers)
- `local` functions/variables are included alongside global ones

### Function Kinds

- `function`: Global function declarations (`function foo()`)
- `local_function`: Local function declarations (`local function foo()`)
- `module_function`: Dot-indexed functions (`function M.foo()`)
- `method`: Colon-indexed methods (`function M:foo()`)

### Body Removal

When `--include-body` flag is not used:

- Functions/Methods: body removed after closing parenthesis `)` of parameter list
- Table assignments (`local M = {}`): preserved as-is
- Function assignments (`local foo = function() end`): preserved as-is

### Import Extraction

- `require()` calls are extracted with `--include-imports` flag
- Format: `local json = require("json")` (full statement preserved)

### Doc Comments

- LuaDoc comments (`---`) and regular line comments (`--`) are both extracted
- Block comments (`--[[ ... ]]`) are also supported
