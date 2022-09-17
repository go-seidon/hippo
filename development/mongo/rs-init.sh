#!/bin/bash

DELAY=25
h1_db_host='localhost'
h1_db_port='27031'
h1_db_root_username='root'
h1_db_root_password='toor'

mongo --host ${h1_db_host} --port ${h1_db_port} -u ${h1_db_root_username} -p ${h1_db_root_password} <<EOF
var config = {
  "_id": "rs-goseidon",
  "version": 1,
  "members": [
    {
      "_id": 1,
      "host": "mongo-db-1:27031"
    },
    {
      "_id": 2,
      "host": "mongo-db-2:27032"
    },
    {
      "_id": 3,
      "host": "mongo-db-3:27033"
    }
  ]
};
rs.initiate(config, { force: true });
EOF

echo "****** Waiting for ${DELAY} seconds for replica set configuration to be applied ******"

sleep $DELAY
