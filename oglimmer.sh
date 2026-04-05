#!/usr/bin/env bash

set -euo pipefail

# Define script metadata
SCRIPT_NAME=$(basename "$0")
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# Default configuration
DEFAULT_REGISTRIES=("registry.oglimmer.com")
DEFAULT_FRONTEND_DEPLOYMENT="coffee-diary-frontend"
DEFAULT_BACKEND_DEPLOYMENT="coffee-diary-backend"

# Configuration variables (can be overridden by parameters)
REGISTRIES=("${DEFAULT_REGISTRIES[@]}")
FRONTEND_IMAGES=()
BACKEND_IMAGES=()
FRONTEND_DEPLOYMENT="$DEFAULT_FRONTEND_DEPLOYMENT"
BACKEND_DEPLOYMENT="$DEFAULT_BACKEND_DEPLOYMENT"

# Directories
BACKEND_DIR="$SCRIPT_DIR/backend"
FRONTEND_DIR="$SCRIPT_DIR/frontend"

# Default options (can be overridden by environment variables)
BUILD_FRONTEND="${BUILD_FRONTEND:-false}"
BUILD_BACKEND="${BUILD_BACKEND:-false}"
VERBOSE="${VERBOSE:-false}"
DRY_RUN="${DRY_RUN:-false}"
RESTART="${RESTART:-true}"
PUSH="${PUSH:-true}"
HELP=false
PLATFORM="${PLATFORM:-arm64}"
RELEASE_MODE=false
SHOW_VERSIONS=false
E2E_MODE=false
COPY_DB_MODE=false

# Color output (only if terminal supports it)
if [[ -t 1 ]] && command -v tput >/dev/null 2>&1; then
  BOLD="$(tput bold)"
  GREEN="$(tput setaf 2)"
  YELLOW="$(tput setaf 3)"
  RED="$(tput setaf 1)"
  BLUE="$(tput setaf 4)"
  RESET="$(tput sgr0)"
else
  BOLD="" GREEN="" YELLOW="" RED="" BLUE="" RESET=""
fi

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${RESET} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${RESET} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${RESET} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${RESET} $1" >&2
}

# Verbose logging
log_verbose() {
    if [[ "$VERBOSE" == true ]]; then
        echo -e "${BLUE}[VERBOSE]${RESET} $1"
    fi
}

# Execute command with dry-run and verbose support
execute_cmd() {
    local cmd="$1"

    if [[ "$DRY_RUN" == true ]]; then
        echo -e "${YELLOW}[DRY-RUN]${RESET} ${cmd}"
        return 0
    else
        log_verbose "Executing: $cmd"
        if [[ "$VERBOSE" == true ]]; then
            eval "$cmd"
        else
            eval "$cmd" >/dev/null 2>&1
        fi
    fi
}

# Show usage information
show_help() {
    cat << EOF
Usage: ${SCRIPT_NAME} [OPTIONS] [COMMAND]

Build, deploy, and release Coffee Diary application components.

COMMANDS:
    build               Build and deploy components (default)
    release             Create a new release with version bumping and build
    show                Show current backend and frontend versions
    e2e                 Run Playwright e2e tests (starts DB + backend, runs tests, tears down)
    copy-db             Copy production database to local Docker Compose MariaDB

BUILD OPTIONS:
    -f, --frontend          Build and deploy frontend only
    -b, --backend           Build and deploy backend only
    -a, --all               Build and deploy both frontend and backend (default if no component specified)
    -v, --verbose           Enable verbose output
    -n, --no-restart        Skip Kubernetes deployment restart
    --no-push               Skip pushing images to registry
    --dry-run               Show what would be done without executing

    # Registry configuration options
    --registries "REG1,REG2"    Comma-separated list of registries to push to (default: ${DEFAULT_REGISTRIES[0]})
    --frontend-deploy NAME      Frontend deployment name (default: $DEFAULT_FRONTEND_DEPLOYMENT)
    --backend-deploy NAME       Backend deployment name (default: $DEFAULT_BACKEND_DEPLOYMENT)

    # Platform options
    --platform PLATFORM        Target platform(s) for Docker build:
                               - amd64: Build for AMD64/x86_64 architecture
                               - arm64: Build for ARM64 architecture
                               - multi: Build for both amd64 and arm64 (multi-platform)
                               - auto: Detect current platform automatically

    -h, --help              Show this help message

EXAMPLES:
    ${SCRIPT_NAME} build                                    # Build and deploy both components with defaults
    ${SCRIPT_NAME} build -f                                 # Build and deploy frontend only
    ${SCRIPT_NAME} build -b -v                              # Build and deploy backend with verbose output
    ${SCRIPT_NAME} release                                  # Create a new release with version bump and build
    ${SCRIPT_NAME} show                                     # Show current versions
    ${SCRIPT_NAME} e2e                                      # Run e2e tests with fresh DB
    ${SCRIPT_NAME} copy-db                                   # Copy prod DB to local (requires PROD_DB_PASSWORD)
    ${SCRIPT_NAME} build --registries my-registry.com       # Use custom registry
    ${SCRIPT_NAME} build --platform amd64                   # Build for AMD64 only

ENVIRONMENT VARIABLES:
    FRONTEND_DEPLOYMENT     Override default frontend deployment name
    BACKEND_DEPLOYMENT      Override default backend deployment name
    PLATFORM                Override default platform (amd64|arm64|multi|auto)
    DEFAULT_REGISTRIES_ENV  Override default registries (comma-separated)
    VERBOSE                 Enable verbose mode (true/false)
    DRY_RUN                 Enable dry-run mode (true/false)
    PUSH                    Enable/disable pushing to registry (true/false)
    RESTART                 Enable/disable Kubernetes restart (true/false)
    PROD_DB_PASSWORD        Production MariaDB root password (required for copy-db)

EOF
}

# Parse command line arguments
parse_args() {
    # Check if first argument is a command
    if [[ $# -gt 0 ]]; then
        case $1 in
            build)
                shift
                ;;
            release)
                RELEASE_MODE=true
                shift
                ;;
            show)
                SHOW_VERSIONS=true
                shift
                ;;
            e2e)
                E2E_MODE=true
                shift
                ;;
            copy-db)
                COPY_DB_MODE=true
                shift
                ;;
            help|-h|--help)
                HELP=true
                shift
                ;;
        esac
    fi

    while [[ $# -gt 0 ]]; do
        case $1 in
            -f|--frontend)
                BUILD_FRONTEND=true
                shift
                ;;
            -b|--backend)
                BUILD_BACKEND=true
                shift
                ;;
            -a|--all)
                BUILD_FRONTEND=true
                BUILD_BACKEND=true
                shift
                ;;
            -v|--verbose)
                VERBOSE=true
                shift
                ;;
            -n|--no-restart)
                RESTART=false
                shift
                ;;
            --no-push)
                PUSH=false
                shift
                ;;
            --dry-run)
                DRY_RUN=true
                shift
                ;;
            --registries)
                # Clear existing registries and parse comma-separated list
                REGISTRIES=()
                IFS=',' read -ra ADDR <<< "$2"
                for registry in "${ADDR[@]}"; do
                    REGISTRIES+=("$(echo "$registry" | xargs)")  # trim whitespace
                done
                shift 2
                ;;
            --frontend-deploy)
                FRONTEND_DEPLOYMENT="$2"
                shift 2
                ;;
            --backend-deploy)
                BACKEND_DEPLOYMENT="$2"
                shift 2
                ;;
            --platform)
                PLATFORM="$2"
                shift 2
                ;;
            -h|--help)
                HELP=true
                shift
                ;;
            *)
                log_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done

    # Handle environment variable overrides
    FRONTEND_DEPLOYMENT="${FRONTEND_DEPLOYMENT:-$FRONTEND_DEPLOYMENT}"
    BACKEND_DEPLOYMENT="${BACKEND_DEPLOYMENT:-$BACKEND_DEPLOYMENT}"
    PLATFORM="${DOCKER_PLATFORM:-$PLATFORM}"

    # Override default registries from environment if set
    if [[ -n "${DEFAULT_REGISTRIES_ENV:-}" ]]; then
        REGISTRIES=()
        IFS=',' read -ra ADDR <<< "$DEFAULT_REGISTRIES_ENV"
        for registry in "${ADDR[@]}"; do
            REGISTRIES+=("$(echo "$registry" | xargs)")
        done
    fi

    # Build image arrays from registries
    if [[ ${#REGISTRIES[@]} -gt 0 ]]; then
        FRONTEND_IMAGES=()
        BACKEND_IMAGES=()
        for registry in "${REGISTRIES[@]}"; do
            FRONTEND_IMAGES+=("$registry/coffee-diary-frontend")
            BACKEND_IMAGES+=("$registry/coffee-diary-backend")
        done
    else
        # Fallback to defaults if no registries specified
        FRONTEND_IMAGES=("${DEFAULT_REGISTRIES[0]}/coffee-diary-frontend")
        BACKEND_IMAGES=("${DEFAULT_REGISTRIES[0]}/coffee-diary-backend")
    fi

    # Validate platform parameter
    if [[ -n "$PLATFORM" && ! "$PLATFORM" =~ ^(amd64|arm64|multi|auto)$ ]]; then
        log_error "Invalid platform: $PLATFORM. Must be one of: amd64, arm64, multi, auto"
        exit 1
    fi

    # Validate conflicting options
    if [[ "$PUSH" == false && "$RESTART" == true && "$RELEASE_MODE" == false ]]; then
        log_warning "Cannot restart deployments without pushing images. Setting --no-restart."
        RESTART=false
    fi

    # If no component specified for build mode, build both
    if [[ "$RELEASE_MODE" == false && "$SHOW_VERSIONS" == false && "$COPY_DB_MODE" == false && "$BUILD_FRONTEND" == false && "$BUILD_BACKEND" == false ]]; then
        BUILD_FRONTEND=true
        BUILD_BACKEND=true
    fi
}

# Check if required tools are available
check_prerequisites() {
    local tools=("docker" "kubectl")

    # Add additional tools for release mode
    if [[ "$RELEASE_MODE" == true ]]; then
        tools+=("npm" "git")
    fi

    local missing_deps=()
    for tool in "${tools[@]}"; do
        if ! command -v "$tool" >/dev/null 2>&1; then
            missing_deps+=("$tool")
        fi
    done

    if [[ ${#missing_deps[@]} -gt 0 ]]; then
        log_error "Missing required dependencies: ${missing_deps[*]}"
        echo "Please install the missing dependencies and try again." >&2
        exit 1
    fi

    # Check if Docker daemon is running (skip in dry-run mode)
    if [[ "$DRY_RUN" != true ]] && ! docker info >/dev/null 2>&1; then
        log_error "Docker daemon is not running"
        echo "Please start Docker and try again." >&2
        exit 1
    fi

    # Check if buildx is available for multi-platform builds
    if [[ "$PLATFORM" == "multi" ]]; then
        if ! docker buildx version &> /dev/null; then
            log_error "Docker buildx is required for multi-platform builds but not available"
            log_info "Please install Docker Desktop or enable buildx plugin"
            exit 1
        fi

        # Ensure buildx builder is available
        if ! docker buildx inspect &> /dev/null; then
            log_info "Creating buildx builder instance..."
            docker buildx create --use --name multiplatform-builder 2>/dev/null || true
        fi
    fi

    log_verbose "All required tools are available"
}

# Show current versions
show_versions() {
    # Backend version (from VERSION file or default)
    local backend_version
    if [[ -f "$BACKEND_DIR/VERSION" ]]; then
        backend_version=$(cat "$BACKEND_DIR/VERSION")
    else
        backend_version="0.0.1-SNAPSHOT"
    fi

    # Frontend version
    local frontend_version
    frontend_version=$(grep '"version"' "$FRONTEND_DIR/package.json" | head -1 | sed -E 's/.*"version": *"([^"]+)".*/\1/')

    echo "Backend version: $backend_version"
    echo "Frontend version: $frontend_version"
}

# Bump semantic version
bump_version() {
    local current_version="$1"
    local bump_type="$2"
    # Strip -SNAPSHOT suffix if present
    current_version="${current_version%-SNAPSHOT}"
    IFS='.' read -r major minor patch <<< "$current_version"

    case "$bump_type" in
        major)
            major=$((major + 1)); minor=0; patch=0;
            ;;
        minor)
            minor=$((minor + 1)); patch=0;
            ;;
        bugfix|patch)
            patch=$((patch + 1));
            ;;
        *)
            echo "Unknown bump type: $bump_type" >&2
            exit 1
            ;;
    esac
    echo "$major.$minor.$patch"
}

# Get platform arguments for docker build
get_platform_args() {
    local platform_args=""

    case "$PLATFORM" in
        "amd64")
            platform_args="--platform linux/amd64"
            ;;
        "arm64")
            platform_args="--platform linux/arm64"
            ;;
        "multi")
            platform_args="--platform linux/amd64,linux/arm64"
            ;;
        "auto"|"")
            # Let Docker detect the platform automatically
            platform_args=""
            ;;
    esac

    echo "$platform_args"
}

# Build Docker image for multiple targets
build_image() {
    local component="$1"
    local dockerfile_args="$2"
    local platform_args=$(get_platform_args)

    # Create array of image tags - passed as remaining arguments
    shift 2
    local image_tags=("$@")
    local primary_tag="${image_tags[0]}"

    log_info "Building $component image for ${#image_tags[@]} target(s):"
    for tag in "${image_tags[@]}"; do
        log_info "  - $tag"
    done
    if [[ -n "$platform_args" ]]; then
        log_info "Target platform(s): $PLATFORM"
    fi

    local build_cmd=""

    # Use buildx for multi-platform builds or when platform is specified
    if [[ "$PLATFORM" == "multi" || (-n "$PLATFORM" && "$PLATFORM" != "auto") ]]; then
        build_cmd="docker buildx build $platform_args"

        # Add all tags
        for tag in "${image_tags[@]}"; do
            build_cmd="$build_cmd --tag $tag"
        done

        if [[ "$PUSH" == true ]]; then
            build_cmd="$build_cmd --push"
        else
            # For local builds with buildx, we need to load the image
            if [[ "$PLATFORM" != "multi" ]]; then
                build_cmd="$build_cmd --load"
            else
                log_warning "Multi-platform builds cannot be loaded locally, forcing push to registry"
                build_cmd="$build_cmd --push"
            fi
        fi

        # Add dockerfile arguments
        build_cmd="$build_cmd $dockerfile_args"

    else
        # Use regular docker build for single platform or auto-detection
        # Build with primary tag first
        build_cmd="docker build $platform_args --tag $primary_tag $dockerfile_args"

        # Tag for additional registries
        if [[ ${#image_tags[@]} -gt 1 ]]; then
            for tag in "${image_tags[@]:1}"; do
                build_cmd="$build_cmd && docker tag $primary_tag $tag"
            done
        fi

        # Push to all registries if requested
        if [[ "$PUSH" == true ]]; then
            for tag in "${image_tags[@]}"; do
                build_cmd="$build_cmd && docker push $tag"
            done
        fi
    fi

    log_verbose "Build command: $build_cmd"

    if execute_cmd "$build_cmd"; then
        log_success "$component image built successfully"
        if [[ "$PUSH" == false && "$PLATFORM" != "multi" ]]; then
            log_info "$component image tagged locally (not pushed)"
        elif [[ "$PUSH" == true ]]; then
            log_success "$component image pushed to ${#image_tags[@]} target(s)"
        fi
    else
        log_error "Failed to build $component image"
        exit 1
    fi
}

# Restart Kubernetes deployment
restart_deployment() {
    local deployment="$1"

    log_info "Restarting deployment: $deployment"

    if execute_cmd "kubectl rollout restart deployment/$deployment"; then
        log_success "Deployment $deployment restarted successfully"

        # Wait for rollout to complete if verbose
        if [[ "$VERBOSE" == true ]]; then
            log_info "Waiting for rollout to complete..."
            kubectl rollout status deployment/"$deployment" --timeout=300s
        fi
    else
        log_error "Failed to restart deployment: $deployment"
        exit 1
    fi
}

# Execute build process
execute_build() {
    # Display configuration
    echo -e "${BOLD}=== Build Configuration ===${RESET}"
    echo "Registries:        ${REGISTRIES[*]}"
    echo "Platform:          ${PLATFORM:-auto}"
    echo "Build Frontend:    $BUILD_FRONTEND"
    echo "Build Backend:     $BUILD_BACKEND"
    echo "Push to Registry:  $PUSH"
    echo "Restart K8s:       $RESTART"
    echo "Dry-run:           $DRY_RUN"
    echo "Verbose:           $VERBOSE"
    if [[ "$BUILD_FRONTEND" == true ]]; then
        echo "Frontend Deploy:   $FRONTEND_DEPLOYMENT"
    fi
    if [[ "$BUILD_BACKEND" == true ]]; then
        echo "Backend Deploy:    $BACKEND_DEPLOYMENT"
    fi
    echo -e "${BOLD}===========================${RESET}"
    echo

    log_info "Starting build process..."

    # Build frontend
    if [[ "$BUILD_FRONTEND" == true ]]; then
        build_image "frontend" "--build-arg GIT_COMMIT=$(git rev-parse --short HEAD) frontend/" "${FRONTEND_IMAGES[@]}"
    fi

    # Build backend
    if [[ "$BUILD_BACKEND" == true ]]; then
        build_image "backend" "--build-arg BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ) --build-arg GIT_COMMIT=$(git rev-parse --short HEAD) backend/" "${BACKEND_IMAGES[@]}"
    fi

    # Restart deployments if requested
    if [[ "$RESTART" == true ]]; then
        if [[ "$BUILD_FRONTEND" == true ]]; then
            restart_deployment "$FRONTEND_DEPLOYMENT"
        fi

        if [[ "$BUILD_BACKEND" == true ]]; then
            restart_deployment "$BACKEND_DEPLOYMENT"
        fi
    else
        log_info "Skipping deployment restarts (--no-restart specified)"
    fi

    echo
    echo -e "${BOLD}${GREEN}All operations completed successfully${RESET}"
}

# Execute release process
execute_release() {
    log_info "Starting release process..."

    # Show current versions
    echo "Current versions:"; show_versions; echo

    # Explain bump types
    echo "Select which part to bump (semantic versioning):"
    echo "  1) major  - incompatible API changes"
    echo "  2) minor  - backwards-compatible new features"
    echo "  3) bugfix - backwards-compatible bug fixes"
    PS3="Enter choice (1-3): "
    select bump in major minor bugfix; do
        if [[ -n "$bump" ]]; then
            echo "Chosen bump type: $bump"; break
        else
            echo "Invalid choice. Please select 1, 2, or 3.";
        fi
    done

    # Compute new version
    if [[ -f "$BACKEND_DIR/VERSION" ]]; then
        current_version=$(cat "$BACKEND_DIR/VERSION")
    else
        current_version="0.0.1-SNAPSHOT"
    fi
    new_version=$(bump_version "$current_version" "$bump")
    log_info "Releasing version $new_version..."

    # Update backend version file
    log_info "Updating backend version to $new_version..."
    echo "$new_version" > "$BACKEND_DIR/VERSION"

    # Update frontend
    log_info "Updating frontend version to $new_version..."
    (cd "$FRONTEND_DIR" && npm version "$new_version" --no-git-tag-version)

    # Commit and tag release
    log_info "Committing version changes and creating tag..."
    git add "$BACKEND_DIR/VERSION" "$FRONTEND_DIR/package.json" "$FRONTEND_DIR/package-lock.json"
    git commit -m "Release v$new_version"
    git tag -a "v$new_version" -m "Release v$new_version"

    # Build and upload after version commit
    log_info "Building and uploading release version $new_version..."
    BUILD_FRONTEND=true
    BUILD_BACKEND=true
    execute_build

    # Bump backend to SNAPSHOT
    log_info "Setting backend to SNAPSHOT version..."
    snapshot="${new_version}-SNAPSHOT"
    echo "$snapshot" > "$BACKEND_DIR/VERSION"
    git add "$BACKEND_DIR/VERSION"
    git commit -m "Set backend to $snapshot"

    log_success "Release v$new_version complete. Backend is now $snapshot."
}

# Execute e2e test process
execute_e2e() {
    local db_container="coffee-diary-e2e-db"
    local db_port="3307"
    local backend_pid=""
    local exit_code=0

    # Cleanup function — always runs on exit
    cleanup_e2e() {
        log_info "Tearing down e2e environment..."
        if [[ -n "$backend_pid" ]] && kill -0 "$backend_pid" 2>/dev/null; then
            kill "$backend_pid" 2>/dev/null
            wait "$backend_pid" 2>/dev/null || true
            log_info "Backend stopped"
        fi
        if docker inspect "$db_container" &>/dev/null; then
            docker rm -f "$db_container" >/dev/null 2>&1
            log_info "Database container removed"
        fi
    }
    trap cleanup_e2e EXIT

    # 1. Start fresh MariaDB container
    log_info "Starting fresh MariaDB on port $db_port..."
    docker rm -f "$db_container" 2>/dev/null || true
    docker run -d \
        --name "$db_container" \
        -p "$db_port:3306" \
        -e MARIADB_ROOT_PASSWORD=root \
        -e MARIADB_DATABASE=coffeediary_e2e \
        -e MARIADB_USER=app \
        -e MARIADB_PASSWORD=app \
        mariadb:11 >/dev/null

    # 2. Wait for MariaDB to be ready
    log_info "Waiting for MariaDB to be ready..."
    local retries=30
    while ! docker exec "$db_container" mariadb -uroot -proot -e "SELECT 1" &>/dev/null; do
        retries=$((retries - 1))
        if [[ $retries -le 0 ]]; then
            log_error "MariaDB failed to start"
            exit 1
        fi
        sleep 1
    done
    log_success "MariaDB is ready"

    # 3. Start backend with test config
    log_info "Starting backend..."
    (
        cd "$BACKEND_DIR"
        DB_PORT="$db_port" \
        DB_NAME="coffeediary_e2e" \
        OIDC_ISSUER_URL="https://id.oglimmer.de/realms/playwright-tests" \
        OIDC_CLIENT_ID="test" \
        OIDC_CLIENT_SECRET="3VEZSovF5lkiurjJFsDgW61JCkd8UTdY" \
        go run ./cmd/server
    ) &
    backend_pid=$!

    # 4. Wait for backend health check
    log_info "Waiting for backend to be ready..."
    retries=30
    while ! curl -sf http://localhost:8080/actuator/health >/dev/null 2>&1; do
        if ! kill -0 "$backend_pid" 2>/dev/null; then
            log_error "Backend process died"
            exit 1
        fi
        retries=$((retries - 1))
        if [[ $retries -le 0 ]]; then
            log_error "Backend failed to start"
            exit 1
        fi
        sleep 1
    done
    log_success "Backend is ready"

    # 5. Run Playwright tests
    log_info "Running Playwright e2e tests..."
    (cd "$FRONTEND_DIR" && npx playwright test) || exit_code=$?

    if [[ $exit_code -eq 0 ]]; then
        log_success "All e2e tests passed"
    else
        log_error "E2e tests failed (exit code: $exit_code)"
    fi

    # Cleanup happens via trap
    exit "$exit_code"
}

# Execute copy-db process: dump prod DB and import into local Docker Compose MariaDB
execute_copy_db() {
    local prod_db_password="${PROD_DB_PASSWORD:-}"
    local prod_db_name="coffeediary"
    local local_container="coffee-diary-db-1"
    local local_db_name="coffeediary"
    local dump_file

    if [[ -z "$prod_db_password" ]]; then
        log_error "PROD_DB_PASSWORD environment variable is required"
        echo "Usage: PROD_DB_PASSWORD=<password> ${SCRIPT_NAME} copy-db" >&2
        exit 1
    fi

    # Verify local DB container is running
    if ! docker inspect "$local_container" &>/dev/null; then
        log_error "Local DB container '$local_container' is not running. Start it with: docker compose up -d db"
        exit 1
    fi

    dump_file=$(mktemp /tmp/coffee-diary-prod-dump.XXXXXX.sql)
    trap "rm -f '$dump_file'" EXIT

    # Dump production database via kubectl
    log_info "Dumping production database..."
    if ! kubectl exec -i mariadb-0 -- mariadb-dump -u root -p"$prod_db_password" --single-transaction "$prod_db_name" > "$dump_file"; then
        log_error "Failed to dump production database"
        exit 1
    fi

    local dump_size
    dump_size=$(wc -c < "$dump_file" | xargs)
    log_success "Production dump complete (${dump_size} bytes)"

    # Import into local container
    log_info "Importing into local database..."
    if ! docker exec -i "$local_container" mariadb -u root -proot "$local_db_name" < "$dump_file"; then
        log_error "Failed to import dump into local database"
        exit 1
    fi

    log_success "Production database copied to local successfully"
}

# Main execution function
main() {
    # Show help if no arguments provided
    if [[ $# -eq 0 ]]; then
        show_help
        exit 0
    fi

    parse_args "$@"

    if [[ "$HELP" == true ]]; then
        show_help
        exit 0
    fi

    if [[ "$SHOW_VERSIONS" == true ]]; then
        show_versions
        exit 0
    fi

    if [[ "$E2E_MODE" == true ]]; then
        execute_e2e
        exit 0
    fi

    if [[ "$COPY_DB_MODE" == true ]]; then
        execute_copy_db
        exit 0
    fi

    check_prerequisites

    if [[ "$RELEASE_MODE" == true ]]; then
        execute_release
    else
        execute_build
    fi
}

# Run main function with all arguments
main "$@"
