# TypeScript सपोर्ट

🌐 [English](../../languages/typescript.md) | [한국어](../../ko/languages/typescript.md) | [日本語](../../ja/languages/typescript.md) | [हिन्दी](typescript.md) | [Deutsch](../../de/languages/typescript.md)

## समर्थित एक्सटेंशन

- `.ts`
- `.tsx`
- `.js` (JavaScript)
- `.jsx` (JSX)

## निष्कर्षण लक्ष्य

| तत्व | Kind | उदाहरण |
|------|------|--------|
| फंक्शन डिक्लेरेशन | `function` | `function greet()` |
| एरो फंक्शन | `arrow` | `const greet = () => {}` |
| मेथड | `method` | `class A { method() {} }` |
| क्लास | `class` | `class User {}` |
| इंटरफेस | `interface` | `interface Props {}` |
| टाइप एलियास | `type` | `type ID = string` |
| मॉड्यूल-लेवल const/let | `variable` | `const API_URL = "..."` |
| कमेंट | `doc` | `// Comment` |

## उदाहरण

### इनपुट

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

### आउटपुट (XML)

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

## नोट्स

### एक्सपोर्ट डिटेक्शन

- केवल `export` कीवर्ड वाले तत्व निकाले जाते हैं (डिफ़ॉल्ट)
- JavaScript फाइलों के लिए सभी तत्व निकाले जाते हैं

### एरो फंक्शन

- `const`/`let`/`var` से घोषित एरो फंक्शन कैप्चर होते हैं
- वेरिएबल नाम को फंक्शन नाम के रूप में उपयोग किया जाता है
- `export const` फॉर्मेट भी समर्थित है

### बॉडी रिमूवल

`--include-body` फ्लैग का उपयोग न करने पर:

- फंक्शन/मेथड: ओपनिंग ब्रेस `{` के बाद बॉडी हटा दी जाती है
- एरो फंक्शन: `=>` के बाद बॉडी हटा दी जाती है
- क्लास/इंटरफेस: ओपनिंग ब्रेस `{` के बाद कंटेंट हटा दिया जाता है

`--include-private` का उपयोग करके गैर-निर्यातित/निजी सिंबल शामिल करें।

### JSDoc सपोर्ट

- `/** ... */` स्टाइल JSDoc कमेंट्स ऑटोमैटिकली लिंक होते हैं
- फंक्शन/क्लास से पहले के कमेंट्स doc के रूप में कैप्चर होते हैं

### JavaScript कम्पैटिबिलिटी

- `.js`, `.jsx` फाइलें TypeScript पार्सर से प्रोसेस होती हैं
- टाइप इंफॉर्मेशन के बिना भी फंक्शन/क्लास निकाले जा सकते हैं
