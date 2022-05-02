# POSTGRES_TO_GOMIGRATE

Simple script to prepare [gomigrate](https://github.com/golang-migrate/migrate) migration files from existing postgres database


1. Setup Postgres (we setup using docker-compose)
    ```bash
    # start postgres
    docker compose up -d

    # stop and clean postgres 
    docker compose down -v

    # execute initial database script
    PGPASSWORD=pass psql -h localhost -p 5434 -U user -d user -f initial_database.sql

    # postgres cli
    PGPASSWORD=pass psql -h localhost -p 5434 -U user -d user 
    ```

2. Generate migrations file (edit the script for different configs)
    ```bash
    # remove the previous generated files
    rm -rf migrations

    # run the script
    go run generate_migration.go
    ```

3. Test for GoMigrate (don't forget to clean your database before)
    ```bash
    # install gomigrate
    brew install golang-migrate

    # migrate database
    migrate -path migrations -database "postgresql://user:pass@localhost:5434/user?sslmode=disable" up

    # rollback database
    migrate -path migrations -database "postgresql://user:pass@localhost:5434/user?sslmode=disable" down
    ```


