---
sidebar_position: 10
title: Scala
---

# Scala Support

Brf.it provides full support for Scala with Tree-sitter based parsing.

## Supported Features

- Class definitions
- Case classes
- Object definitions
- Trait definitions
- Function definitions
- Type aliases
- Package declarations
- Import statements (with `--include-imports`)
- Scaladoc comments

## Example

### Input

```scala
package com.example.users

import scala.concurrent.Future

/** Represents a system user */
case class User(
  id: Long,
  name: String,
  email: String
)

/** Trait for user repository implementations */
trait UserRepository {
  /** Find a user by ID
   *  @param id the user ID
   *  @return the user if found
   */
  def findById(id: Long): Future[Option[User]]

  /** Save a user
   *  @param user the user to save
   *  @return the saved user
   */
  def save(user: User): Future[User]
}

/** Service for user operations */
class UserService(repository: UserRepository) {
  def getUser(id: Long): Future[Option[User]] = {
    repository.findById(id)
  }

  def createUser(name: String, email: String): Future[User] = {
    val user = User(id = 0L, name = name, email = email)
    repository.save(user)
  }
}
```

### Output (Brf.it)

```scala
// com/example/users/User.scala
package com.example.users
import scala.concurrent.Future

case class User(id: Long, name: String, email: String)

trait UserRepository {
  def findById(id: Long): Future[Option[User]]
  def save(user: User): Future[User]
}

class UserService(repository: UserRepository) {
  def getUser(id: Long): Future[Option[User]]
  def createUser(name: String, email: String): Future[User]
}
```

## Extensions

Files with `.scala` and `.sc` extensions are processed.
