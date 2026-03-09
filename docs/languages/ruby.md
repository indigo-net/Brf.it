---
layout: default
title: Ruby
parent: Language Guides
nav_order: 15
---

# Ruby Support

[English](ruby.md) | [한국어](../ko/languages/ruby.md) | [日本語](../ja/languages/ruby.md) | [हिन्दी](../hi/languages/ruby.md) | [Deutsch](../de/languages/ruby.md)

## Supported Extensions

- `.rb`

## Extraction Targets

| Element | Kind | XML Tag | Example |
|---------|------|---------|---------|
| Method | `method` | `<function>` | `def greet(name)` |
| Class Method | `method` | `<function>` | `def self.create(attrs)` |
| Class | `class` | `<type>` | `class User < ActiveRecord::Base` |
| Module | `namespace` | `<type>` | `module Authentication` |
| Constant (top-level) | `variable` | `<variable>` | `MAX_RETRIES = 3` |
| Comment | `doc` | | `# Description` |
| require | | `<imports>` | `require "json"` |
| require_relative | | `<imports>` | `require_relative "helpers"` |

## Example

### Input

```ruby
require "json"
require_relative "helpers"

MAX_RETRIES = 3

# Represents a user in the system.
class User
  # Creates a new user from attributes.
  def self.create(attrs)
    new(attrs).save
  end

  # Initializes the user.
  def initialize(name)
    @name = name
  end

  # Greets another person.
  def greet(other)
    "Hello, #{other}! I'm #{@name}."
  end
end

module Authentication
  def authenticate(password)
    password == @secret
  end
end
```

### Output (XML)

```xml
<file path="example.rb" language="ruby">
  <type>class User</type>
  <function>def self.create(attrs)</function>
  <function>def initialize(name)</function>
  <function>def greet(other)</function>
  <variable>MAX_RETRIES = 3</variable>
  <type>module Authentication</type>
  <function>def authenticate(password)</function>
</file>
```

## Notes

### Visibility

- All methods are extracted regardless of visibility (`public`, `protected`, `private`)
- Both instance methods (`def foo`) and class methods (`def self.foo`) are captured

### Method Kinds

- Both instance methods (`def foo`) and class methods (`def self.foo`) use kind `method`

### Body Removal

When `--include-body` flag is not used:

- Methods: body removed after closing parenthesis `)` of parameter list (or after method name if no parameters)
- Classes/Modules: only the declaration line is preserved
- Constants: preserved as-is

### Import Extraction

- `require` and `require_relative` statements are extracted with `--include-imports` flag
- Format: `require "json"` / `require_relative "helpers"` (full statement preserved)

### Doc Comments

- YARD-style comments (`#`) directly above methods/classes are extracted
- Multi-line `#` comments are supported
- `=begin`...`=end` block comments are also recognized
