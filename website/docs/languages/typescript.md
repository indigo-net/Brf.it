---
sidebar_position: 2
title: TypeScript
---

# TypeScript Support

Brf.it provides full support for TypeScript with Tree-sitter based parsing.

## Supported Features

- Function declarations
- Arrow functions
- Class definitions
- Interface definitions
- Type aliases
- Enums
- Generics
- Import/export statements (with `--include-imports`)
- JSDoc comments

## Example

### Input

```typescript
import { Request, Response } from 'express';

interface User {
  id: number;
  name: string;
  email: string;
}

interface UserRepository {
  findById(id: number): Promise<User | null>;
  save(user: User): Promise<User>;
}

/**
 * UserService handles business logic for user operations
 */
export class UserService implements UserRepository {
  constructor(private db: Database) {}

  async findById(id: number): Promise<User | null> {
    return this.db.query('SELECT * FROM users WHERE id = $1', [id]);
  }

  async save(user: User): Promise<User> {
    // implementation...
    return user;
  }
}
```

### Output (Brf.it)

```typescript
// user.service.ts
import { Request, Response } from 'express';

interface User { id: number; name: string; email: string }
interface UserRepository { findById(id: number): Promise<User | null>; save(user: User): Promise<User> }
export class UserService implements UserRepository {
  constructor(private db: Database)
  async findById(id: number): Promise<User | null>
  async save(user: User): Promise<User>
}
```

## Extensions

Files with `.ts` and `.tsx` extensions are processed.
