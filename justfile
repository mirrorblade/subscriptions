set dotenv-load
set shell := ["zsh", "-cu"]

SUCCESS_MESSAGE := "[\\u001b[32mSUCCESS\\u001b[0m]"
INFO_MESSAGE := "[\\u001b[36mINFO\\u001b[0m]"

DATABASE_URL := "postgresql://$PGUSER:$PGPASSWORD@$PGHOST:$PGPORT/postgres?x-multi-statement=true"
MIGRATION_PATH := "./migrations/"

# list all available commands of just
default:
    @just --list --unsorted

# migrate database up
@migrate-up N="":
    echo "{{INFO_MESSAGE}} Starts migrate up database" 
    migrate -verbose -path "{{MIGRATION_PATH}}" -database "{{DATABASE_URL}}" up {{N}}
    echo "{{SUCCESS_MESSAGE}} Migration up was successful" 

# create new migration
@migrate-create MIGRATION_NAME:
    echo "{{INFO_MESSAGE}} Starts creating new migration" 
    migrate create -ext sql -dir "{{MIGRATION_PATH}}" {{MIGRATION_NAME}}
    echo "{{SUCCESS_MESSAGE}} Creating migrtaion was successful" 

# force database version
@migrate-force VERSION:
    echo "{{INFO_MESSAGE}} Starts forcing migration" 
    migrate -verbose -path "{{MIGRATION_PATH}}" -database "{{DATABASE_URL}}" force {{VERSION}}
    echo "{{SUCCESS_MESSAGE}} Forcing migration was successful" 

# fix latest dirty migration 
@migrate-fix:
    echo "{{INFO_MESSAGE}} Starts fixing latest dirty migration" && \
    VERSION_OUTPUT=$(migrate -path "{{MIGRATION_PATH}}" -database "{{DATABASE_URL}}" version 2>&1) && \
    if [[ "$VERSION_OUTPUT" == *"(dirty)"* ]]; then \
        VERSION=$(echo "$VERSION_OUTPUT" | cut -d ' ' -f 1) && \
        just migrate-force $VERSION && \
        just migrate-down 1 && \
        echo "{{SUCCESS_MESSAGE}} Fixing latest migration was successful" && \
    else \
        echo "{{SUCCESS_MESSAGE}} Your database is already clear" && \
    fi 

# migrate database down
@migrate-down N="-all":
    echo "{{INFO_MESSAGE}} Starts migrate down database" 
    migrate -verbose -path "{{MIGRATION_PATH}}" -database "{{DATABASE_URL}}" down {{N}}
    echo "{{SUCCESS_MESSAGE}} Migration down was successful" 
