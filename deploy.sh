#!/bin/bash

ps auwx | grep "kgb" | grep -v "grep" | grep "wyt" | awk '{print $2}' | xargs kill -9

echo "git pull code..."
git pull origin main

echo "start go build..."

go build kgb.go

nohup ./kgb > /dev/null 2>&1 &
echo "run kgb success"