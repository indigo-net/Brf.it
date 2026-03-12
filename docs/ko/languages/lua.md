# Lua 지원

[English](../../languages/lua.md) | [한국어](lua.md) | [日本語](../../ja/languages/lua.md) | [हिन्दी](../../hi/languages/lua.md) | [Deutsch](../../de/languages/lua.md)

## 지원 확장자

- `.lua`

## 추출 대상

| 요소 | Kind | 예시 |
|------|------|------|
| 전역 함수 | `function` | `function globalFunc(a, b)` |
| 로컬 함수 | `local_function` | `local function helper()` |
| 모듈 함수 | `module_function` | `function M.create(name)` |
| 메서드 | `method` | `function M:init(config)` |
| 테이블 할당 | `variable` | `local M = {}` |
| 함수 할당 | `variable` | `local foo = function() end` |
| LuaDoc 주석 | `doc` | `--- 설명` |
| 줄 주석 | `doc` | `-- 설명` |
| require() | `import` | `local json = require("json")` |

## 예시

### 입력

```lua
local M = {}

--- 이름으로 사람에게 인사합니다.
-- @param name string 사람의 이름
function M.greet(name)
    print("Hello, " .. name)
end

function M:init(config)
    self.config = config
end

local function helper()
    return 42
end

function globalFunc(a, b)
    return a + b
end

local json = require("json")
```

### 출력 (XML)

```xml
<file path="example.lua" language="lua">
  <variable kind="variable" line="1">
    <name>M</name>
    <text>local M = {}</text>
  </variable>
  <function kind="module_function" line="4">
    <name>greet</name>
    <text>function M.greet(name)</text>
  </function>
  <function kind="method" line="9">
    <name>init</name>
    <text>function M:init(config)</text>
  </function>
  <function kind="local_function" line="13">
    <name>helper</name>
    <text>local function helper()</text>
  </function>
  <function kind="function" line="17">
    <name>globalFunc</name>
    <text>function globalFunc(a, b)</text>
  </function>
</file>
```

## 참고사항

### 가시성 (Visibility)

- 모든 선언이 추출됩니다 (Lua는 접근 수식어가 없습니다)
- `local` 함수/변수도 전역 함수와 함께 포함됩니다

### 함수 종류

- `function`: 전역 함수 선언 (`function foo()`)
- `local_function`: 로컬 함수 선언 (`local function foo()`)
- `module_function`: 점 인덱스 함수 (`function M.foo()`)
- `method`: 콜론 인덱스 메서드 (`function M:foo()`)

### 본문 제거

`--include-body` 플래그를 사용하지 않을 때:

- 함수/메서드: 파라미터 리스트의 닫는 괄호 `)` 이후 본문 제거
- 테이블 할당 (`local M = {}`): 그대로 보존
- 함수 할당 (`local foo = function() end`): 그대로 보존

### Import 추출

- `require()` 호출은 `--include-imports` 플래그로 추출 가능
- `--include-private`를 사용하여 비공개/unexported 심볼 포함
- 형식: `local json = require("json")` (전체 구문 보존)

### 문서 주석

- LuaDoc 주석 (`---`)과 일반 줄 주석 (`--`) 모두 추출됩니다
- 블록 주석 (`--[[ ... ]]`)도 지원됩니다
