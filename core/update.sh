#!/bin/bash
newVersion=`wget -q -O - https://raw.githubusercontent.com/cdle/binary/main/compile_time.go | tr -cd "[0-9]"`
oldVersion=`cat /opt/amd64/compile_time.go | tr -cd "[0-9]"`
if [[ ${#newVersion} == 13 && $newVersion != $oldVersion ]]; then
    rm -rf /opt/binary
    cd /root
    git clone https://github.com/cdle/binary.git

    cd /opt/amd64
    git checkout --orphan latest_branch
    rm -rf *
    cp /opt/binary/sillyGirl_linux_amd64_$newVersion /opt/amd64
    cp /opt/binary/compile_time.go /opt/amd64
    git add -A
    git commit -am "commit message"
    git branch -D main
    git branch -m main
    git push -f origin main

    cd /root/arm64
    git checkout --orphan latest_branch
    rm -rf *
    cp /opt/binary/sillyGirl_linux_arm64_$newVersion /root/arm64
    cp /opt/binary/compile_time.go /root/arm64
    git add -A
    git commit -am "commit message"
    git branch -D main
    git branch -m main
    git push -f origin main

    rm -rf /opt/binary
fi

