#!/bin/bash

# Test script for petrock generators
# Creates a new project, adds a posts feature, and tests the generator command

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

main() {
    local project_name="generator-test"
    local tmp_dir="tmp"
    local project_path="${tmp_dir}/${project_name}"

    log_info "Starting generator test suite..."

    # Clean up any existing test project
    if [[ -d "$project_path" ]]; then
        log_warn "Removing existing test project at $project_path"
        rm -rf "$project_path"
    fi

    # Ensure tmp directory exists
    mkdir -p "$tmp_dir"

    # Step 1: Create new project
    log_info "Creating new project '$project_name' in $tmp_dir/"
    cd "$tmp_dir"
    ../petrock new "$project_name" "github.com/example/$project_name"

    # Change to project directory
    cd "$project_name"
    log_info "Changed to project directory: $(pwd)"

    # Step 2: Check if posts feature exists, create if not
    if [[ ! -d "internal/posts" ]]; then
        log_info "Posts feature not found, generating posts feature..."
        ../../petrock feature posts
    else
        log_info "Posts feature already exists"
    fi

    # Step 3: Test the component generator commands
    log_info "Testing component generators..."
    
    log_info "Generating command component: posts/publish"
    ../../petrock new command posts/publish
    git add . && git commit -m "Add command component posts/publish"
    
    log_info "Generating query component: posts/search"
    ../../petrock new query posts/search
    git add . && git commit -m "Add query component posts/search"
    
    log_info "Generating worker component: posts/analytics"
    ../../petrock new worker posts/analytics
    git add . && git commit -m "Add worker component posts/analytics"
    
    log_info "Testing collision detection by trying to generate same worker again"
    ../../petrock new worker posts/analytics || log_info "Collision correctly detected!"

    log_info "Generator test completed successfully!"
}

# Run main function
main "$@"
