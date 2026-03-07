---
sidebar_position: 7
title: Ruby
---

# Ruby Support

Brf.it provides full support for Ruby with Tree-sitter based parsing.

## Supported Features

- Class definitions
- Module definitions
- Method definitions
- Singleton methods
- Include/extend statements
- Attribute accessors (attr_reader, attr_writer, attr_accessor)
- YARD/RDoc comments
- Require statements (with `--include-imports`)

## Example

### Input

```ruby
require 'date'

# Represents a system user
class User
  attr_reader :id, :name, :email
  attr_accessor :status

  # Creates a new user instance
  # @param id [Integer] the user ID
  # @param name [String] the user name
  # @param email [String] the user email
  def initialize(id, name, email)
    @id = id
    @name = name
    @email = email
    @status = :active
  end

  # Validates the user's email format
  # @return [Boolean] true if email is valid
  def valid_email?
    email.include?('@')
  end
end

# Repository for user persistence
module UserRepository
  def find_by_id(id)
    raise NotImplementedError
  end

  def save(user)
    raise NotImplementedError
  end
end
```

### Output (Brf.it)

```ruby
# user.rb
require 'date'

class User
  attr_reader :id, :name, :email
  attr_accessor :status
  def initialize(id, name, email)
  def valid_email?
end

module UserRepository
  def find_by_id(id)
  def save(user)
end
```

## Extensions

Files with `.rb` and `.rake` extensions are processed.
