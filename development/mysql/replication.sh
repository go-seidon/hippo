#!/bin/bash

# variable definition
h1_ct_name='hippo_mysql-db-1_1'
h1_db_root_username='root'
h1_db_root_password='toor'

h2_ct_name='hippo_mysql-db-2_1'
h2_db_host='localhost'
h2_db_port='3312'
h2_db_root_username='root'
h2_db_root_password='toor'
h2_db_username='hippo-h2'
h2_db_password='123456'

h3_ct_name='hippo_mysql-db-3_1'
h3_db_host='localhost'
h3_db_port='3313'
h3_db_root_username='root'
h3_db_root_password='toor'
h3_db_username='hippo-h3'
h3_db_password='123456'

h4_ct_name='hippo_mysql-db-4_1'
h4_db_host='localhost'
h4_db_port='3314'
h4_db_root_username='root'
h4_db_root_password='toor'
h4_db_username='hippo-h4'
h4_db_password='123456'

# function definition
docker-ip() {
  docker inspect --format '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' "$@"
}

# 1. turn off docker-compose (if any)
echo "[1] turning off docker compose...  "
docker-compose stop mysql-db-1 mysql-db-2 mysql-db-3 mysql-db-4
docker-compose rm -v -f mysql-db-1 mysql-db-2 mysql-db-3 mysql-db-4
printf "[DONE]\n\n"

# 2. remove mysql-db data volume (if any)
echo "[2] removing data volume...  "
docker volume rm hippo_mysql-db-1-data
docker volume rm hippo_mysql-db-2-data
docker volume rm hippo_mysql-db-3-data
docker volume rm hippo_mysql-db-4-data
printf "[DONE]\n\n"

# 3. rebuild and run docker-compose
echo "[3] rebuilding docker compose...  "
docker-compose build
docker-compose up -d
printf "[DONE]\n\n"

sleep 4

# 4. try to connect to m1 db
echo "[4] connecting to $h1_ct_name...  "
until docker exec "$h1_ct_name" sh -c "mysql -u$h1_db_root_username -p$h1_db_root_password -e ';'"
do
  echo "Waiting for '$h1_ct_name' database connection..."
  sleep 4
done
printf "[DONE]\n\n"

# 5. grant user for replication in m1 db
echo "[5] granting user in $h1_ct_name for replication...  "
grant_r1_stmt='CREATE USER "'$h2_db_username'"@"'%'" IDENTIFIED WITH mysql_native_password BY "'$h2_db_password'"; GRANT REPLICATION SLAVE ON *.* TO "'$h2_db_username'"@"%"; FLUSH PRIVILEGES;'
docker exec $h1_ct_name sh -c "mysql -u$h1_db_root_username -p$h1_db_root_password -e '$grant_r1_stmt'"

grant_r1_res_stmt='SHOW GRANTS FOR "'$h2_db_username'"@"%"'
grant_r1_res=`docker exec $h1_ct_name sh -c "mysql -u$h1_db_root_username -p$h1_db_root_password -e '$grant_r1_res_stmt'"`
echo "$grant_r1_res"
printf "[DONE]\n\n"

# 6. check master status in m1 db
echo "[6] checking master status in $h1_ct_name...  "
h1_status_res=`docker exec $h1_ct_name sh -c "mysql -u$h1_db_root_username -p$h1_db_root_password -e 'SHOW MASTER STATUS'"`
h1_log_file=`echo $h1_status_res | awk '{print $6}'`
h1_log_position=`echo $h1_status_res | awk '{print $7}'`
echo "log file            : $h1_log_file"
echo "log position        : $h1_log_position"
printf "[DONE]\n\n"

sleep 3

# 7. try to connect to r1 db
echo "[7] connecting to $h2_ct_name...  "
until docker exec "$h2_ct_name" sh -c "mysql -u$h2_db_root_username -p$h2_db_root_password -e ';'"
do
  echo "Waiting for '$h2_ct_name' database connection..."
  sleep 4
done
printf "[DONE]\n\n"

# 8. adding replica in h2 db
echo "[8] adding replica in $h2_ct_name...  "
echo "master host         : $(docker-ip $h1_ct_name)"
echo "username            : $h2_db_username"
echo "password            : $h2_db_password"
echo "log file            : $h1_log_file"
echo "log position        : $h1_log_position"
add_h2_stmt='STOP SLAVE; CHANGE MASTER TO MASTER_HOST="'$(docker-ip $h1_ct_name)'", MASTER_USER="'$h2_db_username'", MASTER_PASSWORD="'$h2_db_password'", MASTER_LOG_FILE="'$h1_log_file'", MASTER_LOG_POS='$h1_log_position'; START SLAVE;'
add_h2_res=`docker exec $h2_ct_name sh -c "mysql -u$h2_db_root_username -p$h2_db_root_password -e '$add_h2_stmt'"`
echo "$add_h2_res"
printf "[DONE]\n\n"

# 9. showing replica status
echo "[9] showing $h2_ct_name status...  "
h2_status_res=`docker exec $h2_ct_name sh -c "mysql -u$h2_db_root_username -p$h2_db_root_password -e 'SHOW SLAVE STATUS\G'"`
echo "$h2_status_res"
printf "[DONE]\n\n"

# 10. grant user for replication in h2 db
echo "[10] granting user in $h2_ct_name for replication...  "

# granting h3 in h2
grant_h3_stmt='CREATE USER "'$h3_db_username'"@"'%'" IDENTIFIED WITH mysql_native_password BY "'$h3_db_password'"; GRANT REPLICATION SLAVE ON *.* TO "'$h3_db_username'"@"%"; FLUSH PRIVILEGES;'
docker exec $h2_ct_name sh -c "mysql -u$h2_db_root_username -p$h2_db_root_password -e '$grant_h3_stmt'"

grant_h3_res_stmt='SHOW GRANTS FOR "'$h3_db_username'"@"%"'
grant_h3_res=`docker exec $h2_ct_name sh -c "mysql -u$h2_db_root_username -p$h2_db_root_password -e '$grant_h3_res_stmt'"`
echo "$grant_h3_res"
printf "[DONE]\n\n"

# granting h4 in h2
grant_h4_stmt='CREATE USER "'$h4_db_username'"@"'%'" IDENTIFIED WITH mysql_native_password BY "'$h4_db_password'"; GRANT REPLICATION SLAVE ON *.* TO "'$h4_db_username'"@"%"; FLUSH PRIVILEGES;'
docker exec $h2_ct_name sh -c "mysql -u$h2_db_root_username -p$h2_db_root_password -e '$grant_h4_stmt'"

grant_h4_res_stmt='SHOW GRANTS FOR "'$h4_db_username'"@"%"'
grant_h4_res=`docker exec $h2_ct_name sh -c "mysql -u$h2_db_root_username -p$h2_db_root_password -e '$grant_h4_res_stmt'"`
echo "$grant_h4_res"
printf "[DONE]\n\n"

# 11. check master status in h2 db
echo "[11] checking master status in $h2_ct_name...  "
h2_status_res=`docker exec $h2_ct_name sh -c "mysql -u$h2_db_root_username -p$h2_db_root_password -e 'SHOW MASTER STATUS'"`
h2_log_file=`echo $h2_status_res | awk '{print $6}'`
h2_log_position=`echo $h2_status_res | awk '{print $7}'`
echo "log file            : $h2_log_file"
echo "log position        : $h2_log_position"
printf "[DONE]\n\n"

sleep 3

# 12. try to connect to h3 db
echo "[12] connecting to $h3_ct_name...  "
until docker exec "$h3_ct_name" sh -c "mysql -u$h3_db_root_username -p$h3_db_root_password -e ';'"
do
  echo "Waiting for '$h3_ct_name' database connection..."
  sleep 4
done
printf "[DONE]\n\n"

# 13. adding replica in h3 db
echo "[13] adding replica in $h3_ct_name...  "
echo "master host         : $(docker-ip $h2_ct_name)"
echo "username            : $h3_db_username"
echo "password            : $h3_db_password"
echo "log file            : $h2_log_file"
echo "log position        : $h2_log_position"
add_h3_stmt='STOP SLAVE; CHANGE MASTER TO MASTER_HOST="'$(docker-ip $h2_ct_name)'", MASTER_USER="'$h3_db_username'", MASTER_PASSWORD="'$h3_db_password'", MASTER_LOG_FILE="'$h2_log_file'", MASTER_LOG_POS='$h2_log_position'; START SLAVE;'
add_h3_res=`docker exec $h3_ct_name sh -c "mysql -u$h3_db_root_username -p$h3_db_root_password -e '$add_h3_stmt'"`
echo "$add_h3_res"
printf "[DONE]\n\n"

# 14. showing replica status
echo "[14] showing $h3_ct_name status...  "
h3_status_res=`docker exec $h3_ct_name sh -c "mysql -u$h3_db_root_username -p$h3_db_root_password -e 'SHOW SLAVE STATUS\G'"`
echo "$h3_status_res"
printf "[DONE]\n\n"

# 15. try to connect to h4 db
echo "[15] connecting to $h4_ct_name...  "
until docker exec "$h4_ct_name" sh -c "mysql -u$h4_db_root_username -p$h4_db_root_password -e ';'"
do
  echo "Waiting for '$h4_ct_name' database connection..."
  sleep 4
done
printf "[DONE]\n\n"

# 16. adding replica in h4 db
echo "[16] adding replica in $h4_ct_name...  "
echo "master host         : $(docker-ip $h2_ct_name)"
echo "username            : $h4_db_username"
echo "password            : $h4_db_password"
echo "log file            : $h2_log_file"
echo "log position        : $h2_log_position"
add_h4_stmt='STOP SLAVE; CHANGE MASTER TO MASTER_HOST="'$(docker-ip $h2_ct_name)'", MASTER_USER="'$h4_db_username'", MASTER_PASSWORD="'$h4_db_password'", MASTER_LOG_FILE="'$h2_log_file'", MASTER_LOG_POS='$h2_log_position'; START SLAVE;'
add_h4_res=`docker exec $h4_ct_name sh -c "mysql -u$h4_db_root_username -p$h4_db_root_password -e '$add_h4_stmt'"`
echo "$add_h4_res"
printf "[DONE]\n\n"

# 17. showing replica status
echo "[17] showing $h4_ct_name status...  "
h4_status_res=`docker exec $h4_ct_name sh -c "mysql -u$h4_db_root_username -p$h4_db_root_password -e 'SHOW SLAVE STATUS\G'"`
echo "$h4_status_res"
printf "[DONE]\n\n"

# 18. done
echo "exiting...  "
