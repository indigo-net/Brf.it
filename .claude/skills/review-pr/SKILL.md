---
name: review-pr
description: PR을 프로젝트 컨벤션 기준으로 리뷰하고, 리뷰 결과에 응답(수정 또는 반박)하는 2단계 워크플로우
argument-hint: "[PR-number or leave blank for current branch]"
---

# PR Review Automation

2단계 자동화 리뷰 워크플로우를 실행합니다.

## PR 컨텍스트 프리로드

아래 명령으로 PR 정보를 수집하세요:

```
PR_NUMBER=$ARGUMENTS

# PR 번호가 없으면 현재 브랜치의 PR을 찾기
if [ -z "$PR_NUMBER" ]; then
  PR_NUMBER=$(gh pr view --json number -q .number 2>/dev/null)
fi
```

### PR 메타데이터

!`gh pr view $ARGUMENTS --json title,body,baseRefName,headRefName,number,url`

### PR Diff

!`gh pr diff $ARGUMENTS`

### 기존 코멘트

!`gh pr view $ARGUMENTS --comments`

---

## Phase 1: 코드 리뷰 (pr-reviewer)

위에서 수집한 PR 컨텍스트를 기반으로 `pr-reviewer` 에이전트에 리뷰를 위임하세요.

**위임 방법:**
```
Agent(pr-reviewer)에게 다음을 전달:
- PR 번호: $PR_NUMBER
- PR diff (위 프리로드 결과)
- PR 메타데이터 (위 프리로드 결과)
- 기존 코멘트 (위 프리로드 결과)
- 리포지토리 owner/repo: gh repo view --json nameWithOwner -q .nameWithOwner 로 확인
```

리뷰어 에이전트는 리뷰 결과를 `gh pr review --comment` 으로 PR에 게시합니다.

**에이전트 완료 대기 후 결과를 확인하세요.**

---

## Phase 2: 리뷰 응답 (pr-responder)

리뷰어가 게시한 코멘트를 기반으로 `pr-responder` 에이전트에 응답을 위임하세요.

**위임 방법:**
```
Agent(pr-responder)에게 다음을 전달:
- PR 번호: $PR_NUMBER
- 리뷰어가 게시한 리뷰 코멘트 내용
- PR의 head 브랜치명 (push 대상)
- 리포지토리 owner/repo
```

응답자 에이전트는:
- 유효한 지적: 코드 수정 후 커밋 + push
- 잘못된 지적: 반박 코멘트 게시

---

## 완료 후 요약

두 에이전트가 모두 완료되면 아래 형식으로 요약을 출력하세요:

```
## PR Review Summary

**PR**: #$PR_NUMBER - $PR_TITLE
**리뷰어 판정**: [Ready to merge / Needs fixes / ...]

### 응답 결과
- 수정 반영: N건
- 반박 코멘트: N건
- 커밋: [커밋 해시 목록 또는 "없음"]

### 상세
| # | 이슈 | 심각도 | 응답 | 상태 |
|---|------|--------|------|------|
| 1 | [이슈 요약] | Critical/Important/Minor | 수정/반박 | 완료 |
```
