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

| Element | Kind | Example |
|---------|------|---------|
| Method | `method` | `def greet(name)` |
| Class Method | `class_method` | `def self.create(attrs)` |
| Class | `class` | `class User < ActiveRecord::Base` |
| Module | `module` | `module Authentication` |
| Constant | `variable` | `MAX_RETRIES = 3` |
| YARD Comment | `doc` | `# Description` |
| require | `import` | `require "json"` |
| require_relative | `import` | `require_relative "helpers"` |

## Example

### Input

```ruby
require "json"
require_relative "helpers"

# Represents a user in the system.
class User
  MAX_RETRIES = 3

  # Creates a new user from attributes.
  # @param attrs [Hash] user attributes
  def self.create(attrs)
    new(attrs).save
  end

  # Initializes the user.
  # @param name [String] the user's name
  def initialize(name)
    @name = name
  end

  # Greets another person.
  # @param other [String] the other person's name
  # @return [String] a greeting message
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
  <class kind="class" line="5">
    <name>User</name>
    <text>class User</text>
  </class>
  <variable kind="variable" line="6">
    <name>MAX_RETRIES</name>
    <text>MAX_RETRIES = 3</text>
  </variable>
  <function kind="class_method" line="10">
    <name>create</name>
    <text>def self.create(attrs)</text>
  </function>
  <function kind="method" line="15">
    <name>initialize</name>
    <text>def initialize(name)</text>
  </function>
  <function kind="method" line="21">
    <name>greet</name>
    <text>def greet(other)</text>
  </function>
  <module kind="module" line="27">
    <name>Authentication</name>
    <text>module Authentication</text>
  </module>
  <function kind="method" line="28">
    <name>authenticate</name>
    <text>def authenticate(password)</text>
  </function>
</file>
```

## Notes

### Visibility

- All methods are extracted regardless of visibility (`public`, `protected`, `private`)
- Both instance methods (`def foo`) and class methods (`def self.foo`) are captured

### Method Kinds

- `method`: Instance method declarations (`def foo`)
- `class_method`: Class-level method declarations (`def self.foo`)

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
