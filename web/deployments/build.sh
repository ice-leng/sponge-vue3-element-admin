#!/bin/bash

pnpm build-only

cp -f deployments/deploy.sh dist
chmod +x dist/deploy.sh

find dist -name "._*" -delete 2>/dev/null || true
xattr -cr dist 2>/dev/null || true

TEMP_DIR=$(mktemp -d)
cp -r dist ${TEMP_DIR}/

# Ensure clean tar archive
COPYFILE_DISABLE=1 tar --no-xattrs --exclude='._*' --exclude='.DS_Store' -zcvf dist.tar.gz -C ${TEMP_DIR} dist

# Cleanup
rm -rf ${TEMP_DIR}
rm -rf dist

echo ""
echo "build successfully, output file = dist.tar.gz"
