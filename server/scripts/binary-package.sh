#!/bin/bash

serviceName="admin"

mkdir -p ${serviceName}-binary/configs

cp -f deployments/binary/run.sh ${serviceName}-binary
chmod +x ${serviceName}-binary/run.sh

cp -f deployments/binary/deploy.sh ${serviceName}-binary
chmod +x ${serviceName}-binary/deploy.sh

cp -f cmd/${serviceName}/${serviceName} ${serviceName}-binary
cp -f configs/${serviceName}.yml ${serviceName}-binary/configs

# 复制枚举目录到二进制包中
go test -v ./internal/pkg/util -run TestEnumSave
mv enum.json ${serviceName}-binary/

# Clean macOS metadata files and extended attributes
find ${serviceName}-binary -name "._*" -delete 2>/dev/null || true
xattr -cr ${serviceName}-binary 2>/dev/null || true

# Create tar archive using a temporary directory to avoid macOS metadata
TEMP_DIR=$(mktemp -d)
cp -r ${serviceName}-binary ${TEMP_DIR}/

# Ensure clean tar archive
COPYFILE_DISABLE=1 tar --no-xattrs --exclude='._*' --exclude='.DS_Store' -zcvf ${serviceName}-binary.tar.gz -C ${TEMP_DIR} ${serviceName}-binary

# Cleanup
rm -rf ${TEMP_DIR}
rm -rf ${serviceName}-binary

echo ""
echo "package binary successfully, output file = ${serviceName}-binary.tar.gz"
