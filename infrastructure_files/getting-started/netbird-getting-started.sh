#!/bin/bash

set -e

handle_request_command_status() {
  PARSED_RESPONSE=$1
  FUNCTION_NAME=$2
  RESPONSE=$3
  if [[ $PARSED_RESPONSE -ne 0 ]]; then
    echo "ERROR calling $FUNCTION_NAME: $(echo $RESPONSE | jq -r '.message')" > /dev/stderr
    exit 1
  fi
}

handle_zitadel_request_response() {
  PARSED_RESPONSE=$1
  FUNCTION_NAME=$2
  RESPONSE=$3
  if [[ $PARSED_RESPONSE == "null" ]]; then
    echo "ERROR calling $FUNCTION_NAME: $(echo $RESPONSE | jq -r '.message')" > /dev/stderr
    exit 1
  fi
}

check_docker_compose() {
  if command -v docker-compose &> /dev/null
  then
      echo "docker-compose"
      return
  fi
  if docker compose --help &> /dev/null
  then
      echo "docker compose"
      return
  fi

  echo "docker-compose is not installed or not in PATH. Please follow the steps from the official guide: https://docs.docker.com/engine/install/"
  exit 1
}

check_jq() {
  if ! command -v jq &> /dev/null
  then
    echo "jq is not installed or not in PATH, please install with your package manager. e.g. sudo apt install jq"
    exit 1
  fi
}

init_crdb() {
  echo "Initializing crdb"
  $DOCKER_COMPOSE_COMMAND up -d crdb --wait --wait-timeout 90
  $DOCKER_COMPOSE_COMMAND exec -t crdb /bin/bash -c "cp /cockroach/certs/* /zitadel-certs/ && cockroach cert create-client --overwrite --certs-dir /zitadel-certs/ --ca-key /zitadel-certs/ca.key zitadel_user && chown -R 1000:1000 /zitadel-certs/"
  handle_request_command_status $? "init_crdb failed" ""
}

get_main_ip_address() {
  if [[ "$OSTYPE" == "darwin"* ]]; then
    interface=$(route -n get default | grep 'interface:' | awk '{print $2}')
    ip_address=$(ifconfig $interface | grep 'inet ' | awk '{print $2}')
  else
    interface=$(ip route | grep default | awk '{print $5}' | head -n 1)
    ip_address=$(ip addr show $interface | grep 'inet ' | awk '{print $2}' | cut -d'/' -f1)
  fi

  echo $ip_address
}

wait_pat() {
  PAT_PATH=$1
  set +e
  while true; do
    if [[ -f "$PAT_PATH" ]]; then
      break
    fi
    echo -n " ."
    sleep 1
  done
  echo " done"
  set -e
}

wait_api() {
    INSTANCE_URL=$1
    PAT=$2
    set +e
    while true; do
      curl -s --fail -o /dev/null "$INSTANCE_URL/auth/v1/users/me" -H "Authorization: Bearer $PAT"
      if [[ $? -eq 0 ]]; then
        break
      fi
      echo -n " ."
      sleep 1
    done
    echo " done"
    set -e
}

create_new_project() {
  INSTANCE_URL=$1
  PAT=$2
  PROJECT_NAME="NETBIRD"

  RESPONSE=$(
    curl -sS -X POST "$INSTANCE_URL/management/v1/projects" \
      -H "Authorization: Bearer $PAT" \
      -H "Content-Type: application/json" \
      -d '{"name": "'"$PROJECT_NAME"'"}'
  )
  PARSED_RESPONSE=$(echo "$RESPONSE" | jq -r '.id')
  handle_zitadel_request_response $PARSED_RESPONSE "create_new_project" "$RESPONSE"
  echo "$PARSED_RESPONSE"
}

create_new_application() {
  INSTANCE_URL=$1
  PAT=$2
  APPLICATION_NAME="netbird"

  RESPONSE=$(
    curl -sS -X POST "$INSTANCE_URL/management/v1/projects/$PROJECT_ID/apps/oidc" \
      -H "Authorization: Bearer $PAT" \
      -H "Content-Type: application/json" \
      -d '{
    "name": "'"$APPLICATION_NAME"'",
    "redirectUris": [
      "'"$BASE_REDIRECT_URL"'/nb-auth",
      "'"$BASE_REDIRECT_URL"'/nb-silent-auth",
      "http://localhost:53000/",
      "http://localhost:54000/"
    ],
    "postLogoutRedirectUris": [
       "'"$BASE_REDIRECT_URL"'/nb-auth"
    ],
    "RESPONSETypes": [
      "OIDC_RESPONSE_TYPE_CODE"
    ],
    "grantTypes": [
      "OIDC_GRANT_TYPE_AUTHORIZATION_CODE",
      "OIDC_GRANT_TYPE_REFRESH_TOKEN"
    ],
    "appType": "OIDC_APP_TYPE_USER_AGENT",
    "authMethodType": "OIDC_AUTH_METHOD_TYPE_NONE",
    "version": "OIDC_VERSION_1_0",
    "devMode": '"$ZITADEL_DEV_MODE"',
    "accessTokenType": "OIDC_TOKEN_TYPE_JWT",
    "accessTokenRoleAssertion": true,
    "skipNativeAppSuccessPage": true
  }'
  )

  PARSED_RESPONSE=$(echo "$RESPONSE" | jq -r '.clientId')
  handle_zitadel_request_response $PARSED_RESPONSE "create_new_application" "$RESPONSE"
  echo "$PARSED_RESPONSE"
}

create_service_user() {
  INSTANCE_URL=$1
  PAT=$2

  RESPONSE=$(
    curl -sS -X POST "$INSTANCE_URL/management/v1/users/machine" \
      -H "Authorization: Bearer $PAT" \
      -H "Content-Type: application/json" \
      -d '{
            "userName": "netbird-service-account",
            "name": "Netbird Service Account",
            "description": "Netbird Service Account for IDP management",
            "accessTokenType": "ACCESS_TOKEN_TYPE_JWT"
      }'
  )
  PARSED_RESPONSE=$(echo "$RESPONSE" | jq -r '.userId')
  handle_zitadel_request_response $PARSED_RESPONSE "create_service_user" "$RESPONSE"
  echo "$PARSED_RESPONSE"
}

create_service_user_secret() {
  INSTANCE_URL=$1
  PAT=$2
  USER_ID=$3

  RESPONSE=$(
    curl -sS -X PUT "$INSTANCE_URL/management/v1/users/$USER_ID/secret" \
      -H "Authorization: Bearer $PAT" \
      -H "Content-Type: application/json" \
      -d '{}'
  )
  SERVICE_USER_CLIENT_ID=$(echo "$RESPONSE" | jq -r '.clientId')
  handle_zitadel_request_response $SERVICE_USER_CLIENT_ID "create_service_user_secret" "$RESPONSE"
  SERVICE_USER_CLIENT_SECRET=$(echo "$RESPONSE" | jq -r '.clientSecret')
  handle_zitadel_request_response $SERVICE_USER_CLIENT_SECRET "create_service_user_secret" "$RESPONSE"
}

add_organization_user_manager() {
  INSTANCE_URL=$1
  PAT=$2
  USER_ID=$3

  RESPONSE=$(
    curl -sS -X POST "$INSTANCE_URL/management/v1/orgs/me/members" \
      -H "Authorization: Bearer $PAT" \
      -H "Content-Type: application/json" \
      -d '{
            "userId": "'$USER_ID'",
            "roles": [
              "ORG_USER_MANAGER"
            ]
      }'
  )
  PARSED_RESPONSE=$(echo "$RESPONSE" | jq -r '.details.creationDate')
  handle_zitadel_request_response $PARSED_RESPONSE "add_organization_user_manager" "$RESPONSE"
  echo "$PARSED_RESPONSE"
}

create_admin_user() {
    INSTANCE_URL=$1
    PAT=$2
    USERNAME=$3
    PASSWORD=$4
    RESPONSE=$(
        curl -sS -X POST "$INSTANCE_URL/management/v1/users/human/_import" \
          -H "Authorization: Bearer $PAT" \
          -H "Content-Type: application/json" \
          -d '{
                "userName": "'$USERNAME'",
                "profile": {
                  "firstName": "Zitadel",
                  "lastName": "Admin"
                },
                "email": {
                  "email": "'$USERNAME'",
                  "isEmailVerified": true
                },
                "password": "'$PASSWORD'",
                "passwordChangeRequired": true
          }'
      )
      PARSED_RESPONSE=$(echo "$RESPONSE" | jq -r '.userId')
      handle_zitadel_request_response $PARSED_RESPONSE "create_admin_user" "$RESPONSE"
      echo "$PARSED_RESPONSE"
}

add_instance_admin() {
  INSTANCE_URL=$1
  PAT=$2
  USER_ID=$3

  RESPONSE=$(
    curl -sS -X POST "$INSTANCE_URL/admin/v1/members" \
      -H "Authorization: Bearer $PAT" \
      -H "Content-Type: application/json" \
      -d '{
            "userId": "'$USER_ID'",
            "roles": [
              "IAM_OWNER"
            ]
      }'
  )
  PARSED_RESPONSE=$(echo "$RESPONSE" | jq -r '.details.creationDate')
  handle_zitadel_request_response $PARSED_RESPONSE "add_instance_admin" "$RESPONSE"
  echo "$PARSED_RESPONSE"
}

delete_auto_service_user() {
  INSTANCE_URL=$1
  PAT=$2

  RESPONSE=$(
    curl -sS -X GET "$INSTANCE_URL/auth/v1/users/me" \
      -H "Authorization: Bearer $PAT" \
      -H "Content-Type: application/json" \
  )
  USER_ID=$(echo "$RESPONSE" | jq -r '.user.id')
  handle_zitadel_request_response $USER_ID "delete_auto_service_user_get_user" "$RESPONSE"

  RESPONSE=$(
      curl -sS -X DELETE "$INSTANCE_URL/admin/v1/members/$USER_ID" \
        -H "Authorization: Bearer $PAT" \
        -H "Content-Type: application/json" \
  )
  PARSED_RESPONSE=$(echo "$RESPONSE" | jq -r '.details.changeDate')
  handle_zitadel_request_response $PARSED_RESPONSE "delete_auto_service_user_remove_instance_permissions" "$RESPONSE"

  RESPONSE=$(
      curl -sS -X DELETE "$INSTANCE_URL/management/v1/orgs/me/members/$USER_ID" \
        -H "Authorization: Bearer $PAT" \
        -H "Content-Type: application/json" \
  )
  PARSED_RESPONSE=$(echo "$RESPONSE" | jq -r '.details.changeDate')
  handle_zitadel_request_response $PARSED_RESPONSE "delete_auto_service_user_remove_org_permissions" "$RESPONSE"
  echo "$PARSED_RESPONSE"
}

init_zitadel() {
  echo "Initializing Zitadel IDP"
  INSTANCE_URL="$NETBIRD_HTTP_PROTOCOL://$NETBIRD_DOMAIN:$NETBIRD_PORT"

  TOKEN_PATH=./machinekey/zitadel-admin-sa.token

  # shellcheck disable=SC2028
  echo -n "Waiting for Zitadel's PAT to be created "
  wait_pat "$TOKEN_PATH"
  echo "Reading Zitadel PAT"
  PAT=$(cat $TOKEN_PATH)
  if [ "$PAT" = "null" ]; then
    echo "Failed requesting getting Zitadel PAT"
    exit 1
  fi

  # shellcheck disable=SC2028
  echo -n "Waiting for Zitadel to become ready "
  wait_api "$INSTANCE_URL" "$PAT"

  #  create the zitadel project
  echo "Creating new zitadel project"
  PROJECT_ID=$(create_new_project "$INSTANCE_URL" "$PAT")

  ZITADEL_DEV_MODE=false
  BASE_REDIRECT_URL=$NETBIRD_HTTP_PROTOCOL://$NETBIRD_DOMAIN
  if [[ $NETBIRD_HTTP_PROTOCOL == "http" ]]; then
    ZITADEL_DEV_MODE=true
  fi

  # create zitadel spa application
  echo "Creating new zitadel spa application"
  APPLICATION_CLIENT_ID=$(create_new_application "$INSTANCE_URL" "$PAT")

  MACHINE_USER_ID=$(create_service_user "$INSTANCE_URL" "$PAT")

  SERVICE_USER_CLIENT_ID="null"
  SERVICE_USER_CLIENT_SECRET="null"

  create_service_user_secret "$INSTANCE_URL" "$PAT" "$MACHINE_USER_ID"

  DATE=$(add_organization_user_manager "$INSTANCE_URL" "$PAT" "$MACHINE_USER_ID")

  ZITADEL_ADMIN_USERNAME="admin@$NETBIRD_DOMAIN"
  ZITADEL_ADMIN_PASSWORD="$(openssl rand -base64 32 | sed 's/=//g')@"

  HUMAN_USER_ID=$(create_admin_user "$INSTANCE_URL" "$PAT" "$ZITADEL_ADMIN_USERNAME" "$ZITADEL_ADMIN_PASSWORD")

  DATE="null"

  DATE=$(add_instance_admin "$INSTANCE_URL" "$PAT" "$HUMAN_USER_ID")

  DATE="null"
  DATE=$(delete_auto_service_user "$INSTANCE_URL" "$PAT")
  if [ "$DATE" = "null" ]; then
      echo "Failed deleting auto service user"
      echo "Please remove it manually"
  fi

  export NETBIRD_AUTH_CLIENT_ID=$APPLICATION_CLIENT_ID
  export NETBIRD_IDP_MGMT_CLIENT_ID=$SERVICE_USER_CLIENT_ID
  export NETBIRD_IDP_MGMT_CLIENT_SECRET=$SERVICE_USER_CLIENT_SECRET
  export ZITADEL_ADMIN_USERNAME
  export ZITADEL_ADMIN_PASSWORD
}

initEnvironment() {
  CADDY_SECURE_DOMAIN=""
  ZITADEL_EXTERNALSECURE="false"
  ZITADEL_TLS_MODE="disabled"
  ZITADEL_MASTERKEY="$(openssl rand -base64 32 | head -c 32)"
  USING_DOMAIN="true"
  NETBIRD_PORT=80
  NETBIRD_HTTP_PROTOCOL="http"
  TURN_USER="self"
  TURN_PASSWORD=$(openssl rand -base64 32 | sed 's/=//g')
  TURN_MIN_PORT=49152
  TURN_MAX_PORT=65535

  NETBIRD_DOMAIN=$NETBIRD_DOMAIN
  if [ "$NETBIRD_DOMAIN-x" == "-x" ] ; then
    echo "NETBIRD_DOMAIN is not set, using the main IP address"
    NETBIRD_DOMAIN=$(get_main_ip_address)
    USING_DOMAIN="false"
  fi

  if [ "$NETBIRD_DOMAIN" == "localhost" ]; then
    USING_DOMAIN="false"
  fi

  if [ $USING_DOMAIN == "true" ]; then
    ZITADEL_EXTERNALSECURE="true"
    ZITADEL_TLS_MODE="external"
    NETBIRD_PORT=443
    CADDY_SECURE_DOMAIN=", $NETBIRD_DOMAIN:$NETBIRD_PORT"
    NETBIRD_HTTP_PROTOCOL="https"
  fi

  if [[ "$OSTYPE" == "darwin"* ]]; then
      ZIDATE_TOKEN_EXPIRATION_DATE=$(date -u -v+30M "+%Y-%m-%dT%H:%M:%SZ")
  else
      ZIDATE_TOKEN_EXPIRATION_DATE=$(date -u -d "+30 minutes" "+%Y-%m-%dT%H:%M:%SZ")
  fi

  DOCKER_COMPOSE_COMMAND=$(check_docker_compose)

  if [ -f zitadel.env ]; then
    echo "Generated files already exists, if you want to reinitialize the environment, please remove them first."
    echo "You can use the following commands:"
    echo "  $DOCKER_COMPOSE_COMMAND down --volumes # to remove all containers and volumes"
    echo "  rm -f docker-compose.yml Caddyfile zitadel.env dashboard.env machinekey/zitadel-admin-sa.token"
    echo "Be aware that this will remove all data from the database and you will have to reconfigure the dashboard."
    exit 1
  fi

  echo Rendering initial files...
  renderDockerCompose > docker-compose.yml
  renderCaddyfile > Caddyfile
  renderZitadelEnv > zitadel.env
  echo "" > dashboard.env

  mkdir -p machinekey
  chmod 777 machinekey

  init_crdb

  echo Starting zidatel
  $DOCKER_COMPOSE_COMMAND up -d caddy zitadel crdb
  init_zitadel

  echo Rendering NetBird files...
  renderTurnServerConf > turnserver.conf
  renderManagementJson > management.json
  renderDashboardEnv > dashboard.env

  echo Starting NetBird services
  $DOCKER_COMPOSE_COMMAND up -d
  echo "Done!"
  echo "You can access the NetBird dashboard at $NETBIRD_HTTP_PROTOCOL://$NETBIRD_DOMAIN:$NETBIRD_PORT"
  echo "Login with the following credentials:"
  echo "Username: $ZITADEL_ADMIN_USERNAME" | tee .env
  echo "Password: $ZITADEL_ADMIN_PASSWORD" | tee -a .env
}

renderCaddyfile() {
  cat <<EOF
{
  debug
	servers :80,:443 {
    protocols h1 h2c
  }
}

:80${CADDY_SECURE_DOMAIN} {
    # Signal
    reverse_proxy /signalexchange.SignalExchange/* h2c://signal:10000
    # Management
    reverse_proxy /api/* management:80
    reverse_proxy /management.ManagementService/* h2c://management:80
    # Zitadel
    reverse_proxy /zitadel.admin.v1.AdminService/* h2c://zitadel:8080
    reverse_proxy /admin/v1/* h2c://zitadel:8080
    reverse_proxy /zitadel.auth.v1.AuthService/* h2c://zitadel:8080
    reverse_proxy /auth/v1/* h2c://zitadel:8080
    reverse_proxy /zitadel.management.v1.ManagementService/* h2c://zitadel:8080
    reverse_proxy /management/v1/* h2c://zitadel:8080
    reverse_proxy /zitadel.system.v1.SystemService/* h2c://zitadel:8080
    reverse_proxy /system/v1/* h2c://zitadel:8080
    reverse_proxy /assets/v1/* h2c://zitadel:8080
    reverse_proxy /ui/* h2c://zitadel:8080
    reverse_proxy /oidc/v1/* h2c://zitadel:8080
    reverse_proxy /saml/v2/* h2c://zitadel:8080
    reverse_proxy /oauth/v2/* h2c://zitadel:8080
    reverse_proxy /.well-known/openid-configuration h2c://zitadel:8080
    reverse_proxy /openapi/* h2c://zitadel:8080
    reverse_proxy /debug/* h2c://zitadel:8080
    # Dashboard
    reverse_proxy /* dashboard:80
}
EOF
}

renderTurnServerConf() {
  cat <<EOF
listening-port=3478
tls-listening-port=5349
min-port=$TURN_MIN_PORT
max-port=$TURN_MAX_PORT
fingerprint
lt-cred-mech
user=$TURN_USER:$TURN_PASSWORD
realm=wiretrustee.com
cert=/etc/coturn/certs/cert.pem
pkey=/etc/coturn/private/privkey.pem
log-file=stdout
no-software-attribute
pidfile="/var/tmp/turnserver.pid"
no-cli
EOF
}

renderManagementJson() {
  cat <<EOF
{
    "Stuns": [
        {
            "Proto": "udp",
            "URI": "stun:$NETBIRD_DOMAIN:3478"
        }
    ],
    "TURNConfig": {
        "Turns": [
            {
                "Proto": "udp",
                "URI": "turn:$NETBIRD_DOMAIN:3478",
                "Username": "$TURN_USER",
                "Password": "$TURN_PASSWORD"
            }
        ],
        "TimeBasedCredentials": false
    },
    "Signal": {
        "Proto": "$NETBIRD_HTTP_PROTOCOL",
        "URI": "$NETBIRD_DOMAIN:$NETBIRD_PORT"
    },
    "HttpConfig": {
        "AuthIssuer": "$NETBIRD_HTTP_PROTOCOL://$NETBIRD_DOMAIN",
        "AuthAudience": "$NETBIRD_AUTH_CLIENT_ID",
        "OIDCConfigEndpoint":"$NETBIRD_HTTP_PROTOCOL://$NETBIRD_DOMAIN/.well-known/openid-configuration"
    },
    "IdpManagerConfig": {
        "ManagerType": "zitadel",
        "ClientConfig": {
            "Issuer": "$NETBIRD_HTTP_PROTOCOL://$NETBIRD_DOMAIN:$NETBIRD_PORT",
            "TokenEndpoint": "$NETBIRD_HTTP_PROTOCOL://$NETBIRD_DOMAIN:$NETBIRD_PORT/oauth/v2/token",
            "ClientID": "$NETBIRD_IDP_MGMT_CLIENT_ID",
            "ClientSecret": "$NETBIRD_IDP_MGMT_CLIENT_SECRET",
            "GrantType": "client_credentials"
        },
        "ExtraConfig": {
            "ManagementEndpoint": "$NETBIRD_HTTP_PROTOCOL://$NETBIRD_DOMAIN:$NETBIRD_PORT/management/v1"
        }
     },
    "PKCEAuthorizationFlow": {
        "ProviderConfig": {
            "Audience": "$NETBIRD_AUTH_CLIENT_ID",
            "ClientID": "$NETBIRD_AUTH_CLIENT_ID",
            "Scope": "openid profile email offline_access",
            "RedirectURLs": ["http://localhost:53000/","http://localhost:54000/"]
        }
    }
}
EOF
}

renderDashboardEnv() {
  cat <<EOF
# Endpoints
NETBIRD_MGMT_API_ENDPOINT=$NETBIRD_HTTP_PROTOCOL://$NETBIRD_DOMAIN:$NETBIRD_PORT
NETBIRD_MGMT_GRPC_API_ENDPOINT=$NETBIRD_HTTP_PROTOCOL://$NETBIRD_DOMAIN:$NETBIRD_PORT
# OIDC
AUTH_AUDIENCE=$NETBIRD_AUTH_CLIENT_ID
AUTH_CLIENT_ID=$NETBIRD_AUTH_CLIENT_ID
AUTH_AUTHORITY=$NETBIRD_HTTP_PROTOCOL://$NETBIRD_DOMAIN:$NETBIRD_PORT
USE_AUTH0=false
AUTH_SUPPORTED_SCOPES="openid profile email offline_access"
AUTH_REDIRECT_URI=/nb-auth
AUTH_SILENT_REDIRECT_URI=/nb-silent-auth
# SSL
NGINX_SSL_PORT=443
# Letsencrypt
LETSENCRYPT_DOMAIN=none
EOF
}

renderZitadelEnv() {
  cat <<EOF
ZITADEL_LOG_LEVEL=debug
ZITADEL_MASTERKEY=$ZITADEL_MASTERKEY
ZITADEL_DATABASE_COCKROACH_HOST=crdb
ZITADEL_DATABASE_COCKROACH_USER_USERNAME=zitadel_user
ZITADEL_DATABASE_COCKROACH_USER_SSL_MODE=verify-full
ZITADEL_DATABASE_COCKROACH_USER_SSL_ROOTCERT="/crdb-certs/ca.crt"
ZITADEL_DATABASE_COCKROACH_USER_SSL_CERT="/crdb-certs/client.zitadel_user.crt"
ZITADEL_DATABASE_COCKROACH_USER_SSL_KEY="/crdb-certs/client.zitadel_user.key"
ZITADEL_DATABASE_COCKROACH_ADMIN_SSL_MODE=verify-full
ZITADEL_DATABASE_COCKROACH_ADMIN_SSL_ROOTCERT="/crdb-certs/ca.crt"
ZITADEL_DATABASE_COCKROACH_ADMIN_SSL_CERT="/crdb-certs/client.root.crt"
ZITADEL_DATABASE_COCKROACH_ADMIN_SSL_KEY="/crdb-certs/client.root.key"
ZITADEL_EXTERNALSECURE=$ZITADEL_EXTERNALSECURE
ZITADEL_TLS_ENABLED="false"
ZITADEL_EXTERNALPORT=$NETBIRD_PORT
ZITADEL_EXTERNALDOMAIN=$NETBIRD_DOMAIN
ZITADEL_FIRSTINSTANCE_PATPATH=/machinekey/zitadel-admin-sa.token
ZITADEL_FIRSTINSTANCE_ORG_MACHINE_MACHINE_USERNAME=zitadel-admin-sa
ZITADEL_FIRSTINSTANCE_ORG_MACHINE_MACHINE_NAME=Admin
ZITADEL_FIRSTINSTANCE_ORG_MACHINE_PAT_SCOPES=openid
ZITADEL_FIRSTINSTANCE_ORG_MACHINE_PAT_EXPIRATIONDATE=$ZIDATE_TOKEN_EXPIRATION_DATE
EOF
}

renderDockerCompose() {
  cat <<EOF
version: "3.4"
services:
  # Caddy reverse proxy
  caddy:
    image: caddy
    restart: unless-stopped
    networks: [ netbird ]
    ports:
      - '443:443'
      - '80:80'
      - '8080:8080'
    volumes:
      - netbird_caddy_data:/data
      - ./Caddyfile:/etc/caddy/Caddyfile
  #UI dashboard
  dashboard:
    image: wiretrustee/dashboard:latest
    restart: unless-stopped
    networks: [netbird]
    env_file:
      - ./dashboard.env
  # Signal
  signal:
    image: netbirdio/signal:latest
    restart: unless-stopped
    networks: [netbird]
  # Management
  management:
    image: netbirdio/management:latest
    restart: unless-stopped
    networks: [netbird]
    volumes:
      - netbird_management:/var/lib/netbird
      - ./management.json:/etc/netbird/management.json
    command: [
      "--port", "80",
      "--log-file", "console",
      "--log-level", "debug",
      "--disable-anonymous-metrics=false",
      "--single-account-mode-domain=netbird.selfhosted",
      "--dns-domain=netbird.selfhosted",
    ]
  # Coturn, AKA relay server
  coturn:
    image: coturn/coturn
    restart: unless-stopped
    domainname: netbird.relay.selfhosted
    volumes:
      - ./turnserver.conf:/etc/turnserver.conf:ro
    network_mode: host
    command:
      - -c /etc/turnserver.conf
  # Zitadel - identity provider
  zitadel:
    restart: 'always'
    networks: [netbird]
    image: 'ghcr.io/zitadel/zitadel:v2.31.3'
    command: 'start-from-init --masterkeyFromEnv --tlsMode $ZITADEL_TLS_MODE'
    env_file:
      - ./zitadel.env
    depends_on:
      crdb:
        condition: 'service_healthy'
    volumes:
      - ./machinekey:/machinekey
      - netbird_zitadel_certs:/crdb-certs:ro
    healthcheck:
      test: [ "CMD", "curl", "-f", "http://localhost:8080/debug/healthz" ]
      interval: '10s'
      timeout: '30s'
      retries: 5
      start_period: '20s'
  # CockroachDB for zitadel
  crdb:
    restart: 'always'
    networks: [netbird]
    image: 'cockroachdb/cockroach:v22.2.2'
    command: 'start-single-node --advertise-addr crdb'
    volumes:
      - netbird_crdb_data:/cockroach/cockroach-data
      - netbird_crdb_certs:/cockroach/certs
      - netbird_zitadel_certs:/zitadel-certs
    healthcheck:
      test: [ "CMD", "curl", "-f", "http://localhost:8080/health?ready=1" ]
      interval: '10s'
      timeout: '30s'
      retries: 5
      start_period: '20s'

volumes:
  netbird_management:
  netbird_caddy_data:
  netbird_crdb_data:
  netbird_crdb_certs:
  netbird_zitadel_certs:

networks:
  netbird:
EOF
}

initEnvironment
