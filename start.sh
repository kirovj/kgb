git pull
source /etc/profile
go build bin/kgb.go
nohup ./kgb > /dev/null 2>&1 &