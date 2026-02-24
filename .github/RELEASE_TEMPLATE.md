# Release Notes Template

> 릴리즈 노트 작성 시 아래 템플릿을 복사하여 사용하세요.

---

## GitHub Release Notes

```markdown
## What's New

### Added
- 새로 추가된 기능

### Changed
- 기존 기능의 변경사항

### Fixed
- 버그 수정

### Removed
- 제거된 기능

## Usage

\`\`\`bash
# 예시 명령어
brfit . -f md
\`\`\`

**Full Changelog**: https://github.com/indigo-net/Brf.it/compare/v0.0.0...v0.0.0
```

---

## CHANGELOG.md Entry

```markdown
## [0.0.0] - YYYY-MM-DD

### Added
- 새로 추가된 기능

### Changed
- 기존 기능의 변경사항

### Fixed
- 버그 수정

### Removed
- 제거된 기능
```

---

## 릴리즈 명령어

```bash
# 1. 커밋
git add . && git commit -m "feat: 변경 내용"

# 2. 태그 생성 및 푸시 (GitHub Actions 트리거)
git tag v0.0.0 && git push origin main --tags

# 3. GitHub Release 노트 추가
gh release edit v0.0.0 --notes "릴리즈 노트 내용"

# 4. CHANGELOG.md 업데이트 후 푸시
git add CHANGELOG.md && git commit -m "docs: CHANGELOG v0.0.0 추가" && git push
```
