# TypeScript Unterstützung

🌐 [English](typescript.md) | [한국어](typescript.ko.md) | [日本語](typescript.ja.md) | [हिन्दी](typescript.hi.md) | [Deutsch](typescript.de.md)

## Unterstützte Erweiterungen

- `.ts`
- `.tsx`
- `.js` (JavaScript)
- `.jsx` (JSX)

## Extraktionsziele

| Element | Kind | Beispiel |
|---------|------|----------|
| Funktionsdeklaration | `function` | `function greet()` |
| Arrow-Funktion | `arrow` | `const greet = () => {}` |
| Methode | `method` | `class A { method() {} }` |
| Klasse | `class` | `class User {}` |
| Interface | `interface` | `interface Props {}` |
| Typ-Alias | `type` | `type ID = string` |
| Modul-Level const/let | `variable` | `const API_URL = "..."` |
| Kommentar | `doc` | `// Comment` |

## Beispiel

### Eingabe

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

### Ausgabe (XML)

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

## Hinweise

### Export-Erkennung

- Nur Elemente mit `export`-Schlüsselwort werden extrahiert (Standard)
- Bei JavaScript-Dateien werden alle Elemente extrahiert

### Arrow-Funktionen

- Mit `const`/`let`/`var` deklarierte Arrow-Funktionen werden erfasst
- Variablenname wird als Funktionsname verwendet
- `export const`-Format wird ebenfalls unterstützt

### Body-Entfernung

Wenn `--include-body` Flag nicht verwendet wird:

- Funktionen/Methoden: Body nach öffnender Klammer `{` entfernt
- Arrow-Funktionen: Body nach `=>` entfernt
- Klassen/Interfaces: Inhalt nach öffnender Klammer `{` entfernt

### JSDoc-Unterstützung

- `/** ... */` Stil JSDoc-Kommentare werden automatisch verknüpft
- Kommentare direkt vor Funktionen/Klassen werden als doc erfasst

### JavaScript-Kompatibilität

- `.js`, `.jsx` Dateien werden mit TypeScript-Parser verarbeitet
- Funktionen/Klassen können auch ohne Typinformationen extrahiert werden
