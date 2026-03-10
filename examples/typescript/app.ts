/** Configuration options for the app. */
export interface AppConfig {
  port: number;
  host: string;
  debug: boolean;
}

/** User entity with authentication data. */
export interface User {
  id: string;
  email: string;
  name: string;
  createdAt: Date;
}

/** Create a default configuration. */
export function createDefaultConfig(): AppConfig {
  return { port: 3000, host: "localhost", debug: false };
}

/** Generic repository for CRUD operations. */
export class Repository<T extends { id: string }> {
  private items: Map<string, T> = new Map();

  /** Find an item by its ID. */
  async findById(id: string): Promise<T | undefined> {
    return this.items.get(id);
  }

  /** Save or update an item. */
  async save(item: T): Promise<T> {
    this.items.set(item.id, item);
    return item;
  }

  /** Delete an item by ID. Returns true if found and deleted. */
  async delete(id: string): Promise<boolean> {
    return this.items.delete(id);
  }

  /** List all items. */
  async findAll(): Promise<T[]> {
    return Array.from(this.items.values());
  }
}

/** Format a user's display name. */
export const formatUser = (user: User): string =>
  `${user.name} <${user.email}>`;

/** Type guard to check if a value is a valid User. */
export function isUser(value: unknown): value is User {
  return (
    typeof value === "object" &&
    value !== null &&
    "id" in value &&
    "email" in value
  );
}
