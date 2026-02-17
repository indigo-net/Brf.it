# Brf.it

_"Don't just feed raw code. Brief it for your AI."_

---

## (한국어)

### Brf.it이란?

Brf.it은 AI 코딩 어시스턴트에게 코드를 더 효율적으로 전달하기 위한 CLI 도구입니다.

프로젝트 코드베이스에서 함수 시그니처와 문서만 추출하여, AI가 핵심 내용을 빠르게 파악할 수 있는 형태로 변환합니다. 불필요한 구현 디테일을 제거해 토큰 사용량을 크게 줄일 수 있습니다.

### 왜 필요한가요?

AI에게 큰 프로젝트의 코드를 전달할 때, 모든 파일을 그대로 복사하면:

- 컨텍스트가 너무 길어집니다
- 토큰이 낭비됩니다
- AI가 중요한 부분을 놓칠 수 있습니다

Brf.it은 이 문제를 해결합니다.

### 계획된 기능

- **지능형 코드 압축**: 함수 시그니처와 주석만 추출 (Tree-sitter 기반)
- **다양한 출력 포맷**: XML, Markdown 선택 가능
- **구조화된 메타데이터**: 디렉토리 트리와 심볼 리스트 포함
- **간편한 사용**: `npx brfit <path>` 한 줄로 실행

### 사용 예시 (계획)

```bash
# 기본 사용법
npx brfit ./src

# Markdown 포맷으로 출력
npx brfit ./src -f md

# 결과를 파일로 저장
npx brfit . -o output.xml
```

### 개발 현황

현재 개발 중입니다. MVP 완성 후 npm에 배포될 예정입니다.

---

## (English)

### What is Brf.it?

Brf.it is a CLI tool designed to deliver code to AI coding assistants more efficiently.

It extracts only function signatures and documentation from your codebase, transforming them into a format that AI can quickly understand. By removing unnecessary implementation details, it significantly reduces token usage.

### Why do you need it?

When sharing a large project's code with AI, copying all files directly leads to:

- Overwhelmingly long context
- Wasted tokens
- AI potentially missing important parts

Brf.it solves this problem.

### Planned Features

- **Smart Code Compression**: Extract only function signatures and comments (Tree-sitter based)
- **Multiple Output Formats**: Choose between XML and Markdown
- **Structured Metadata**: Includes directory tree and symbol list
- **Simple Usage**: Run with just `npx brfit <path>`

### Usage Examples (Planned)

```bash
# Basic usage
npx brfit ./src

# Output in Markdown format
npx brfit ./src -f md

# Save output to file
npx brfit . -o output.xml
```

### Development Status

Currently in development. Will be published to npm after MVP completion.
