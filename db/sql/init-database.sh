#!/usr/bin/env bash
#wait for the MySQL Server to come up

#run the setup script to create the DB and the schema in the DB
mysql -u daily -pdaily daily < "/docker-entrypoint-initdb.d/delete-tables.sql"
mysql -u daily -pdaily daily < "/docker-entrypoint-initdb.d/create_users_table.sql"
mysql -u daily -pdaily daily < "/docker-entrypoint-initdb.d/create_articles_table.sql"
mysql -u daily -pdaily daily < "/docker-entrypoint-initdb.d/create_images_table.sql"

