#!/bin/bash

cat > $CONFFILE <<- EOF
{
   "version": "$VERSION",
   "dev_addr": "$DEV_ADDR",
   "dev_port": $DEV_PORT,
   "db_config": {
      "conn_string": "mongodb://$DEV_ADDR:$MONGODB_PORT",
      "database": "spotdb",
      "collection": "garages"
   }
}
EOF
