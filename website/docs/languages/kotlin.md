---
sidebar_position: 5
title: Kotlin
---

# Kotlin Support

Brf.it provides full support for Kotlin with Tree-sitter based parsing.

## Supported Features

- Class definitions
- Data classes
- Interface definitions
- Function declarations
- Extension functions
- Properties
- Annotations
- Import statements (with `--include-imports`)
- KDoc comments

## Example

### Input

```kotlin
package com.example.users

data class User(
    val id: Long,
    val name: String,
    val email: String
)

interface UserRepository {
    suspend fun findById(id: Long): User?
    suspend fun findAll(): List<User>
    suspend fun save(user: User): User
}

class UserService(
    private val repository: UserRepository
) {
    suspend fun getUser(id: Long): User? {
        return repository.findById(id)
    }

    suspend fun createUser(name: String, email: String): User {
        val user = User(id = 0, name = name, email = email)
        return repository.save(user)
    }
}

// Extension function
fun User.toDto(): UserDto = UserDto(
    id = this.id.toString(),
    name = this.name,
    email = this.email
)
```

### Output (Brf.it)

```kotlin
// com/example/users/User.kt
package com.example.users

data class User(val id: Long, val name: String, val email: String)

interface UserRepository {
    suspend fun findById(id: Long): User?
    suspend fun findAll(): List<User>
    suspend fun save(user: User): User
}

class UserService(private val repository: UserRepository) {
    suspend fun getUser(id: Long): User?
    suspend fun createUser(name: String, email: String): User
}

fun User.toDto(): UserDto
```

## Extensions

Files with `.kt` and `.kts` extensions are processed.
