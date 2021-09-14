#!/bin/bash
is_boolean_yes() {
    local -r bool="${1:-}"
    # comparison is performed without regard to the case of alphabetic characters
    shopt -s nocasematch
    if [[ "$bool" = 1 || "$bool" =~ ^(yes|true)$ ]]; then
        true
    else
        false
    fi
}

HEADLESS_SERVICE="open-match-redis-headless.open-match.svc.cluster.local"
REDIS_SERVICE="open-match-redis.open-match.svc.cluster.local"

export REDIS_REPLICATION_MODE="slave"
if [[ -z "$(getent ahosts "$HEADLESS_SERVICE" | grep -v "^$(hostname -i) ")" ]]; then
    export REDIS_REPLICATION_MODE="master"
fi

if [[ -n $REDIS_PASSWORD_FILE ]]; then
    password_aux=`cat ${REDIS_PASSWORD_FILE}`
    export REDIS_PASSWORD=$password_aux
fi

if [[ -n $REDIS_MASTER_PASSWORD_FILE ]]; then
    password_aux=`cat ${REDIS_MASTER_PASSWORD_FILE}`
    export REDIS_MASTER_PASSWORD=$password_aux
fi

if [[ "$REDIS_REPLICATION_MODE" == "master" ]]; then
    echo "I am master"
    if [[ ! -f /opt/bitnami/redis/etc/master.conf ]];then
    cp /opt/bitnami/redis/mounted-etc/master.conf /opt/bitnami/redis/etc/master.conf
    fi
else
    if [[ ! -f /opt/bitnami/redis/etc/replica.conf ]];then
    cp /opt/bitnami/redis/mounted-etc/replica.conf /opt/bitnami/redis/etc/replica.conf
    fi

    if is_boolean_yes "$REDIS_TLS_ENABLED"; then
    sentinel_info_command="redis-cli -h $REDIS_SERVICE -p 26379 --tls --cert ${REDIS_TLS_CERT_FILE} --key ${REDIS_TLS_KEY_FILE} --cacert ${REDIS_TLS_CA_FILE} sentinel get-master-addr-by-name om-redis-master"
    else
    sentinel_info_command="redis-cli -h $REDIS_SERVICE -p 26379 sentinel get-master-addr-by-name om-redis-master"
    fi
    REDIS_SENTINEL_INFO=($($sentinel_info_command))
    REDIS_MASTER_HOST=${REDIS_SENTINEL_INFO[0]}
    REDIS_MASTER_PORT_NUMBER=${REDIS_SENTINEL_INFO[1]}


    # Immediately attempt to connect to the reported master. If it doesn't exist the connection attempt will either hang
    # or fail with "port unreachable" and give no data. The liveness check will then timeout waiting for the redis
    # container to be ready and restart the it. By then the new master will likely have been elected
    if is_boolean_yes "$REDIS_TLS_ENABLED"; then
    sentinel_info_command="redis-cli -h $REDIS_MASTER_HOST -p 26379 --tls --cert ${REDIS_TLS_CERT_FILE} --key ${REDIS_TLS_KEY_FILE} --cacert ${REDIS_TLS_CA_FILE} sentinel get-master-addr-by-name om-redis-master"
    else
    sentinel_info_command="redis-cli -h $REDIS_MASTER_HOST -p 26379 sentinel get-master-addr-by-name om-redis-master"
    fi

    if [[ ! ($($sentinel_info_command)) ]]; then
    # master doesn't actually exist, this probably means the remaining pods haven't elected a new one yet
    # and are reporting the old one still. Once this happens the container will get stuck and never see the new
    # master. We stop here to allow the container to not pass the liveness check and be restarted.
    exit 1
    fi
fi

if [[ ! -f /opt/bitnami/redis/etc/redis.conf ]];then
    cp /opt/bitnami/redis/mounted-etc/redis.conf /opt/bitnami/redis/etc/redis.conf
fi
ARGS=("--port" "${REDIS_PORT}")

if [[ "$REDIS_REPLICATION_MODE" == "slave" ]]; then
    ARGS+=("--slaveof" "${REDIS_MASTER_HOST}" "${REDIS_MASTER_PORT_NUMBER}")
fi
ARGS+=("--protected-mode" "no")

if [[ "$REDIS_REPLICATION_MODE" == "master" ]]; then
    ARGS+=("--include" "/opt/bitnami/redis/etc/master.conf")
else
    ARGS+=("--include" "/opt/bitnami/redis/etc/replica.conf")
fi

ARGS+=("--include" "/opt/bitnami/redis/etc/redis.conf")
exec /run.sh "${ARGS[@]}"