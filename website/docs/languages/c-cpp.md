---
sidebar_position: 11
title: C/C++
---

# C/C++ Support

Brf.it provides full support for C and C++ with Tree-sitter based parsing.

## Supported Features

- Function declarations
- Class/struct definitions
- Template definitions
- Namespace declarations
- Preprocessor directives (with `--include-imports`)
- Doxygen comments

## Example (C++)

### Input

```cpp
#include <string>
#include <vector>
#include <optional>

namespace users {

/**
 * @brief Represents a system user
 */
struct User {
    long id;
    std::string name;
    std::string email;

    /**
     * @brief Validates the user's email format
     * @return true if email is valid
     */
    bool hasValidEmail() const;
};

/**
 * @brief Repository interface for user data access
 */
class UserRepository {
public:
    virtual ~UserRepository() = default;

    virtual std::optional<User> findById(long id) = 0;
    virtual std::vector<User> findAll() = 0;
    virtual User save(const User& user) = 0;
};

/**
 * @brief Template for cached repository
 */
template<typename T>
class CachedRepository {
public:
    CachedRepository(T& repository) : repo(repository) {}

    std::optional<User> findById(long id);

private:
    T& repo;
};

} // namespace users
```

### Output (Brf.it)

```cpp
// users/user.hpp
#include <string>
#include <vector>
#include <optional>

namespace users {

struct User { long id; std::string name; std::string email;
    bool hasValidEmail() const
}

class UserRepository {
public:
    virtual ~UserRepository() = default
    virtual std::optional<User> findById(long id) = 0
    virtual std::vector<User> findAll() = 0
    virtual User save(const User& user) = 0
}

template<typename T>
class CachedRepository {
public:
    CachedRepository(T& repository)
    std::optional<User> findById(long id)
private:
    T& repo
}

} // namespace users
```

## Extensions

- C: `.c`, `.h`
- C++: `.cpp`, `.cc`, `.cxx`, `.hpp`, `.h`
