#!/bin/bash -e

# central-set-go (c7) installer script
# Installs the latest central-set-go release binary as `c7` into ~/.c7/
# Issues/PRs: https://github.com/realdatadriven/central-set-go

REPO="realdatadriven/central-set-go"
BINARY_NAME="c7"
INSTALL_DIR="${HOME}/.c7"
API_URL="https://api.github.com/repos/${REPO}/releases/latest"

# ── helpers ──────────────────────────────────────────────────────────────────

die() { echo >&2 "ERROR: $*"; exit 1; }

require() {
    for cmd in "$@"; do
        command -v "$cmd" >/dev/null 2>&1 || die "Required tool '${cmd}' not found. Please install it and retry."
    done
}

# ── detect OS / arch ─────────────────────────────────────────────────────────

detect_platform() {
    OS=$(uname -s)
    ARCH=$(uname -m)

    case "${OS}" in
        Linux)
            case "${ARCH}" in
                x86_64|amd64) PLATFORM="linux-amd64" ;;
                # arm64/aarch64 builds are not yet published; fail gracefully
                aarch64|arm64) die "Linux arm64 builds are not yet available. Follow https://github.com/${REPO}/releases for updates." ;;
                *) die "Unsupported Linux architecture: ${ARCH}" ;;
            esac
            ;;
        Darwin)
            case "${ARCH}" in
                x86_64)        PLATFORM="macos-amd64" ;;
                arm64)         PLATFORM="macos-arm64" ;;
                *) die "Unsupported macOS architecture: ${ARCH}" ;;
            esac
            ;;
        *)
            die "Unsupported operating system: ${OS}. This installer supports Linux and macOS only."
            ;;
    esac
}

# ── fetch latest release info from GitHub API ─────────────────────────────────

fetch_release() {
    echo "Fetching latest release information from GitHub..."

    RELEASE_JSON=$(curl --fail --silent --location \
        --header "Accept: application/vnd.github+json" \
        "${API_URL}") || die "Failed to fetch release information from ${API_URL}"
    ASSET_NAME="central-set-${PLATFORM}.zip"
    # extract tag name (works with or without jq)
    if command -v jq >/dev/null 2>&1; then
        VERSION=$(echo "${RELEASE_JSON}" | jq -r '.tag_name')
        DOWNLOAD_URL=$(echo "${RELEASE_JSON}" | jq -r \
            --arg asset "central-set-${PLATFORM}.zip" \
            '.assets[] | select(.name == $asset) | .browser_download_url')
    else
        # pure bash/sed fallback — good enough for well-formed GitHub API JSON
        VERSION=$(echo "${RELEASE_JSON}" | grep -o '"tag_name": *"[^"]*"' | head -1 | grep -o '"[^"]*"$' | tr -d '"')
        DOWNLOAD_URL=$(echo "${RELEASE_JSON}" | grep -o '"browser_download_url": *"[^"]*'"${ASSET_NAME}"'[^"]*"' | head -1 | grep -o 'https://[^"]*')
    fi
    if [ -z "${DOWNLOAD_URL}" ] || [ "${DOWNLOAD_URL}" = "null" ]; then
        echo "Asset URL not found in GitHub API response. Falling back to standard release URL..."
        DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ASSET_NAME}"
        echo "URL: ${DOWNLOAD_URL}"
    fi
    [ -n "${VERSION}" ]      || die "Could not determine latest version from GitHub API response."
    [ -n "${DOWNLOAD_URL}" ] || die "No asset found for platform '${PLATFORM}' in release ${VERSION} (${DOWNLOAD_URL}). Check https://github.com/${REPO}/releases for available binaries."
}

# ── banner ────────────────────────────────────────────────────────────────────

print_banner() {
    echo
    echo "╔══════════════════════════════════════════════════════╗"
    echo "║          central-set-go installer — installing as c7           ║"
    echo "╚══════════════════════════════════════════════════════╝"
    echo
    echo "  Version  : ${VERSION}"
    echo "  Platform : ${PLATFORM}"
    echo "  Dest     : ${INSTALL_DIR}/${BINARY_NAME}"
    echo
}

# ── install ───────────────────────────────────────────────────────────────────

install_binary() {
    VERSIONED_DIR="${INSTALL_DIR}/${VERSION}"
    LATEST_LINK="${INSTALL_DIR}/latest"
    DEST_BIN="${VERSIONED_DIR}/${BINARY_NAME}"
    TMP_ZIP="${VERSIONED_DIR}/central-set-${PLATFORM}.zip"

    mkdir -p "${VERSIONED_DIR}" || die "Failed to create directory ${VERSIONED_DIR}"

    if [ -f "${DEST_BIN}" ]; then
        echo "Binary already exists at ${DEST_BIN} — skipping download."
    else
        echo "Downloading central-set ${VERSION} for ${PLATFORM}..."
        curl --fail --location --progress-bar "${DOWNLOAD_URL}" -o "${TMP_ZIP}" || die "Download failed from ${DOWNLOAD_URL}"

        echo "Extracting..."
        # The zip contains a single binary named central-set-<platform>
        unzip -o -j "${TMP_ZIP}" "central-set-${PLATFORM}" -d "${VERSIONED_DIR}"  || die "Failed to extract archive."

        # rename the extracted binary to c7
        mv "${VERSIONED_DIR}/central-set-${PLATFORM}" "${DEST_BIN}" || die "Failed to rename binary to ${BINARY_NAME}."

        chmod a+x "${DEST_BIN}" || die "Failed to set executable permissions."
        rm -f "${TMP_ZIP}"

        [ -f "${DEST_BIN}" ] || die "Binary not found after extraction: ${DEST_BIN}"
        echo "Successfully installed ${BINARY_NAME} ${VERSION} to ${DEST_BIN}"
    fi

    # update the 'latest' symlink
    rm -f "${LATEST_LINK}"
    ln -s "${VERSIONED_DIR}" "${LATEST_LINK}" \
        || die "Failed to create symlink ${LATEST_LINK} → ${VERSIONED_DIR}"
}

# ── PATH setup ────────────────────────────────────────────────────────────────

setup_path() {
    LATEST_BIN="${INSTALL_DIR}/latest"
    LOCAL_BIN="${HOME}/.local/bin"

    echo
    echo "──────────────────────────────────────────────────────"
    echo "  PATH setup"
    echo "──────────────────────────────────────────────────────"

    # Detect current shell profile
    PROFILE=""
    case "${SHELL}" in
        */zsh)  PROFILE="${HOME}/.zshrc" ;;
        */bash)
            if [ -f "${HOME}/.bash_profile" ]; then
                PROFILE="${HOME}/.bash_profile"
            else
                PROFILE="${HOME}/.bashrc"
            fi
            ;;
        */fish) PROFILE="${HOME}/.config/fish/config.fish" ;;
    esac

    PATH_LINE='export PATH="'"${LATEST_BIN}"':$PATH"'

    # Try ~/.local/bin symlink first (already-on-PATH convenience)
    if [ -d "${LOCAL_BIN}" ] && [ -w "${LOCAL_BIN}" ] && [ ! -e "${LOCAL_BIN}/${BINARY_NAME}" ]; then
        ln -s "${LATEST_BIN}/${BINARY_NAME}" "${LOCAL_BIN}/${BINARY_NAME}" \
            && echo "  Created symlink: ${LOCAL_BIN}/${BINARY_NAME} → ${LATEST_BIN}/${BINARY_NAME}" \
            || true
    fi

    echo
    echo "  To make '${BINARY_NAME}' available in every new shell session,"
    echo "  add the following line to your shell profile"
    if [ -n "${PROFILE}" ]; then
        echo "  (${PROFILE}):"
    fi
    echo
    echo "    ${PATH_LINE}"
    echo

    # Offer to append automatically
    if [ -n "${PROFILE}" ] && [ -f "${PROFILE}" ]; then
        if grep -qF "${LATEST_BIN}" "${PROFILE}" 2>/dev/null; then
            echo "  ✓ PATH entry already present in ${PROFILE}"
        else
            printf "  Append automatically? [y/N] "
            read -r REPLY
            case "${REPLY}" in
                [Yy]*)
                    echo "" >> "${PROFILE}"
                    echo "# central-set / c7 — added by install.sh" >> "${PROFILE}"
                    echo "${PATH_LINE}" >> "${PROFILE}"
                    echo "  ✓ Appended to ${PROFILE}"
                    echo "  Run:  source ${PROFILE}  (or open a new terminal)"
                    ;;
                *)
                    echo "  Skipped. Append it manually when ready."
                    ;;
            esac
        fi
    fi

    echo
    echo "  To launch c7 right now:"
    echo "    ${LATEST_BIN}/${BINARY_NAME} --help"
    echo
}

# ── main ──────────────────────────────────────────────────────────────────────

main() {
    require curl unzip

    detect_platform
    fetch_release
    print_banner
    install_binary
    setup_path
}

main
