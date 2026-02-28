# TypeScript Support

🌐 [English](typescript.md) | [한국어](../ko/languages/typescript.md) | [日本語](../ja/languages/typescript.md) | [हिन्दी](../hi/languages/typescript.md) | [Deutsch](../de/languages/typescript.md)

## Supported Extensions

- `.ts`
- `.tsx`
- `.js` (JavaScript)
- `.jsx` (JSX)

## Extraction Targets

| Element | Kind | Example |
|---------|------|---------|
| Function declaration | `function` | `function greet()` |
| Arrow function | `arrow` | `const greet = () => {}` |
| Method | `method` | `class A { method() {} }` |
| Class | `class` | `class User {}` |
| Interface | `interface` | `interface Props {}` |
| Type alias | `type` | `type ID = string` |
| Module-level const/let | `variable` | `const API_URL = "..."` |
| Comment | `doc` | `// Comment` |

## Example

### Input

```typescript
/**
 * User interface representing a user entity.
 */
interface User {
    id: number;
    name: string;
}

/**
 * Creates a new user with the given name.
 */
function createUser(name: string): User {
    return { id: Date.now(), name };
}

/**
 * User service for managing users.
 */
class UserService {
    private users: User[] = [];

    /**
     * Adds a user to the service.
     */
    addUser(user: User): void {
        this.users.push(user);
    }
}

// Arrow function example
const formatName = (user: User): string => {
    return user.name.toUpperCase();
};
```

### Output (XML)

```xml
<file path="user.ts" language="typescript">
  <signature kind="interface" line="4">
    <name>User</name>
    <text>interface User</text>
    <doc>User interface representing a user entity.</doc>
  </signature>
  <signature kind="function" line="12">
    <name>createUser</name>
    <text>function createUser(name: string): User</text>
    <doc>Creates a new user with the given name.</doc>
  </signature>
  <signature kind="class" line="18">
    <name>UserService</name>
    <text>class UserService</text>
    <doc>User service for managing users.</doc>
  </signature>
  <signature kind="method" line="24">
    <name>addUser</name>
    <text>addUser(user: User): void</text>
    <doc>Adds a user to the service.</doc>
  </signature>
  <signature kind="arrow" line="30">
    <name>formatName</name>
    <text>const formatName = (user: User): string => </text>
    <doc>Arrow function example</doc>
  </signature>
</file>
```

## Notes

### Export Detection

- Only elements with `export` keyword are extracted (default)
- For JavaScript files, all elements are extracted

### Arrow Functions

- Arrow functions declared with `const`/`let`/`var` are captured
- Variable name is used as function name
- `export const` format is also supported

### Body Removal

When `--include-body` flag is not used:

- Functions/Methods: body removed after opening brace `{`
- Arrow functions: body removed after `=>`
- Classes/Interfaces: content removed after opening brace `{`

### JSDoc Support

- `/** ... */` style JSDoc comments are automatically linked
- Comments immediately before functions/classes are captured as doc

### JavaScript Compatibility

- `.js`, `.jsx` files are processed with TypeScript parser
- Functions/classes can be extracted even without type information
