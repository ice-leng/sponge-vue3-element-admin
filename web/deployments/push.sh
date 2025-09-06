#!/usr/bin/expect

set serviceName "admin"

# parameters
set username [lindex $argv 0]
set password [lindex $argv 1]
set hostname [lindex $argv 2]

set timeout 120

spawn scp -r ./dist.tar.gz ${username}@${hostname}:/tmp/
#expect "*yes/no*"
##send  "yes\r"
expect "*password:*"
send  "${password}\r"
expect eof

spawn ssh ${username}@${hostname}
#expect "*yes/no*"
#send  "yes\r"
expect "*password:*"
send  "${password}\r"

# execute a command or script
expect "*${username}@*"
send "cd /tmp && tar --no-same-owner --no-same-permissions -zxvf dist.tar.gz && find /tmp/dist -name '._*' -delete 2>/dev/null || true\r"
expect "*${username}@*"
send "bash /tmp/dist/deploy.sh\r"

# logging out of a session
expect "*${username}@*"
send "exit\r"

expect eof
