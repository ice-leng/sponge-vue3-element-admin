#!/bin/bash

#!/bin/bash
script_path=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
internal_path=$(dirname "$script_path")
output="${internal_path}/internal/database/admin.sql"

mysqlUsername="root"
mysqlDatabase="hyperf"

export MYSQL_PWD=""
mysqldump -u ${mysqlUsername} ${mysqlDatabase} > ${output}
unset MYSQL_PWD

git add ${output}

echo ""
echo "mysqldump successfully, output file = ${output}"
