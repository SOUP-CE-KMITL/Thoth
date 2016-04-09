#!/bin/bash
export LD_LIBRARY_PATH=/usr/local/lib
go run app/profil/main.go &
go run app/ctl/main.go
