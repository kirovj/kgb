ps auwx | grep "kgb" | grep -v "grep" | grep "wyt" | awk '{print $2}' | xargs kill -9
git pull
cd bin
go build kgb.go
nohup ./kgb > /dev/null 2>&1 &