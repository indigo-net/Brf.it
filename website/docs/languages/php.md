---
sidebar_position: 8
title: PHP
---

# PHP Support

Brf.it provides full support for PHP with Tree-sitter based parsing.

## Supported Features

- Class definitions
- Interface definitions
- Trait definitions
- Function definitions
- Method definitions
- Namespaces
- Use statements
- PHPDoc comments
- Attributes (PHP 8+)

## Example

### Input

```php
<?php

namespace App\Models;

use DateTimeInterface;

/**
 * Represents a system user
 */
class User
{
    private int $id;
    private string $name;
    private string $email;

    public function __construct(int $id, string $name, string $email)
    {
        $this->id = $id;
        $this->name = $name;
        $this->email = $email;
    }

    /**
     * Get the user's ID
     */
    public function getId(): int
    {
        return $this->id;
    }

    /**
     * Validate the user's email
     */
    public function hasValidEmail(): bool
    {
        return filter_var($this->email, FILTER_VALIDATE_EMAIL) !== false;
    }
}

interface UserRepositoryInterface
{
    public function findById(int $id): ?User;
    public function save(User $user): User;
}
```

### Output (Brf.it)

```php
// App/Models/User.php
<?php
namespace App\Models;
use DateTimeInterface;

class User {
    public function __construct(int $id, string $name, string $email)
    public function getId(): int
    public function hasValidEmail(): bool
}

interface UserRepositoryInterface {
    public function findById(int $id): ?User
    public function save(User $user): User
}
```

## Extension

Files with `.php` extension are processed.
