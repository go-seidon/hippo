#!/bin/bash

# variable definition
m1_ct_name='hippo_mysql-db_1'
m1_db_root_username='root'
m1_db_root_password='toor'

r1_ct_name='hippo_mysql-db-r1_1'
r1_db_host='localhost'
r1_db_port='3311'
r1_db_root_username='root'
r1_db_root_password='toor'
r1_db_username='hippo-r1'
r1_db_password='123456'

r2_ct_name='hippo_mysql-db-r2_1'
r2_db_host='localhost'
r2_db_port='3312'
r2_db_root_username='root'
r2_db_root_password='toor'
r2_db_username='hippo-r2'
r2_db_password='123456'

r3_ct_name='hippo_mysql-db-r3_1'
r3_db_host='localhost'
r3_db_port='3313'
r3_db_root_username='root'
r3_db_root_password='toor'
r3_db_username='hippo-r3'
r3_db_password='123456'

# function definition
docker-ip() {
  docker inspect --format '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' "$@"
}

# 1. turn off docker-compose (if any)
echo "[1] turning off docker compose...  "
docker-compose stop mysql-db mysql-db-r1 mysql-db-r2 mysql-db-r3
docker-compose rm -v -f mysql-db-r1 mysql-db-r2 mysql-db-r3
printf "[DONE]\n\n"

# 2. remove mysql-db replica data volume (if any)
echo "[2] removing $m1_ct_name data volume...  "
docker volume rm hippo_mysql-db-r1-data
docker volume rm hippo_mysql-db-r2-data
docker volume rm hippo_mysql-db-r3-data
printf "[DONE]\n\n"

# 3. rebuild and run docker-compose
echo "[3] rebuilding docker compose...  "
docker-compose build
docker-compose up -d
printf "[DONE]\n\n"

sleep 3

# 4. try to connect to m1 db
echo "[4] connecting to $m1_ct_name...  "
until docker exec "$m1_ct_name" sh -c "export MYSQL_PWD=$m1_db_root_password; mysql -u $m1_db_root_username -e ';'"
do
  echo "Waiting for '$m1_ct_name' database connection..."
  sleep 4
done
printf "[DONE]\n\n"

# 5. grant user for replication in m1 db
echo "[5] granting user in $m1_ct_name for replication...  "
grant_r1_stmt='GRANT REPLICATION SLAVE ON *.* TO "'$r1_db_username'"@"%" IDENTIFIED BY "'$r1_db_password'"; FLUSH PRIVILEGES;'
docker exec $m1_ct_name sh -c "export MYSQL_PWD=$m1_db_root_password; mysql -u $m1_db_root_username -e '$grant_r1_stmt'"

grant_r1_res_stmt='SHOW GRANTS FOR "'$r1_db_username'"@"%"'
grant_r1_res=`docker exec $m1_ct_name sh -c "export MYSQL_PWD=$m1_db_root_password; mysql -u $m1_db_root_username -e '$grant_r1_res_stmt'"`
echo "$grant_r1_res"
printf "[DONE]\n\n"

# 6. check master status in m1 db
echo "[6] checking master status in $m1_ct_name...  "
m1_status_res=`docker exec $m1_ct_name sh -c "export MYSQL_PWD=$m1_db_root_password; mysql -u $m1_db_root_username -e 'SHOW MASTER STATUS'"`
m1_log_file=`echo $m1_status_res | awk '{print $5}'`
m1_log_position=`echo $m1_status_res | awk '{print $6}'`
echo "log file            : $m1_log_file"
echo "log position        : $m1_log_position"
printf "[DONE]\n\n"

sleep 3

# 7. try to connect to r1 db
echo "[7] connecting to $r1_ct_name...  "
until docker exec "$r1_ct_name" sh -c "export MYSQL_PWD=$r1_db_root_password; mysql -u $r1_db_root_username -e ';'"
do
  echo "Waiting for '$r1_ct_name' database connection..."
  sleep 4
done
printf "[DONE]\n\n"

# 8. adding replica in r1 db
echo "[8] adding replica in $r1_ct_name...  "
echo "master host         : $(docker-ip $m1_ct_name)"
echo "username            : $r1_db_username"
echo "password            : $r1_db_password"
echo "log file            : $m1_log_file"
echo "log position        : $m1_log_position"
add_r1_stmt='STOP SLAVE; CHANGE MASTER TO MASTER_HOST="'$(docker-ip $m1_ct_name)'", MASTER_USER="'$r1_db_username'", MASTER_PASSWORD="'$r1_db_password'", MASTER_LOG_FILE="'$m1_log_file'", MASTER_LOG_POS='$m1_log_position'; START SLAVE;'
add_r1_res=`docker exec $r1_ct_name sh -c "export MYSQL_PWD=$r1_db_root_password; mysql -u $r1_db_root_username -e '$add_r1_stmt'"`
echo "$add_r1_res"
printf "[DONE]\n\n"

# 9. showing replica status
echo "[9] showing $r1_ct_name status...  "
r1_status_res=`docker exec $r1_ct_name sh -c "export MYSQL_PWD=$r1_db_root_password; mysql -u $r1_db_root_username -e 'SHOW SLAVE STATUS\G'"`
echo "$r1_status_res"
printf "[DONE]\n\n"

# 10. grant user for replication in r1 db
echo "[10] granting user in $r1_ct_name for replication...  "

# granting r2 in r1
grant_r2_stmt='GRANT REPLICATION SLAVE ON *.* TO "'$r2_db_username'"@"%" IDENTIFIED BY "'$r2_db_password'"; FLUSH PRIVILEGES;'
docker exec $r1_ct_name sh -c "export MYSQL_PWD=$r1_db_root_password; mysql -u $r1_db_root_username -e '$grant_r2_stmt'"

grant_r2_res_stmt='SHOW GRANTS FOR "'$r2_db_username'"@"%"'
grant_r2_res=`docker exec $r1_ct_name sh -c "export MYSQL_PWD=$r1_db_root_password; mysql -u $r1_db_root_username -e '$grant_r2_res_stmt'"`
echo "$grant_r2_res"
printf "[DONE]\n\n"

# granting r3 in r1
grant_r3_stmt='GRANT REPLICATION SLAVE ON *.* TO "'$r3_db_username'"@"%" IDENTIFIED BY "'$r3_db_password'"; FLUSH PRIVILEGES;'
docker exec $r1_ct_name sh -c "export MYSQL_PWD=$r1_db_root_password; mysql -u $r1_db_root_username -e '$grant_r3_stmt'"

grant_r3_res_stmt='SHOW GRANTS FOR "'$r3_db_username'"@"%"'
grant_r3_res=`docker exec $r1_ct_name sh -c "export MYSQL_PWD=$r1_db_root_password; mysql -u $r1_db_root_username -e '$grant_r3_res_stmt'"`
echo "$grant_r3_res"
printf "[DONE]\n\n"

# 11. check master status in r1 db
echo "[11] checking master status in $r1_ct_name...  "
r1_status_res=`docker exec $r1_ct_name sh -c "export MYSQL_PWD=$r1_db_root_password; mysql -u $r1_db_root_username -e 'SHOW MASTER STATUS'"`
r1_log_file=`echo $r1_status_res | awk '{print $5}'`
r1_log_position=`echo $r1_status_res | awk '{print $6}'`
echo "log file            : $r1_log_file"
echo "log position        : $r1_log_position"
printf "[DONE]\n\n"

sleep 3

# 12. try to connect to r2 db
echo "[12] connecting to $r2_ct_name...  "
until docker exec "$r2_ct_name" sh -c "export MYSQL_PWD=$r2_db_root_password; mysql -u $r2_db_root_username -e ';'"
do
  echo "Waiting for '$r2_ct_name' database connection..."
  sleep 4
done
printf "[DONE]\n\n"

# 13. adding replica in r2 db
echo "[13] adding replica in $r2_ct_name...  "
echo "master host         : $(docker-ip $r1_ct_name)"
echo "username            : $r2_db_username"
echo "password            : $r2_db_password"
echo "log file            : $r1_log_file"
echo "log position        : $r1_log_position"
add_r2_stmt='STOP SLAVE; CHANGE MASTER TO MASTER_HOST="'$(docker-ip $r1_ct_name)'", MASTER_USER="'$r2_db_username'", MASTER_PASSWORD="'$r2_db_password'", MASTER_LOG_FILE="'$r1_log_file'", MASTER_LOG_POS='$r1_log_position'; START SLAVE;'
add_r2_res=`docker exec $r2_ct_name sh -c "export MYSQL_PWD=$r2_db_root_password; mysql -u $r2_db_root_username -e '$add_r2_stmt'"`
echo "$add_r2_res"
printf "[DONE]\n\n"

# 14. showing replica status
echo "[14] showing $r2_ct_name status...  "
r2_status_res=`docker exec $r2_ct_name sh -c "export MYSQL_PWD=$r2_db_root_password; mysql -u $r2_db_root_username -e 'SHOW SLAVE STATUS\G'"`
echo "$r2_status_res"
printf "[DONE]\n\n"

# 15. try to connect to r3 db
echo "[15] connecting to $r3_ct_name...  "
until docker exec "$r3_ct_name" sh -c "export MYSQL_PWD=$r3_db_root_password; mysql -u $r3_db_root_username -e ';'"
do
  echo "Waiting for '$r3_ct_name' database connection..."
  sleep 4
done
printf "[DONE]\n\n"

# 16. adding replica in r3 db
echo "[16] adding replica in $r3_ct_name...  "
echo "master host         : $(docker-ip $r1_ct_name)"
echo "username            : $r3_db_username"
echo "password            : $r3_db_password"
echo "log file            : $r1_log_file"
echo "log position        : $r1_log_position"
add_r3_stmt='STOP SLAVE; CHANGE MASTER TO MASTER_HOST="'$(docker-ip $r1_ct_name)'", MASTER_USER="'$r3_db_username'", MASTER_PASSWORD="'$r3_db_password'", MASTER_LOG_FILE="'$r1_log_file'", MASTER_LOG_POS='$r1_log_position'; START SLAVE;'
add_r3_res=`docker exec $r3_ct_name sh -c "export MYSQL_PWD=$r3_db_root_password; mysql -u $r3_db_root_username -e '$add_r3_stmt'"`
echo "$add_r3_res"
printf "[DONE]\n\n"

# 17. showing replica status
echo "[17] showing $r3_ct_name status...  "
r3_status_res=`docker exec $r3_ct_name sh -c "export MYSQL_PWD=$r3_db_root_password; mysql -u $r3_db_root_username -e 'SHOW SLAVE STATUS\G'"`
echo "$r3_status_res"
printf "[DONE]\n\n"

# 18. done
echo "exiting...  "
