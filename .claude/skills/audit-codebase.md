---
name: audit-codebase
description: Use when auditing the codebase for issues - performance bottlenecks, error handling gaps, uncontrollable factors, token optimization. Also use when user says "코드 점검", "이슈 찾기", "코드 감사", or wants to discover and file GitHub issues for improvements.
---

# Audit Codebase

## Overview

3단계 파이프라인으로 코드베이스를 감사하고, 검증된 이슈만 GitHub에 생성한다.

**핵심 원칙: 검증 에이전트 승인 없이 이슈 생성 금지**

```
탐색 (3 Explore agents) → 검증 (1 agent) → 이슈 생성
```

## HARD-GATE

**Stage 2 검증 완료 전 `gh issue create` 실행 절대 금지.**

검증 단계에서 필터링되지 않은 항목은 이슈로 생성하지 않는다. 이 게이트를 우회하는 것은 어떤 이유로도 허용되지 않는다.

---

## Stage 1: 병렬 탐색

3개의 Explore 에이전트를 **동시에** 실행한다. 각 에이전트에게 아래 프롬프트를 전달하세요.

### Agent 1: 성능 병목 + 동시성 패턴

```
Brf.it 코드베이스에서 성능 병목과 동시성 문제를 탐색하세요.

탐색 대상 디렉토리:
- pkg/scanner/ (파일 시스템 스캔, goroutine/channel 사용)
- pkg/parser/treesitter/ (Tree-sitter CGO 바인딩, 파서 풀링)
- pkg/extractor/ (시그니처 추출 파이프라인)
- pkg/formatter/ (XML/Markdown 출력 생성)
- internal/ (내부 유틸리티)
- cmd/brfit/ (CLI 진입점)

탐색 항목:
1. 불필요한 메모리 할당 (큰 슬라이스/맵 재할당, 불필요한 복사)
2. goroutine 누수 가능성 (닫히지 않는 채널, 대기 상태 고루틴)
3. CGO 호출 오버헤드 (빈번한 Go↔C 경계 횡단)
4. Tree-sitter 리소스 해제 누락 (Parser, Tree, Query, QueryCursor의 Close/defer)
5. 동시성 안전성 (공유 상태 접근, race condition 가능성)
6. I/O 병목 (파일 읽기 패턴, 버퍼링 전략)

출력 형식:
각 발견 항목마다:
- 파일 경로 + 라인 번호
- 문제 설명 (구체적으로)
- 심각도 추정 (CRITICAL / HIGH / MEDIUM / LOW)
- 개선 방안
```

### Agent 2: 에러 처리 + 엣지 케이스 + CGO 안전성

```
Brf.it 코드베이스에서 에러 처리 갭과 엣지 케이스를 탐색하세요.

탐색 대상 디렉토리:
- pkg/scanner/ (파일 시스템 에러, 권한 문제)
- pkg/parser/treesitter/ (파싱 실패, 잘못된 입력)
- pkg/extractor/ (추출 실패 케이스)
- pkg/formatter/ (출력 생성 에러)
- internal/ (내부 에러 전파)
- cmd/brfit/ (CLI 에러 처리, 사용자 입력 검증)

탐색 항목:
1. 무시된 에러 반환값 (err가 _ 로 버려지는 곳)
2. 에러 래핑 없이 반환 (fmt.Errorf + %w 미사용)
3. nil 포인터 역참조 가능성 (Tree-sitter 노드 접근 시)
4. CGO 안전성 (C 메모리 접근, 해제 후 사용)
5. 경계값 처리 (빈 파일, 0바이트, 매우 긴 줄, 바이너리 파일)
6. 심볼릭 링크 순환 참조 처리
7. 동시 파일 시스템 변경 (스캔 중 파일 삭제/수정)
8. panic 가능성 (인덱스 범위 초과, 타입 단언 실패)

출력 형식:
각 발견 항목마다:
- 파일 경로 + 라인 번호
- 문제 설명 (구체적으로)
- 재현 시나리오 (어떤 입력이 문제를 유발하는지)
- 심각도 추정 (CRITICAL / HIGH / MEDIUM / LOW)
- 개선 방안
```

### Agent 3: 토큰 최적화 + 기존 이슈 수집

```
Brf.it 코드베이스에서 토큰 최적화 기회를 탐색하고, 기존 GitHub 이슈 목록을 수집하세요.

Part A — 토큰 최적화 탐색:

탐색 대상:
- pkg/formatter/ (XML/Markdown 출력 효율성)
- pkg/extractor/ (불필요한 정보 포함 여부)
- pkg/parser/treesitter/languages/ (쿼리 패턴 효율성)

탐색 항목:
1. 출력에 불필요한 공백/줄바꿈이 포함되는 패턴
2. 중복 정보 출력 (같은 시그니처가 여러 번 나오는 경우)
3. 더 압축된 출력 형식 가능성
4. import 문 최적화 (중복 import, 불필요한 import 포함)
5. 빈 파일/빈 시그니처 처리 효율성

Part B — 기존 이슈 수집:

아래 명령어를 실행하여 기존 이슈 목록을 수집하세요:
  gh issue list --state all --limit 100 --json number,title,state,labels

이 목록은 Stage 2 검증에서 중복 체크에 사용됩니다.

출력 형식:
Part A: 각 발견 항목마다 파일 경로 + 라인 번호 + 문제 설명 + 개선 방안
Part B: 기존 이슈 번호/제목/상태 전체 목록
```

---

## Stage 2: 검증 에이전트

Stage 1의 3개 에이전트 결과를 모두 수집한 후, **단일 검증 에이전트**를 실행한다.

### 검증 기준

각 발견 항목에 대해 4가지 검증을 수행:

| # | 검증 | 방법 | 실패 시 |
|---|------|------|---------|
| 1 | 기존 이슈 중복 | Agent 3의 이슈 목록과 대조 | 제외 (이슈 번호 기록) |
| 2 | 코드 경로 존재 | 파일 경로 + 라인 번호 실제 확인 | 제외 |
| 3 | 코드 기반 타당성 | 실제 코드를 읽고 문제가 존재하는지 확인 | 제외 |
| 4 | 심각도 분류 | CRITICAL / HIGH / MEDIUM / LOW 재평가 | 심각도 조정 |

### 검증 에이전트 프롬프트

```
아래는 코드베이스 감사에서 발견된 항목 목록입니다. 각 항목을 검증하세요.

[여기에 Stage 1 결과물을 삽입]

검증 절차:
1. 기존 GitHub 이슈와 중복 여부 확인 (아래 이슈 목록 참조)
2. 파일 경로와 라인 번호가 실제로 존재하는지 확인 (Read 도구 사용)
3. 해당 코드를 직접 읽고, 보고된 문제가 실제로 존재하는지 판단
4. 심각도를 재평가 (CRITICAL / HIGH / MEDIUM / LOW)

기존 이슈 목록:
[Agent 3에서 수집한 이슈 목록 삽입]

출력 형식:

## 검증 통과 항목
각 항목마다:
- 원본 ID
- 파일 경로:라인
- 카테고리 (performance / error-handling / token-optimization)
- 심각도
- 요약
- 이슈 제목 (한글, 간결하게)
- 이슈 라벨 (bug 또는 enhancement)

## 검증 실패 항목 (제외 사유 포함)
각 항목마다:
- 원본 ID
- 제외 사유 (중복:#XX / 경로 불일치 / 문제 미존재 / 기타)
```

---

## Stage 3: 이슈 일괄 생성

검증 통과 항목만 GitHub 이슈로 생성한다.

### 이슈 생성 규칙

- `--assignee indigo-net` 필수
- `--label` 적절한 라벨 (bug / enhancement / performance)
- 본문 형식:

```markdown
## 개요
[문제 요약 1-2줄]

## 현재 상태
[현재 코드가 어떻게 동작하는지]

## 영향
[이 문제가 미치는 영향: 성능, 안정성, UX 등]

## 개선 방안
[구체적인 해결 방법]

## 관련 파일
- `파일경로:라인` — 설명
```

### 생성 후 사용자 확인

이슈 생성 전에 검증 통과 항목 목록을 사용자에게 보여주고 승인을 받는다.

---

## 출력 포맷

모든 단계 완료 후 아래 형식으로 요약을 출력한다:

```
## Audit Summary

### 발견 항목 통계
| 카테고리 | 발견 | 검증 통과 | 이슈 생성 |
|----------|------|-----------|-----------|
| 성능     | N    | N         | N         |
| 에러 처리 | N   | N         | N         |
| 토큰 최적화 | N | N         | N         |
| **합계** | N    | N         | N         |

### 생성된 이슈
| # | 이슈 번호 | 제목 | 심각도 | 라벨 |
|---|-----------|------|--------|------|
| 1 | #XX | ... | HIGH | enhancement |

### 제외된 항목
| # | 사유 | 설명 |
|---|------|------|
| 1 | 중복 (#XX) | ... |

### 기존 이슈와의 관계
- #XX (기존) ← 관련: #YY (신규)
```

---

## Common Mistakes

| 실수 | 방지책 |
|------|--------|
| 기존 이슈 중복 체크 누락 | Agent 3에서 반드시 `gh issue list` 수집, Stage 2에서 대조 |
| 검증 없이 이슈 바로 생성 | HARD-GATE: Stage 2 완료 전 `gh issue create` 금지 |
| 에이전트 프롬프트에 파일 경로 누락 | 위 템플릿의 탐색 대상 디렉토리 목록을 반드시 포함 |
| 심각도 과대평가 | Stage 2에서 실제 코드를 읽고 심각도 재평가 |
| 이슈 본문에 코드 경로 누락 | 관련 파일 섹션에 `파일:라인` 형식으로 명시 |
| 너무 많은 LOW 이슈 생성 | LOW는 사용자 확인 시 선별적으로 제외 가능 |

## Red Flags — 멈추고 재확인

- Stage 1 결과를 그대로 이슈로 생성하려는 충동
- "검증은 이미 에이전트가 했으니 Stage 2는 생략해도 된다"
- 기존 이슈 목록을 수집하지 않고 진행
- 파일 경로를 확인하지 않고 이슈 본문 작성
- 사용자 승인 없이 이슈 일괄 생성
