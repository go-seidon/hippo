#!/bin/bash

# variable definition
replica_set_name='rs-hippo'

h1_ct_name='hippo_mongo-db-1_1'
h1_db_host='localhost'
h1_db_port='27031'
h1_db_root_username='root'
h1_db_root_password='toor'

h2_ct_name='hippo_mongo-db-2_1'
h2_db_host='localhost'
h2_db_port='27032'
h2_db_root_username='root'
h2_db_root_password='toor'

h3_ct_name='hippo_mongo-db-3_1'
h3_db_host='localhost'
h3_db_port='27033'
h3_db_root_username='root'
h3_db_root_password='toor'

# function definition
docker-ip() {
  docker inspect --format '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' "$@"
}

connect-mongo() {
  docker exec $1 sh -c 'mongo --host '$2' --port '$3' -u '$4' -p '$5''
}

execute-query() {
  docker exec $1 sh -c 'mongo --host '$2' --port '$3' -u '$4' -p '$5' --eval "'$6'"'
}

# 1. turn off docker-compose (if any)
echo "[1] turning off docker compose...  "
docker-compose stop mongo-db-1 mongo-db-2 mongo-db-3
docker-compose rm -v -f mongo-db-1 mongo-db-2 mongo-db-3
printf "[DONE]\n\n"

# 2. remove mongo-db data volume (if any)
echo "[2] removing data volume...  "
docker volume rm hippo_mongo-db-1-data
docker volume rm hippo_mongo-db-2-data
docker volume rm hippo_mongo-db-3-data
printf "[DONE]\n\n"

# 3. rebuild and run docker-compose
echo "[3] rebuilding docker compose...  "
docker-compose build
docker-compose up -d
printf "[DONE]\n\n"

sleep 3

# 4. try to connect to m1 db
echo "[4] connecting to $h1_ct_name...  "
until connect-mongo $h1_ct_name $h1_db_host $h1_db_port $h1_db_root_username $h1_db_root_password
do
  echo "Waiting for '$h1_ct_name' database connection..."
  sleep 4
done
printf "[DONE]\n\n"

# 5. check master status in m1 db
echo "[5] checking master status in $h1_ct_name...  "
h1_status_res=`execute-query $h1_ct_name $h1_db_host $h1_db_port $h1_db_root_username $h1_db_root_password 'rs.status();'`
echo "rs status            : $h1_status_res"
printf "[DONE]\n\n"

sleep 3

# 6. initiate replica set with default primary member
echo "[6] initiate replica set: $replica_set_name...  "
h1_init_res=`docker exec $h1_ct_name sh -c "chmod +x ./data/script/rs-init.sh && ./data/script/rs-init.sh"`
echo "init status            : $h1_init_res"
printf "[DONE]\n\n"

# 7. check replication status in m1 db
echo "[7] checking replication status in $h1_ct_name...  "
h1_status_res=`execute-query $h1_ct_name $h1_db_host $h1_db_port $h1_db_root_username $h1_db_root_password 'rs.status();'`
echo "rs status            : $h1_status_res"
printf "[DONE]\n\n"

# 8. enable read in 2 db
echo "[8] enable read in $h2_ct_name...  "
h2_enable_res=`execute-query $h2_ct_name $h2_db_host $h2_db_port $h2_db_root_username $h2_db_root_password 'rs.secondaryOk();'`
echo "enable read 2            : $h2_enable_res"
printf "[DONE]\n\n"

# 9. enable read in 3 db
echo "[9] enable read in $h3_ct_name...  "
h3_enable_res=`execute-query $h3_ct_name $h3_db_host $h3_db_port $h3_db_root_username $h3_db_root_password 'rs.secondaryOk();'`
echo "enable read 3            : $h3_enable_res"
printf "[DONE]\n\n"

# 10. done
echo "exiting...  "
