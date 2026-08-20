#!/usr/bin/env sh

set -eu

for migration in /migrations/*.sql; do
    echo "Applying ${migration}"
    sed -n '/^-- +goose Up$/,/^-- +goose Down$/p' "${migration}" \
        | sed '$d' \
        | psql \
            --set ON_ERROR_STOP=on \
            --username "${POSTGRES_USER}" \
            --dbname "${POSTGRES_DB}"
done

touch /var/lib/postgresql/data/.migrations-complete
