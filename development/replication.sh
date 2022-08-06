#!/bin/bash

# variable definition
m1_ct_name='mysql-db'
m1_db_root_username='root'
m1_db_root_password='toor'
r1_ct_name='mysql-db-r1'
r1_db_root_username='root'
r1_db_root_password='toor'
r1_db_username='goseidon-r1'
r1_db_password='123456'

# function definition
docker-ip() {
  docker inspect --format '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' "$@"
}

# 1. turn off docker-compose (if any)
echo "turning off docker compose...  "
docker-compose down
printf "[DONE]\n\n"

# 2. remove mysql-db replica data volume (if any)
echo "removing $m1_ct_name data volume...  "
docker volume rm goseidon-local_mysql-db-r1-data
docker volume rm goseidon-local_mysql-db-r2-data
printf "[DONE]\n\n"

# 3. rebuild and run docker-compose
echo "rebuilding docker compose...  "
docker-compose build
docker-compose up -d
printf "[DONE]\n\n"

# 4. try to connect to mysql-db using root user
echo "connecting to $m1_ct_name...  "
until docker exec "$m1_ct_name" sh -c "export MYSQL_PWD=$m1_db_root_password; mysql -u $m1_db_root_username -e ';'"
do
  echo "Waiting for '$m1_ct_name' database connection..."
  sleep 4
done
printf "[DONE]\n\n"

# 5. create user for replication in mysql-db using root user
echo "creating user for replication...  "
create_rpl_r1_stmt='CREATE USER "'$r1_db_username'"@"%" IDENTIFIED WITH mysql_native_password BY "'$r1_db_password'"'
docker exec $m1_ct_name sh -c "export MYSQL_PWD=$m1_db_root_password; mysql -u $m1_db_root_username -e '$grant_rpl_r1_stmt'"
printf "[DONE]\n\n"

# 6. grant user for replication in mysql-db using root user
echo "granting user for replication...  "
grant_rpl_r1_stmt='GRANT REPLICATION SLAVE ON *.* TO "'$r1_db_username'"@"%" IDENTIFIED BY "'$r1_db_password'"; FLUSH PRIVILEGES;'
docker exec $m1_ct_name sh -c "export MYSQL_PWD=$m1_db_root_password; mysql -u $m1_db_root_username -e '$grant_rpl_r1_stmt'"
printf "[DONE]\n\n"

# 7. show grant result
echo "showing grant result...  "
grant_r1_res_stmt='SHOW GRANTS FOR "'$r1_db_username'"@"%"'
grant_r1_res=`docker exec $m1_ct_name sh -c "export MYSQL_PWD=$m1_db_root_password; mysql -u $m1_db_root_username -e '$grant_r1_res_stmt'"`
echo "$grant_r1_res"
printf "[DONE]\n\n"

# 8. check master status
echo "checking master status...  "
m1_status_res=`docker exec $m1_ct_name sh -c "export MYSQL_PWD=$m1_db_root_password; mysql -u $m1_db_root_username -e 'SHOW MASTER STATUS'"`
m1_log_file=`echo $m1_status_res | awk '{print $5}'`
m1_log_position=`echo $m1_status_res | awk '{print $6}'`
echo "log file            : $m1_log_file"
echo "log position        : $m1_log_position"
printf "[DONE]\n\n"

# 9. try to connect to mysql-db-r1 using root user
echo "connecting to $r1_ct_name...  "
until docker exec "$r1_ct_name" sh -c "export MYSQL_PWD=$r1_db_root_password; mysql -u $r1_db_root_username -e ';'"
do
  echo "Waiting for '$r1_ct_name' database connection..."
  sleep 4
done
printf "[DONE]\n\n"

# 10. adding replica 1
echo "adding replica 1...  "
add_r1_stmt='CHANGE MASTER TO MASTER_HOST="'$(docker-ip $m1_ct_name)'", MASTER_USER="'$r1_db_username'", MASTER_PASSWORD="'$r1_db_password'", MASTER_LOG_FILE="'$m1_log_file'", MASTER_LOG_POS='$m1_log_position'; STOP SLAVE; START SLAVE;'
add_r1_res=`docker exec $r1_ct_name sh -c "export MYSQL_PWD=$r1_db_root_password; mysql -u $r1_db_root_username -e '$add_r1_stmt'"`
echo "$add_r1_res"
printf "[DONE]\n\n"

# 11. showing replica 1 info
echo "showing replica 1 info...  "
echo "master host         : $(docker-ip $m1_ct_name)"
echo "replica 1 username  : $r1_db_username"
echo "replica 1 password  : $r1_db_password"
echo "log file            : $m1_log_file"
echo "log position        : $m1_log_position"
printf "[DONE]\n\n"

# 12. showing replica 1 status
echo "showing replica 1 status...  "
r1_status_res=`docker exec $r1_ct_name sh -c "export MYSQL_PWD=$r1_db_root_password; mysql -u $r1_db_root_username -e 'SHOW SLAVE STATUS\G'"`
echo "$r1_status_res"
printf "[DONE]\n\n"

# 13. running db migration on replica 1
echo "running db migration on replica 1...  "
r1_db_migration='mysql://admin:123456@tcp(localhost:3311)/goseidon_local?x-tls-insecure-skip-verify=true'
migrate -database "$r1_db_migration" -path ./migration/mysql up

# 14. done
echo "exiting...  "
