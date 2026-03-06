---
name: release
description: Use when releasing a new version of Brf.it. Handles preflight checks, version decision, CHANGELOG generation, tagging, CI monitoring, and release notes. Also use when user says "릴리즈", "배포", "release", or "/release".
---

# Release Pipeline

## Overview

6단계 파이프라인으로 안전한 릴리즈를 수행한다. 각 단계 사이에 GATE 조건이 있으며, 실패 시 즉시 중단한다.

```
Stage 0: 프리플라이트 → Stage 1: 버전 결정 → Stage 2: CHANGELOG → Stage 3: 태그 → Stage 4: CI 모니터링 → Stage 5: 릴리즈 노트
```

## HARD-GATE: CI 실패 시 인라인 수정 절대 금지

**CI 실패 시 GitHub Issue를 생성하고 즉시 릴리즈를 중단한다.**

인라인 수정, 재태그, 강제 푸시 등 어떤 수단으로든 실패한 릴리즈를 이어가는 것은 금지된다. 이 게이트를 우회하는 것은 어떤 이유로도 허용되지 않는다.

---

## Stage 0: 프리플라이트 검증

아래 5가지를 순서대로 확인한다. 하나라도 실패하면 중단.

### 0-1. 브랜치 확인

```bash
git branch --show-current
```

`main`이 아니면 중단: "릴리즈는 main 브랜치에서만 가능합니다."

### 0-2. 작업 트리 clean 확인

```bash
git status --porcelain
```

출력이 있으면 중단: "커밋되지 않은 변경사항이 있습니다."

### 0-3. 리모트 동기화

```bash
git fetch origin main
git rev-parse HEAD
git rev-parse origin/main
```

HEAD와 origin/main이 다르면 중단: "로컬과 리모트가 동기화되지 않았습니다. `git pull` 또는 `git push`가 필요합니다."

### 0-4. 테스트 실행

```bash
go test ./...
```

실패 시 → **GitHub Issue 생성 후 중단**:

```bash
gh issue create \
  --title "fix: 릴리즈 전 테스트 실패" \
  --body "## 개요
릴리즈 프리플라이트에서 \`go test ./...\` 실패.

## 테스트 출력
\`\`\`
(실패 로그 붙여넣기)
\`\`\`

## 다음 단계
테스트 수정 후 릴리즈 재시도" \
  --assignee indigo-net \
  --label bug
```

사용자에게 이슈 번호를 알리고 즉시 중단.

### 0-5. 빌드 확인

```bash
go build ./cmd/brfit
```

실패 시 → 0-4와 같은 방식으로 GitHub Issue 생성 후 중단.

### GATE 0

5가지 모두 통과 시에만 Stage 1 진행. 프리플라이트 결과를 사용자에게 요약 출력:

```
✅ 프리플라이트 완료
- 브랜치: main
- 작업 트리: clean
- 리모트 동기화: OK
- 테스트: PASS
- 빌드: OK
```

---

## Stage 1: 버전 결정

### 1-1. 마지막 태그 이후 커밋 분석

```bash
git describe --tags --abbrev=0   # 마지막 태그
git log <last-tag>..HEAD --oneline
```

커밋 메시지를 분석하여 버전을 제안:
- `BREAKING` 또는 `!:` 포함 → **major** bump
- `feat:` 포함 → **minor** bump
- `fix:`, `docs:`, `chore:` 등만 → **patch** bump

### 1-2. 인자로 버전 지정

`/release v0.19.0` 형태로 인자가 주어지면 해당 버전을 사용.

### GATE 1

AskUserQuestion으로 버전을 확인받는다:

```
다음 버전으로 릴리즈할까요?

제안 버전: vX.X.X
근거: feat 커밋 N개, fix 커밋 M개 ...

주요 변경사항:
- feat: ...
- fix: ...
```

사용자 승인 없이 다음 단계로 절대 진행하지 않는다.

---

## Stage 2: CHANGELOG 생성

### 2-1. 커밋 분류

마지막 태그 이후 커밋을 Keep a Changelog 카테고리로 분류:

| 커밋 타입 | CHANGELOG 섹션 |
|-----------|---------------|
| `feat:` | Added |
| `fix:` | Fixed |
| `refactor:`, `style:`, `chore:` | Changed |
| `BREAKING` | ⚠️ Breaking Changes (최상단) |

### 2-2. CHANGELOG.md 수정

기존 `CHANGELOG.md`를 Read로 읽고, 최신 엔트리 바로 위에 새 버전 섹션을 삽입한다.

형식 (Keep a Changelog):
```markdown
## [X.X.X] - YYYY-MM-DD

### Added
- 새 기능 설명 (#이슈번호)

### Fixed
- 버그 수정 설명 (#이슈번호)

### Changed
- 변경사항 설명 (#이슈번호)
```

하단 링크 섹션에도 새 버전 비교 링크 추가:
```markdown
[X.X.X]: https://github.com/indigo-net/Brf.it/compare/vPREV...vX.X.X
```

### 2-3. 커밋 & 푸시

```bash
git add CHANGELOG.md
git commit -m "docs: CHANGELOG vX.X.X 추가"
git push origin main
```

### GATE 2

push 성공 + `git status --porcelain` 출력 없음 확인.

---

## Stage 3: 태그 생성 & 푸시

```bash
git tag vX.X.X
git push origin vX.X.X
```

### GATE 3

아래 두 가지 확인:

1. 태그가 리모트에 존재하는지:
```bash
git ls-remote --tags origin | grep vX.X.X
```

2. GitHub Actions 워크플로우가 트리거되었는지:
```bash
gh run list --workflow=release.yml --limit 1
```

트리거 확인되면 Stage 4 진행.

---

## Stage 4: CI 모니터링 (HARD-GATE)

**이 단계는 HARD-GATE이다. CI 실패 시 인라인 수정 절대 금지.**

### 4-1. CI 완료 대기

```bash
gh run watch <run-id> --exit-status
```

### 4-2A. 성공 시

Stage 5로 진행.

### 4-2B. 실패 시 — 즉시 중단 프로토콜

**절대로 코드를 수정하거나 재태그하지 않는다.** 아래 순서를 엄격히 따른다:

1. 실패 로그 수집:
```bash
gh run view <run-id> --log-failed
```

2. GitHub Issue 생성:
```bash
gh issue create \
  --title "fix: vX.X.X 릴리즈 CI 실패" \
  --body "## 개요
vX.X.X 릴리즈 CI가 실패했습니다.

## CI 로그
\`\`\`
(gh run view --log-failed 출력)
\`\`\`

## 환경
- Workflow: release.yml
- Run ID: <run-id>
- Tag: vX.X.X

## 다음 단계
1. 이슈 원인 분석 및 수정
2. 수정 커밋 후 새 태그로 릴리즈 재시도" \
  --assignee indigo-net \
  --label bug
```

3. 사용자에게 알림:
```
❌ CI 실패 — 릴리즈 중단

- 실패 Run: <run-url>
- 생성된 이슈: #XX
- 다음 단계: 이슈 해결 후 새 버전으로 릴리즈 재시도
```

4. **즉시 중단. Stage 5로 절대 진행하지 않는다.**

---

## Stage 5: 릴리즈 노트 & 후처리

### 5-1. 에셋 확인

```bash
gh release view vX.X.X --json assets --jq '.assets[].name'
```

6개 에셋 확인 (`.goreleaser.yml` 기준):
- `brfit_X.X.X_darwin_amd64.tar.gz`
- `brfit_X.X.X_darwin_arm64.tar.gz`
- `brfit_X.X.X_linux_amd64.tar.gz`
- `brfit_X.X.X_linux_arm64.tar.gz`
- `brfit_X.X.X_windows_amd64.zip`
- `checksums.txt`

에셋이 6개가 아니면 사용자에게 경고하고 확인을 받는다.

### 5-2. 릴리즈 노트 작성

`.github/RELEASE_TEMPLATE.md` 형식을 따라 릴리즈 노트를 작성한다.

이전 태그를 확인하여 Full Changelog 링크를 포함:
```bash
git describe --tags --abbrev=0 vX.X.X^  # 이전 태그
```

```bash
gh release edit vX.X.X --notes "$(cat <<'EOF'
## 🎉 [릴리즈 제목]

[1-2문장 설명]

### ✨ New Features

- 기능 설명

### 🔧 Improvements

- 개선사항

### 🐛 Bug Fixes

- 수정사항

### 📦 Installation

```bash
# macOS/Linux
curl -fsSL https://raw.githubusercontent.com/indigo-net/Brf.it/main/install.sh | bash

# or via Go
go install github.com/indigo-net/Brf.it/cmd/brfit@vX.X.X
```

### 📝 Example

```bash
brfit . -f md
```

**Full Changelog**: https://github.com/indigo-net/Brf.it/compare/vPREV...vX.X.X
EOF
)"
```

### 5-3. Homebrew tap 갱신 확인

```bash
gh api repos/indigo-net/homebrew-tap/contents/Formula/brfit.rb --jq '.content' | base64 -d | head -5
```

GoReleaser가 자동 갱신하므로 `version` 필드가 새 버전인지 확인만 한다. 갱신되지 않았으면 사용자에게 알린다.

### 5-4. SAMPLE 갱신 워크플로우 확인

```bash
gh workflow list --all | grep "Update SAMPLE"
gh run list --workflow="update-code-package.yml" --limit 1
```

자동 갱신이 트리거되었는지 확인. 트리거되지 않았으면:
```bash
gh workflow run "Update SAMPLE"
```

### 5-5. 완료 요약

```
✅ vX.X.X 릴리즈 완료

| 항목 | 상태 |
|------|------|
| 프리플라이트 | ✅ PASS |
| CHANGELOG | ✅ 업데이트 완료 |
| 태그 | ✅ vX.X.X |
| CI | ✅ PASS |
| 에셋 | ✅ 6/6 |
| 릴리즈 노트 | ✅ 작성 완료 |
| Homebrew tap | ✅ / ⚠️ |
| SAMPLE 갱신 | ✅ / ⚠️ |

🔗 릴리즈: https://github.com/indigo-net/Brf.it/releases/tag/vX.X.X
```

---

## Common Mistakes

| 실수 | 방지책 |
|------|--------|
| CI 실패 시 인라인 수정 후 재태그 | HARD-GATE: Issue 생성 후 즉시 중단 |
| 테스트 건너뛰고 릴리즈 | Stage 0에서 `go test ./...` 필수 |
| CHANGELOG 작성 안 하고 태그 | Stage 2 완료 후에만 Stage 3 진행 |
| 리모트 동기화 안 하고 태그 | Stage 0-3에서 fetch + HEAD 비교 |
| 에셋 수 확인 안 함 | Stage 5-1에서 6개 에셋 검증 |
| 릴리즈 노트 템플릿 미준수 | `.github/RELEASE_TEMPLATE.md` 참조 필수 |
| 버전을 사용자 확인 없이 결정 | Stage 1 GATE에서 AskUserQuestion 필수 |

## Red Flags — 멈추고 재확인

- "CI 에러가 사소해서 빠르게 고치고 재태그하면 된다"
- "테스트는 로컬에서 이미 돌렸으니 스킵해도 된다"
- "CHANGELOG은 나중에 작성해도 된다"
- "Stage 0 실패했지만 Stage 1으로 넘어가도 될 것 같다"
- "에셋이 5개인데 하나는 빠져도 괜찮겠다"
- "사용자 확인 안 받고 바로 태그해도 될 것 같다"
