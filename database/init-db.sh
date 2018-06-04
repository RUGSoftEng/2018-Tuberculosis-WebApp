#!/bin/bash

DB_NAME="TestDB"
name="root"
pass="root"

mysql -u ${name} --password=${pass} -Nse 'show tables' ${DB_NAME} | while read table; do mysql -u ${name} --password=${pass} ${DB_NAME} -e "SET foreign_key_checks = 0;drop table $table; SET foreign_key_checks=1;"; done

# Inserting General Statements
for f in DDL_statements.sql test_insert_statements.sql; do
    mysql -u ${name} --password=${pass} ${DB_NAME} < ${f}
done

# Inserting videos and quizzes
for f in vq-en.sql vq-nl.sql; do
    mysql -u ${name} --password=${pass} ${DB_NAME} < video-quiz/${f}
done

#Inserting FAQs
for f in faq-en.sql faq-nl.sql; do
    mysql -u ${name} --password=${pass} ${DB_NAME} < faq/${f}
done

echo "Done"
