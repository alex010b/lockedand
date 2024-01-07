#!/bin/bash

# MySQL connection parameters
MYSQL_USER="root"
MYSQL_PASSWORD="ImABall!@#@@(@11212"
MYSQL_DATABASE="lockedand"
TABLE_TO_DROP="people"
rm -rf database/*
# Check if the MySQL command-line tool is installed
if ! command -v mysql &> /dev/null; then
	    echo "MySQL command-line tool is not installed. Please install it."
	        exit 1
fi

# Log in to MySQL and drop the table
mysql -u "$MYSQL_USER" -p"$MYSQL_PASSWORD" "$MYSQL_DATABASE" <<EOF
DROP TABLE IF EXISTS $TABLE_TO_DROP;
EOF

echo "Table '$TABLE_TO_DROP' dropped successfully."
