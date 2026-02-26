# Release Notes Template

> 릴리즈 노트 작성 시 아래 템플릿을 복사하여 사용하세요.

---

## GitHub Release Notes

```markdown
## 🎉 [기능명/릴리즈 요약]

[1-2문장으로 이번 릴리즈의 핵심 내용 설명]

### ✨ New Features

**[카테고리명]**
- 기능 1
- 기능 2

**[다른 카테고리명]**
- 기능 3

### 🔧 Improvements

- 개선사항 1
- 개선사항 2

### 🐛 Bug Fixes

- 버그 수정 내용 (있을 경우)

### 📦 Installation

\`\`\`bash
# macOS/Linux
curl -fsSL https://raw.githubusercontent.com/indigo-net/Brf.it/main/install.sh | bash

# or via Go
go install github.com/indigo-net/Brf.it/cmd/brfit@vX.X.X
\`\`\`

### 📝 Example

\`\`\`bash
# 사용 예시
brfit . -f md
\`\`\`

**Full Changelog**: https://github.com/indigo-net/Brf.it/compare/vPREV...vNEW
```

---

## 이모지 규칙

| 섹션 | 이모지 | 사용 시점 |
|------|--------|----------|
| 제목 | 🎉 | 항상 |
| New Features | ✨ | 새 기능이 있을 때 |
| Improvements | 🔧 | 개선사항이 있을 때 |
| Bug Fixes | 🐛 | 버그 수정이 있을 때 |
| Breaking Changes | ⚠️ | 호환성이 깨지는 변경이 있을 때 |
| Deprecated | 🗑️ | 기능이 deprecated될 때 |
| Installation | 📦 | 항상 |
| Example | 📝 | 사용 예시가 있을 때 |
| Documentation | 📚 | 문서 변경이 있을 때 |

---

## 선택적 섹션

필요에 따라 추가할 수 있는 섹션:

```markdown
### ⚠️ Breaking Changes

- 기존 호환성이 깨지는 변경사항

### 🗑️ Deprecated

- 향후 제거 예정인 기능

### 📚 Documentation

- 문서 관련 변경사항
```

---

## CHANGELOG.md Entry

```markdown
## [X.X.X] - YYYY-MM-DD

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
git tag vX.X.X && git push origin main --tags

# 3. GitHub Release 노트 추가
gh release edit vX.X.X --notes "$(cat <<'EOF'
## 🎉 릴리즈 제목

릴리즈 설명...

### ✨ New Features
...

**Full Changelog**: https://github.com/indigo-net/Brf.it/compare/vPREV...vNEW
EOF
)"

# 4. CHANGELOG.md 업데이트 후 푸시
git add CHANGELOG.md && git commit -m "docs: CHANGELOG vX.X.X 추가" && git push
```
