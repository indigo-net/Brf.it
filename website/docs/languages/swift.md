---
sidebar_position: 9
title: Swift
---

# Swift Support

Brf.it provides full support for Swift with Tree-sitter based parsing.

## Supported Features

- Class definitions
- Struct definitions
- Protocol definitions
- Enum definitions
- Function definitions
- Property declarations
- Extensions
- Doc comments (///)
- Import statements (with `--include-imports`)

## Example

### Input

```swift
import Foundation

/// Represents a system user
struct User: Codable {
    let id: Int
    let name: String
    let email: String

    /// Validates the user's email format
    func hasValidEmail() -> Bool {
        email.contains("@")
    }
}

/// Protocol for user persistence
protocol UserRepository {
    func findById(_ id: Int) async throws -> User?
    func save(_ user: User) async throws -> User
}

/// Service for user operations
class UserService {
    private let repository: UserRepository

    init(repository: UserRepository) {
        self.repository = repository
    }

    /// Fetches a user by ID
    func getUser(id: Int) async throws -> User? {
        try await repository.findById(id)
    }
}
```

### Output (Brf.it)

```swift
// User.swift
import Foundation

struct User: Codable { let id: Int; let name: String; let email: String
    func hasValidEmail() -> Bool
}

protocol UserRepository {
    func findById(_ id: Int) async throws -> User?
    func save(_ user: User) async throws -> User
}

class UserService {
    init(repository: UserRepository)
    func getUser(id: Int) async throws -> User?
}
```

## Extension

Files with `.swift` extension are processed.
