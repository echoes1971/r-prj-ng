#!/bin/sh
# Replace environment variables in env-config.js at runtime

ENV_CONFIG_FILE=/usr/share/nginx/html/env-config.js

# Replace placeholders with actual environment variables
envsubst '${REACT_APP_SITE_TITLE} ${REACT_APP_ENDPOINT} ${REACT_APP_HOME_OBJECT_ID} ${REACT_APP_WEBMASTER_GROUP_ID} ${REACT_APP_APP_NAME} ${REACT_APP_APP_VERSION} ${REACT_APP_SITE_COPYRIGHT}' < $ENV_CONFIG_FILE > $ENV_CONFIG_FILE.tmp
mv $ENV_CONFIG_FILE.tmp $ENV_CONFIG_FILE

echo "Environment variables injected:"
cat $ENV_CONFIG_FILE

# Start nginx
exec nginx -g 'daemon off;'
