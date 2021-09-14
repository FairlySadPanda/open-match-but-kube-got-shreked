#!/bin/bash

# cp /opt/bitnami/redis-sentinel/mounted-etc/original_sentinel.conf /opt/bitnami/redis-sentinel/etc/sentinel.conf

# replace_in_file() {
#     local filename="${1:?filename is required}"
#     local match_regex="${2:?match regex is required}"
#     local substitute_regex="${3:?substitute regex is required}"
#     local posix_regex=${4:-true}

#     local result

#     # We should avoid using 'sed in-place' substitutions
#     # 1) They are not compatible with files mounted from ConfigMap(s)
#     # 2) We found incompatibility issues with Debian10 and "in-place" substitutions
#     del=$'\001' # Use a non-printable character as a 'sed' delimiter to avoid issues
#     if [[ $posix_regex = true ]]; then
#         result="$(sed -E "s${del}${match_regex}${del}${substitute_regex}${del}g" "$filename")"
#     else
#         result="$(sed "s${del}${match_regex}${del}${substitute_regex}${del}g" "$filename")"
#     fi
#     echo "$result" > "$filename"
# }
# sentinel_conf_set() {
#     local -r key="${1:?missing key}"
#     local value="${2:-}"

#     # Sanitize inputs
#     value="${value//\\/\\\\}"
#     value="${value//&/\\&}"
#     value="${value//\?/\\?}"
#     [[ "$value" = "" ]] && value="\"$value\""

#     replace_in_file "/opt/bitnami/redis-sentinel/etc/sentinel.conf" "^#*\s*${key} .*" "${key} ${value}" false
# }
# sentinel_conf_add() {
#     echo $'\n'"$@" >> "/opt/bitnami/redis-sentinel/etc/sentinel.conf"
# }
# is_boolean_yes() {
#     local -r bool="${1:-}"
#     # comparison is performed without regard to the case of alphabetic characters
#     shopt -s nocasematch
#     if [[ "$bool" = 1 || "$bool" =~ ^(yes|true)$ ]]; then
#         true
#     else
#         false
#     fi
# }
# host_id() {
#     echo "$1" | openssl sha1 | awk '{print $2}'
# }

# HEADLESS_SERVICE="open-match-redis-headless.open-match.svc.cluster.local"
# REDIS_SERVICE="open-match-redis.open-match.svc.cluster.local"

# if [[ -n $REDIS_PASSWORD_FILE ]]; then
#     password_aux=`cat ${REDIS_PASSWORD_FILE}`
#     export REDIS_PASSWORD=$password_aux
# fi

# if [[ ! -f /opt/bitnami/redis-sentinel/etc/sentinel.conf ]]; then
#     cp /opt/bitnami/redis-sentinel/mounted-etc/sentinel.conf /opt/bitnami/redis-sentinel/etc/sentinel.conf
# fi

# export REDIS_REPLICATION_MODE="slave"
# if [[ -z "$(getent ahosts "$HEADLESS_SERVICE" | grep -v "^$(hostname -i) ")" ]]; then
#     export REDIS_REPLICATION_MODE="master"
# fi

# if [[ "$REDIS_REPLICATION_MODE" == "master" ]]; then
#     REDIS_MASTER_HOST="$(hostname -i)"
#     REDIS_MASTER_PORT_NUMBER="6379"
# else
#     if is_boolean_yes "$REDIS_SENTINEL_TLS_ENABLED"; then
#     sentinel_info_command="redis-cli -h $REDIS_SERVICE -p 26379 --tls --cert ${REDIS_SENTINEL_TLS_CERT_FILE} --key ${REDIS_SENTINEL_TLS_KEY_FILE} --cacert ${REDIS_SENTINEL_TLS_CA_FILE} sentinel get-master-addr-by-name om-redis-master"
#     else
#     sentinel_info_command="redis-cli -h $REDIS_SERVICE -p 26379 sentinel get-master-addr-by-name om-redis-master"
#     fi
#     REDIS_SENTINEL_INFO=($($sentinel_info_command))
#     REDIS_MASTER_HOST=${REDIS_SENTINEL_INFO[0]}
#     REDIS_MASTER_PORT_NUMBER=${REDIS_SENTINEL_INFO[1]}

#     # Immediately attempt to connect to the reported master. If it doesn't exist the connection attempt will either hang
#     # or fail with "port unreachable" and give no data. The liveness check will then timeout waiting for the sentinel
#     # container to be ready and restart the it. By then the new master will likely have been elected
#     if is_boolean_yes "$REDIS_SENTINEL_TLS_ENABLED"; then
#     sentinel_info_command="redis-cli -h $REDIS_MASTER_HOST -p 26379 --tls --cert ${REDIS_SENTINEL_TLS_CERT_FILE} --key ${REDIS_SENTINEL_TLS_KEY_FILE} --cacert ${REDIS_SENTINEL_TLS_CA_FILE} sentinel get-master-addr-by-name om-redis-master"
#     else
#     sentinel_info_command="redis-cli -h $REDIS_MASTER_HOST -p 26379 sentinel get-master-addr-by-name om-redis-master"
#     fi

#     if [[ ! ($($sentinel_info_command)) ]]; then
#     # master doesn't actually exist, this probably means the remaining pods haven't elected a new one yet
#     # and are reporting the old one still. Once this happens the container will get stuck and never see the new
#     # master. We stop here to allow the container to not pass the liveness check and be restarted.
#     exit 1
#     fi
# fi
# sentinel_conf_set "sentinel monitor" "om-redis-master "$REDIS_MASTER_HOST" "$REDIS_MASTER_PORT_NUMBER" 2"

# add_replica() {
#     if [[ "$1" != "$REDIS_MASTER_HOST" ]]; then
#     sentinel_conf_add "sentinel known-replica om-redis-master $1 6379"
#     fi
# }

cp /opt/bitnami/redis-sentinel/mounted-etc/sentinel.conf /opt/bitnami/redis-sentinel/etc/sentinel.conf
exec redis-server /opt/bitnami/redis-sentinel/etc/sentinel.conf --sentinel