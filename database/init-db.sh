#!/bin/bash

DB_NAME="TestDB"
name="root"
pass="root"

mysql -u ${name} --password=${pass} <<< "DROP DATABASE IF EXISTS ${DB_NAME}"
mysql -u ${name} --password=${pass} <<< "CREATE DATABASE ${DB_NAME}"

# Inserting General Statements
for f in DDL_statements.sql test_insert_statements.sql; do
	 mysql -u ${name} --password=${pass} ${DB_NAME} < ${f}
done

# Inserting videos and quizzes
for f in vq-en.sql vq-nl.sql; do
    	 mysql -u ${name} --password=${pass} ${DB_NAME} < video-quiz/${f}
done

echo "Done"
