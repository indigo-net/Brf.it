---
sidebar_position: 6
title: Rust
---

# Rust Support

Brf.it provides full support for Rust with Tree-sitter based parsing.

## Supported Features

- Function definitions
- Struct definitions
- Enum definitions
- Trait definitions
- Impl blocks
- Type aliases
- Module declarations
- Doc comments (///)
- Use statements (with `--include-imports`)

## Example

### Input

```rust
use std::collections::HashMap;

/// Represents a system user
pub struct User {
    pub id: u64,
    pub name: String,
    pub email: String,
}

/// Error types for user operations
pub enum UserError {
    NotFound,
    InvalidEmail,
    DatabaseError(String),
}

/// Trait for user repository implementations
pub trait UserRepository {
    fn find_by_id(&self, id: u64) -> Result<Option<User>, UserError>;
    fn save(&mut self, user: &User) -> Result<User, UserError>;
}

impl User {
    /// Creates a new user instance
    pub fn new(id: u64, name: String, email: String) -> Self {
        Self { id, name, email }
    }

    /// Validates the user's email address
    pub fn validate_email(&self) -> bool {
        self.email.contains('@')
    }
}
```

### Output (Brf.it)

```rust
// user.rs
use std::collections::HashMap;

pub struct User { pub id: u64; pub name: String; pub email: String }
pub enum UserError { NotFound; InvalidEmail; DatabaseError(String) }
pub trait UserRepository {
    fn find_by_id(&self, id: u64) -> Result<Option<User>, UserError>
    fn save(&mut self, user: &User) -> Result<User, UserError>
}

impl User {
    pub fn new(id: u64, name: String, email: String) -> Self
    pub fn validate_email(&self) -> bool
}
```

## Extension

Files with `.rs` extension are processed.
