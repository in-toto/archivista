#!/bin/bash
# TCR: Test && Commit || Revert
# Optimized for fast refactoring iterations

set -e

# Test only format-related packages for fast feedback
TEST_CMD="go test ./pkg/metadatastorage/format/... ./pkg/sigstorebundle/... -v"

echo "🔄 Running TCR..."
echo "📦 Testing format packages..."

if $TEST_CMD 2>&1 | tee /tmp/tcr-output.log; then
    git add -A

    # Check if there are actually changes to commit
    if git diff --cached --quiet; then
        echo "⚪ No changes to commit"
        exit 0
    fi

    # Create compact commit message with timestamp
    TIMESTAMP=$(date +%H:%M:%S)
    CHANGED_FILES=$(git diff --cached --name-only | wc -l | tr -d ' ')
    git commit -m "TCR $TIMESTAMP ($CHANGED_FILES files)" --no-verify

    echo "✅ Tests passed - Changes committed ($CHANGED_FILES files)"
    echo "📝 Commit: TCR $TIMESTAMP"
else
    echo ""
    echo "❌ Tests failed - Reverting changes..."
    git restore .
    git clean -fd

    echo ""
    echo "🔍 Last 10 lines of test output:"
    tail -10 /tmp/tcr-output.log

    echo ""
    echo "💡 TIP: Make smaller changes. Try reverting to the last commit and taking baby steps."
fi
