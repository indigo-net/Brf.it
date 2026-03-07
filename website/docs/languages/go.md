---
sidebar_position: 1
title: Go
---

# Go Support

Brf.it provides full support for Go with Tree-sitter based parsing.

## Supported Features

- Function declarations
- Method declarations (with receiver types)
- Interface definitions
- Struct definitions
- Type aliases
- Import statements (with `--include-imports`)
- Package documentation

## Example

### Input

```go
package userservice

import (
    "context"
    "errors"
)

// User represents a system user
type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// UserService handles user operations
type UserService interface {
    GetUser(ctx context.Context, id string) (*User, error)
    CreateUser(ctx context.Context, user *User) error
}

// NotFoundError is returned when a user is not found
var NotFoundError = errors.New("user not found")

// GetUser retrieves a user by ID
func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
    // implementation...
    return nil, NotFoundError
}
```

### Output (Brf.it)

```go
// userservice/user.go
package userservice

import ( "context"; "errors" )

type User struct { ID string; Name string; Email string }
type UserService interface { GetUser(ctx context.Context, id string) (*User, error); CreateUser(ctx context.Context, user *User) error }
var NotFoundError = errors.New("user not found")
func (s *Service) GetUser(ctx context.Context, id string) (*User, error)
```

## Extension

Files with `.go` extension are processed.
