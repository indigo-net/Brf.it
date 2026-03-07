---
sidebar_position: 4
title: Java
---

# Java Support

Brf.it provides full support for Java with Tree-sitter based parsing.

## Supported Features

- Class definitions
- Interface definitions
- Method declarations
- Annotations
- Generics
- Javadoc comments
- Import statements (with `--include-imports`)

## Example

### Input

```java
package com.example.users;

import java.util.List;
import java.util.Optional;

/**
 * Represents a system user.
 */
public class User {
    private Long id;
    private String name;
    private String email;

    public User(Long id, String name, String email) {
        this.id = id;
        this.name = name;
        this.email = email;
    }

    public Long getId() { return id; }
    public String getName() { return name; }
    public String getEmail() { return email; }
}

/**
 * Repository interface for user data access.
 */
public interface UserRepository {
    Optional<User> findById(Long id);
    List<User> findAll();
    User save(User user);
    void deleteById(Long id);
}
```

### Output (Brf.it)

```java
// com/example/users/User.java
package com.example.users;
import java.util.List; import java.util.Optional;

public class User {
    public User(Long id, String name, String email)
    public Long getId()
    public String getName()
    public String getEmail()
}

public interface UserRepository {
    Optional<User> findById(Long id)
    List<User> findAll()
    User save(User user)
    void deleteById(Long id)
}
```

## Extension

Files with `.java` extension are processed.
