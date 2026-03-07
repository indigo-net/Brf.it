#!/bin/bash
set -e

# Script to determine the next version based on merged PRs and their linked issues
# Outputs:
#   - new_version: The new version number (e.g., 0.5.0)
#   - bump_type: The type of version bump (major, minor, patch)

REPO_OWNER="indigo-net"
REPO_NAME="Brf.it"

# Get current version from latest tag
get_current_version() {
    local tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
    if [ -n "$tag" ]; then
        echo "${tag#v}"
    else
        echo "0.0.0"
    fi
}

# Parse semver version
parse_version() {
    local version=$1
    echo "$version" | tr '.' ' '
}

# Calculate new version
bump_version() {
    local current=$1
    local bump_type=$2

    read major minor patch <<< $(parse_version "$current")

    case $bump_type in
        major)
            major=$((major + 1))
            minor=0
            patch=0
            ;;
        minor)
            minor=$((minor + 1))
            patch=0
            ;;
        patch)
            patch=$((patch + 1))
            ;;
    esac

    echo "${major}.${minor}.${patch}"
}

# Get linked issues for a PR using GraphQL
get_linked_issues() {
    local pr_number=$1

    gh api graphql -f query='
        query($owner: String!, $repo: String!, $number: Int!) {
            repository(owner: $owner, name: $repo) {
                pullRequest(number: $number) {
                    closingIssuesReferences(first: 20) {
                        nodes {
                            number
                            labels(first: 20) {
                                nodes { name }
                            }
                        }
                    }
                }
            }
        }
    ' -f owner="$REPO_OWNER" -f repo="$REPO_NAME" -F number="$pr_number" \
        --jq '.data.repository.pullRequest.closingIssuesReferences.nodes'
}

# Get version priority from issue labels
get_version_priority() {
    local labels=$1
    local priority=1  # Default: patch

    # Check for version labels
    if echo "$labels" | grep -q "version:major"; then
        echo 3
        return
    elif echo "$labels" | grep -q "version:minor"; then
        echo 2
        return
    fi

    echo 1  # patch
}

# Get bump type from priority
get_bump_type() {
    local priority=$1

    case $priority in
        3) echo "major" ;;
        2) echo "minor" ;;
        *) echo "patch" ;;
    esac
}

# Main logic
main() {
    local current_version=$(get_current_version)
    echo "Current version: $current_version"

    local highest_priority=1  # Default: patch
    local latest_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "")

    # Get merged PRs since last tag (or all if no tag)
    local search_query="is:merged base:main"
    if [ -n "$latest_tag" ]; then
        search_query="$search_query merged:>${latest_tag#v}"
    fi

    # Get merged PRs
    local prs=$(gh pr list --state merged --base main \
        --search "${search_query#is:merged base:main }" \
        --json number \
        --jq '.[].number')

    if [ -z "$prs" ]; then
        echo "No merged PRs found"
        echo "new_version=$(bump_version $current_version patch)" >> $GITHUB_OUTPUT
        echo "bump_type=patch" >> $GITHUB_OUTPUT
        return
    fi

    echo "Analyzing merged PRs..."

    # Check each PR for linked issues
    for pr_number in $prs; do
        echo "Checking PR #$pr_number..."

        local issues=$(get_linked_issues "$pr_number")

        if [ -z "$issues" ] || [ "$issues" = "null" ] || [ "$issues" = "[]" ]; then
            echo "  No linked issues, defaulting to patch"
            continue
        fi

        # Check each linked issue
        echo "$issues" | jq -c '.[]' 2>/dev/null | while read issue; do
            if [ -z "$issue" ] || [ "$issue" = "null" ]; then
                continue
            fi

            local issue_number=$(echo "$issue" | jq -r '.number')
            local labels=$(echo "$issue" | jq -r '.labels.nodes | .[].name' | tr '\n' ' ')

            echo "  Issue #$issue_number labels: $labels"

            local priority=$(get_version_priority "$labels")
            echo "  Priority: $priority"

            if [ "$priority" -gt "$highest_priority" ]; then
                highest_priority=$priority
            fi
        done
    done

    # Determine bump type
    local bump_type=$(get_bump_type $highest_priority)
    local new_version=$(bump_version "$current_version" "$bump_type")

    echo "Bump type: $bump_type"
    echo "New version: $new_version"

    # Set outputs
    echo "new_version=$new_version" >> $GITHUB_OUTPUT
    echo "bump_type=$bump_type" >> $GITHUB_OUTPUT
}

main
