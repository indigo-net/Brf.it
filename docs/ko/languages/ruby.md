# Ruby 지원

[English](../../languages/ruby.md) | [한국어](ruby.md) | [日本語](../../ja/languages/ruby.md) | [हिन्दी](../../hi/languages/ruby.md) | [Deutsch](../../de/languages/ruby.md)

## 지원 확장자

- `.rb`

## 추출 대상

| 요소 | Kind | 예시 |
|------|------|------|
| 메서드 | `method` | `def greet(name)` |
| 클래스 메서드 | `class_method` | `def self.create(attrs)` |
| 클래스 | `class` | `class User < ActiveRecord::Base` |
| 모듈 | `module` | `module Authentication` |
| 상수 | `variable` | `MAX_RETRIES = 3` |
| YARD 주석 | `doc` | `# 설명` |
| require | `import` | `require "json"` |
| require_relative | `import` | `require_relative "helpers"` |

## 예시

### 입력

```ruby
require "json"
require_relative "helpers"

# 시스템의 사용자를 나타냅니다.
class User
  MAX_RETRIES = 3

  # 속성에서 새 사용자를 생성합니다.
  # @param attrs [Hash] 사용자 속성
  def self.create(attrs)
    new(attrs).save
  end

  # 사용자를 초기화합니다.
  # @param name [String] 사용자 이름
  def initialize(name)
    @name = name
  end

  # 다른 사람에게 인사합니다.
  # @param other [String] 상대방 이름
  # @return [String] 인사 메시지
  def greet(other)
    "Hello, #{other}! I'm #{@name}."
  end
end

module Authentication
  def authenticate(password)
    password == @secret
  end
end
```

### 출력 (XML)

```xml
<file path="example.rb" language="ruby">
  <class kind="class" line="5">
    <name>User</name>
    <text>class User</text>
  </class>
  <variable kind="variable" line="6">
    <name>MAX_RETRIES</name>
    <text>MAX_RETRIES = 3</text>
  </variable>
  <function kind="class_method" line="10">
    <name>create</name>
    <text>def self.create(attrs)</text>
  </function>
  <function kind="method" line="15">
    <name>initialize</name>
    <text>def initialize(name)</text>
  </function>
  <function kind="method" line="21">
    <name>greet</name>
    <text>def greet(other)</text>
  </function>
  <module kind="module" line="27">
    <name>Authentication</name>
    <text>module Authentication</text>
  </module>
  <function kind="method" line="28">
    <name>authenticate</name>
    <text>def authenticate(password)</text>
  </function>
</file>
```

## 참고사항

### 가시성 (Visibility)

- 가시성(`public`, `protected`, `private`)에 관계없이 모든 메서드가 추출됩니다
- 인스턴스 메서드(`def foo`)와 클래스 메서드(`def self.foo`) 모두 캡처됩니다

### 메서드 종류

- `method`: 인스턴스 메서드 선언 (`def foo`)
- `class_method`: 클래스 레벨 메서드 선언 (`def self.foo`)

### 본문 제거

`--include-body` 플래그를 사용하지 않을 때:

- 메서드: 파라미터 리스트의 닫는 괄호 `)` 이후 본문 제거 (파라미터가 없으면 메서드 이름 이후)
- 클래스/모듈: 선언 줄만 보존
- 상수: 그대로 보존

### Import 추출

- `require`와 `require_relative` 구문은 `--include-imports` 플래그로 추출 가능
- 형식: `require "json"` / `require_relative "helpers"` (전체 구문 보존)

### 문서 주석

- YARD 스타일 주석(`#`)이 메서드/클래스 바로 위에 있으면 추출됩니다
- 여러 줄 `#` 주석도 지원됩니다
- `=begin`...`=end` 블록 주석도 인식됩니다
