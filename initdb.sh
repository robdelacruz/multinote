#!/bin/sh
if [ "$1" = "" ]; then exit 1; fi
cat create_tables.sql | sqlite3 $1
cat add_testdata.sql | sqlite3 $1

