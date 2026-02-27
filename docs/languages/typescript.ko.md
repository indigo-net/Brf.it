# TypeScript 지원

🌐 [English](typescript.md) | [한국어](typescript.ko.md) | [日本語](typescript.ja.md) | [हिन्दी](typescript.hi.md) | [Deutsch](typescript.de.md)

## 지원 확장자

- `.ts`
- `.tsx`
- `.js` (JavaScript)
- `.jsx` (JSX)

## 추출 대상

| 요소 | Kind | 예시 |
|------|------|------|
| 함수 선언 | `function` | `function greet()` |
| 화살표 함수 | `arrow` | `const greet = () => {}` |
| 메서드 | `method` | `class A { method() {} }` |
| 클래스 | `class` | `class User {}` |
| 인터페이스 | `interface` | `interface Props {}` |
| 타입 별칭 | `type` | `type ID = string` |
| 모듈 레벨 const/let | `variable` | `const API_URL = "..."` |
| 주석 | `doc` | `// Comment` |

## 예시

### 입력

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

### 출력 (XML)

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

## 특이사항

### Export 판별

- `export` 키워드가 있는 요소만 추출 (기본값)
- JavaScript 파일의 경우 모든 요소 추출

### 화살표 함수

- `const`/`let`/`var`로 선언된 화살표 함수 캡처
- 변수명을 함수명으로 사용
- `export const` 형태도 지원

### 본문 제거

`--include-body` 플래그 미사용 시:

- 함수/메서드: 중괄호 `{` 이후 본문 제거
- 화살표 함수: `=>` 이후 본문 제거
- 클래스/인터페이스: 중괄호 `{` 이후 내용 제거

### JSDoc 지원

- `/** ... */` 형태의 JSDoc 주석 자동 연결
- 함수/클래스 직전 주석이 doc으로 캡처됨

### JavaScript 호환

- `.js`, `.jsx` 파일은 TypeScript 파서로 처리
- 타입 정보가 없어도 함수/클래스 추출 가능
