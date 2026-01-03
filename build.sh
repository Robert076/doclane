#!/bin/bash

cd backend/cmd
go run main.go &

BACKEND_PID=$!

cd ../../frontend
npm run dev

kill $BACKEND_PID
