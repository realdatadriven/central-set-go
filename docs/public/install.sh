#!/bin/bash -e

# central-set installer

# Downloads and extracts the latest release into the current directory.

REPO="realdatadriven/central-set-go"
API_URL="https://api.github.com/repos/${REPO}/releases/latest"

die() {
echo >&2 "ERROR: $*"
exit 1
}

require() {
for cmd in "$@"; do
command -v "$cmd" >/dev/null 2>&1 || die "Required tool '${cmd}' not found."
done
}

detect_platform() {
OS=$(uname -s)
ARCH=$(uname -m)


case "${OS}" in
    Linux)
        case "${ARCH}" in
            x86_64|amd64) PLATFORM="linux-amd64" ;;
            arm64|aarch64) PLATFORM="linux-arm64" ;;
            *) die "Unsupported Linux architecture: ${ARCH}" ;;
        esac
        ;;
    Darwin)
        case "${ARCH}" in
            x86_64) PLATFORM="macos-amd64" ;;
            arm64)  PLATFORM="macos-arm64" ;;
            *) die "Unsupported macOS architecture: ${ARCH}" ;;
        esac
        ;;
    *)
        die "Unsupported operating system: ${OS}"
        ;;
esac


}

fetch_release() {
echo "Fetching latest release information..."


RELEASE_JSON=$(curl --fail --silent --location \
    --header "Accept: application/vnd.github+json" \
    "${API_URL}") || die "Failed to fetch release information."

if command -v jq >/dev/null 2>&1; then
    VERSION=$(echo "${RELEASE_JSON}" | jq -r '.tag_name')
else
    VERSION=$(echo "${RELEASE_JSON}" \
        | grep -o '"tag_name": *"[^"]*"' \
        | head -1 \
        | cut -d'"' -f4)
fi

[ -n "${VERSION}" ] || die "Could not determine latest version."

ASSET_NAME="central-set-${PLATFORM}.zip"
ASSET_VERSIONED_NAME="central-set-${PLATFORM}-${VERSION}.zip"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ASSET_NAME}"

echo "Version : ${VERSION}"
echo "Asset   : ${ASSET_NAME}"
echo "URL     : ${DOWNLOAD_URL}"


}

install_release() {
TMP_ZIP="/tmp/${ASSET_VERSIONED_NAME}"


echo
echo "Downloading..."

if [ ! -f "${TMP_ZIP}" ]; then
    curl --fail --location --progress-bar "${DOWNLOAD_URL}" -o "${TMP_ZIP}" || die "Failed to download ${DOWNLOAD_URL}"
fi

echo
echo "Extracting into $(pwd)..."

unzip -o "${TMP_ZIP}" -d . || die "Failed to extract archive."

#rm -f "${TMP_ZIP}"

if [ ! -f "c7" ] && [ -f "central-set-${PLATFORM}" ]; then
    mv "central-set-${PLATFORM}" "c7"
    echo "rename central-set-${PLATFORM} to c7"
fi

if [ ! -f ".env" ] && [ -f "dot-env-example.txt" ]; then
    mv "dot-env-example.txt" ".env"
    echo "Created .env from dot-env-example.txt"
fi

if [ -f "c7" ]; then
    chmod +x "c7"
    echo "chmod +x c7"
fi

if [ -f "c7" ]; then
    ./c7 --init --model "admin_model.md"
    echo "Set up admin model"
fi

if [ -f "c7" ]; then
    ./c7 --init --model "etlx_model.md"
    echo "Set up etlx model"
fi

echo
echo "Installation completed."
echo
echo "Files extracted to:"
echo "  $(pwd)"
echo


}

main() {
require curl unzip


detect_platform
fetch_release
install_release


}

main
