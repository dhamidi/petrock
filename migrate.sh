#!/usr/bin/env bash

# Exit immediately if a command exits with a non-zero status.
set -e

# Define placeholder values
PROJECT_NAME_PLACEHOLDER="petrock_example_project_name"
MODULE_PATH_PLACEHOLDER="github.com/petrock/example_module_path"
TEMPLATE_DIR="internal/template/templates/new"
SKELETON_DIR="internal/skeleton"
SKELETON_CMD_DIR="${SKELETON_DIR}/cmd/${PROJECT_NAME_PLACEHOLDER}"
SKELETON_CORE_DIR="${SKELETON_DIR}/core"

echo "Starting migration from templates to skeleton..."

# --- Step 1: Create Skeleton Directory Structure ---
echo "Creating skeleton directories..."
mkdir -p "${SKELETON_CMD_DIR}"
mkdir -p "${SKELETON_CORE_DIR}"
echo "Skeleton directories created."

# --- Step 2: Populate Example Project Skeleton ---

# Move and rename top-level files
echo "Moving and renaming top-level files (.gitignore, go.mod)..."
mv "${TEMPLATE_DIR}/.gitignore.tmpl" "${SKELETON_DIR}/.gitignore"
mv "${TEMPLATE_DIR}/go.mod.tmpl" "${SKELETON_DIR}/go.mod"

# Move and rename core files
echo "Moving and renaming core files..."
for f in "${TEMPLATE_DIR}/core"/*.go.tmpl; do
  if [ -f "$f" ]; then
    target_name=$(basename "$f" .tmpl)
    mv "$f" "${SKELETON_CORE_DIR}/${target_name}"
    echo "  Moved ${f} -> ${SKELETON_CORE_DIR}/${target_name}"
  fi
done

# Move and rename cmd files
echo "Moving and renaming cmd files..."
for f in "${TEMPLATE_DIR}/cmd"/*.go.tmpl; do
  if [ -f "$f" ]; then
    target_name=$(basename "$f" .tmpl)
    mv "$f" "${SKELETON_CMD_DIR}/${target_name}"
    echo "  Moved ${f} -> ${SKELETON_CMD_DIR}/${target_name}"
  fi
done

# Replace placeholders in all skeleton files
echo "Replacing placeholders in skeleton files..."
# Note: Using sed -i '' for macOS compatibility. Linux users might just use sed -i.
# Use pipe delimiter for module path replacement due to slashes.
find "${SKELETON_DIR}" -type f -print0 | while IFS= read -r -d $'\0' file; do
  echo "  Processing ${file}..."
  sed -i '' "s/{{ \.ProjectName }}/${PROJECT_NAME_PLACEHOLDER}/g" "$file"
  sed -i '' "s|{{ \.ModuleName }}|${MODULE_PATH_PLACEHOLDER}|g" "$file"
done

echo "Placeholder replacement complete."

echo "Migration script finished successfully!"
echo "Please review the changes in the '${SKELETON_DIR}' directory."
echo "You may want to run 'go mod tidy' and 'go build ./...' within '${SKELETON_DIR}' to
verify."
