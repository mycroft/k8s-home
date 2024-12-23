#!/bin/sh

DRY=

NS=temporal

ADMIN_POD=$(kubectl get -n ${NS} pods -l app.kubernetes.io/component=admintools -o name)

POSTGRES_SEEDS=postgres-instance.postgres
DB_PORT=5432

DBNAME=temporal
POSTGRES_USER=$(kubectl get -n ${NS} secret -o yaml postgresql -o json | jq -r .data.username | base64 -d)
POSTGRES_PWD=$(kubectl get -n ${NS} secret -o yaml postgresql -o json | jq -r .data.password | base64 -d)
kubectl exec -n ${NS} -ti ${ADMIN_POD} -- ${DRY} temporal-sql-tool --plugin postgres12 --ep "${POSTGRES_SEEDS}" -u "${POSTGRES_USER}" -pw "${POSTGRES_PWD}" -p "${DB_PORT}" --db "${DBNAME}" update-schema -d "/etc/temporal/schema/postgresql/v12/temporal/versioned"

DBNAME=temporal_visibility
POSTGRES_USER=$(kubectl get -n ${NS} secret -o yaml postgresql_visibility -o json | jq -r .data.username | base64 -d)
POSTGRES_PWD=$(kubectl get -n ${NS} secret -o yaml postgresql_visibility -o json | jq -r .data.password | base64 -d)

kubectl exec -ti -n ${NS} ${ADMIN_POD} -- ${DRY} temporal-sql-tool --plugin postgres12 --ep "${POSTGRES_SEEDS}" -u "${POSTGRES_USER}" -pw "${POSTGRES_PWD}" -p "${DB_PORT}" --db "${DBNAME}" update-schema -d "/etc/temporal/schema/postgresql/v12/visibility/versioned"
