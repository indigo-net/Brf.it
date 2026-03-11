---
name: product-meeting
description: Use when planning the next version. Analyzes project from 3 perspectives (PO, Planner, CTO), creates GitHub milestone + issues automatically. Also use when user says "제품 회의", "다음 버전 기획", "/product-meeting".
---

# Product Meeting — 3관점 자동 기획 파이프라인

## Overview

PO, 기획자, CTO 3가지 관점에서 프로젝트를 병렬 분석하고, 통합 에이전트가 교차검증한 결과를 GitHub 마일스톤 + 이슈로 자동 등록한다.

```
Stage 0: 컨텍스트 수집 (메인)
    ↓
Stage 1: 3개 병렬 에이전트 분석 (PO + 기획자 + CTO)
    ↓
Stage 2: 통합 에이전트 (교차검증 + 우선순위 결정)
    ↓  [HARD-GATE]
Stage 3: 마일스톤 생성 + 이슈 일괄 등록 (메인)
```

**완전 자동** — 사용자 승인 없이 마일스톤과 이슈를 생성한다.

---

## HARD-GATE

**Stage 2 통합 완료 전 `gh issue create` 또는 `gh api` 마일스톤 생성 절대 금지.**

통합 에이전트의 교차검증과 중복 제거를 거치지 않은 항목은 이슈로 생성하지 않는다. 이 게이트를 우회하는 것은 어떤 이유로도 허용되지 않는다.

---

## Stage 0: 컨텍스트 수집

메인에서 직접 수행한다. 아래 데이터를 수집하여 하나의 컨텍스트 블록으로 패키징한다.

### 수집 명령어

```bash
# 1. 오픈 이슈 전체
gh issue list --state open --json number,title,state,labels,assignees

# 2. 최근 클로즈 이슈 20개
gh issue list --state closed --limit 20 --json number,title,state,labels

# 3. 최근 커밋 20개
git log --oneline -20
```

### 파일 읽기

- `CHANGELOG.md` — 버전 히스토리, 릴리즈 간격
- `SAMPLE.md` — 현재 brfit 출력 형태 (코드 구조 파악)
- `README.md` — 제품 소개, 기능 목록

### HARD-GATE 0

- `gh issue list`가 실패하면 (인증 오류, 네트워크) **즉시 중단**
- `CHANGELOG.md`가 없거나 파싱 불가능하면 **즉시 중단**

### 컨텍스트 블록 형식

수집한 데이터를 아래 형식으로 패키징한다:

```
=== OPEN ISSUES ===
[gh issue list 결과]

=== CLOSED ISSUES (RECENT 20) ===
[gh issue list --state closed 결과]

=== RECENT COMMITS ===
[git log 결과]

=== CHANGELOG ===
[CHANGELOG.md 내용]

=== CURRENT OUTPUT (SAMPLE.md) ===
[SAMPLE.md 내용]

=== README ===
[README.md 내용]
```

---

## Stage 1: 병렬 에이전트 분석

3개의 Agent를 **동시에** 실행한다. 각 에이전트에게 Stage 0 컨텍스트 블록을 전달한다.

### Agent 디스패치

```
Agent 1: po-analyst (subagent_type 미사용, general-purpose agent에 po-analyst.md 역할 프롬프트 전달)
Agent 2: planner-analyst
Agent 3: cto-analyst
```

각 에이전트 실행 시 아래 형식으로 프롬프트를 구성한다:

```
당신은 Brf.it 프로젝트의 [역할명]입니다.

아래 프로젝트 컨텍스트를 분석하고, [에이전트 .md 파일의 분석 프로세스]에 따라
3-7개의 제안을 JSON 형식으로 출력하세요.

[Stage 0 컨텍스트 블록 삽입]

출력 형식은 반드시 다음 JSON을 따르세요:
{
  "role": "[PO|Planner|CTO]",
  "proposals": [...]
}
```

**PO 에이전트 특이사항**: WebSearch 도구를 사용하여 경쟁 도구(repomix, code2prompt 등) 최신 기능을 조사한다.

### HARD-GATE 1

- 3개 에이전트 중 **최소 2개**가 유효한 JSON을 출력해야 진행 가능
- 1개 실패 시: 해당 에이전트 **1회 재시도** (JSON 형식 준수 명시)
- 2개 이상 실패 시: **즉시 중단**

---

## Stage 2: 통합 에이전트

Stage 1의 3개 결과를 모두 수집한 후, **meeting-integrator** 에이전트를 실행한다.

### 프롬프트 구성

```
당신은 Brf.it 프로젝트의 제품 회의 의장입니다.

아래 3개 에이전트의 분석 결과를 통합하세요.

[Stage 0 컨텍스트 중 오픈 이슈 목록 — 중복 검사용]

=== PO 분석 결과 ===
[Agent 1 출력]

=== 기획자 분석 결과 ===
[Agent 2 출력]

=== CTO 분석 결과 ===
[Agent 3 출력]

통합 규칙:
1. 2명 이상 제안한 항목은 우선순위 부스트
2. 기존 오픈 이슈와 중복되는 제안은 제외
3. CHANGELOG.md에서 현재 버전을 파싱하여 다음 버전 결정
4. 우선순위 정렬: 교차제안수 > impact > effort(낮은것 우선)

출력 형식은 반드시 다음 JSON을 따르세요:
{
  "milestone": { "title": "vX.Y.Z", "description": "..." },
  "issues": [...],
  "skipped": [...]
}
```

### HARD-GATE 2

- 통합 에이전트가 유효한 JSON을 출력해야 Stage 3 진행 가능
- 이슈 0건 (모두 중복 제외)은 유효한 결과 — "신규 항목 없음" 보고 후 종료
- JSON 파싱 실패 시 **즉시 중단**

---

## Stage 3: GitHub 등록

통합 에이전트의 JSON 결과를 파싱하여 마일스톤과 이슈를 생성한다.

### 3-1. 마일스톤 생성

```bash
gh api repos/{owner}/{repo}/milestones \
  -f title="[milestone.title]" \
  -f description="[milestone.description]"
```

- 이미 동일 제목의 마일스톤이 존재하면 기존 마일스톤 사용
- owner/repo는 `gh repo view --json owner,name`으로 확인

### 3-2. 이슈 일괄 생성

통합 결과의 `issues` 배열을 priority 순서대로 생성:

```bash
gh issue create \
  --title "[issue.title]" \
  --body "[issue.body]" \
  --assignee indigo-net \
  --label "[issue.labels 콤마 구분]" \
  --milestone "[milestone.title]"
```

각 이슈 생성 후 반환된 이슈 번호를 기록한다.

### 3-3. 에러 처리

- 마일스톤 생성 실패: 기존 마일스톤 조회 시도 → 실패 시 마일스톤 없이 이슈 생성
- 개별 이슈 생성 실패: 로그에 기록하고 다음 이슈로 계속
- 모든 이슈 실패 시: 에러 요약 출력

---

## 출력 요약

모든 단계 완료 후 아래 형식으로 요약을 출력한다:

```
## Product Meeting Summary

### 마일스톤: [title]
- 이슈 N개 생성, M개 중복 제외

### 에이전트 통계
| 에이전트 | 제안 | 채택 | 중복 제외 |
|----------|------|------|----------|
| PO       | N    | N    | N        |
| 기획자   | N    | N    | N        |
| CTO      | N    | N    | N        |

### 생성된 이슈
| # | 이슈 번호 | 제목 | 라벨 | 제안 에이전트 |
|---|-----------|------|------|--------------|
| 1 | #XX       | ... | ... | PO, CTO      |

### 교차 검증 하이라이트
- [2+ 에이전트가 독립적으로 제안한 항목 강조]

### 제외 항목
| # | 제목 | 사유 |
|---|------|------|
| 1 | ... | 중복: #XX |
```

---

## Common Mistakes

| 실수 | 방지책 |
|------|--------|
| Stage 2 완료 전 이슈 생성 | HARD-GATE: 통합 에이전트 JSON 출력 전 `gh issue create` 금지 |
| 기존 이슈 중복 체크 누락 | Stage 0에서 오픈 이슈 수집, Stage 2에서 대조 |
| version_tag 라벨 누락 | 모든 이슈에 `version:major/minor/patch` 중 하나 필수 |
| assignee 누락 | 모든 이슈에 `--assignee indigo-net` 필수 |
| 에이전트 JSON 파싱 실패 무시 | HARD-GATE: 최소 2/3 에이전트 유효 출력 필수 |
| 마일스톤 없이 이슈 생성 | 마일스톤 생성을 이슈 생성보다 먼저 수행 |

## Red Flags — 멈추고 재확인

- Stage 1 결과를 그대로 이슈로 생성하려는 충동
- "통합 에이전트는 생략해도 된다"
- 기존 이슈 목록을 확인하지 않고 진행
- 에이전트 출력이 JSON이 아닌데 무시하고 진행
- 한 에이전트의 제안만으로 이슈 우선순위 결정
