# Code Summary: examples/typescript

*brf.it dev*

---

## Files

### examples/typescript/app.ts

```typescript
interface AppConfig {
  port: number;
  host: string;
  debug: boolean;
}
interface User {
  id: string;
  email: string;
  name: string;
  createdAt: Date;
}
export function createDefaultConfig(): AppConfig
class Repository<T extends { id: string }>
async findById(id: string): Promise<T | undefined>
async save(item: T): Promise<T>
async delete(id: string): Promise<boolean>
async findAll(): Promise<T[]>
export const formatUser = (user: User): string
export function isUser(value: unknown): value is User
```

---

