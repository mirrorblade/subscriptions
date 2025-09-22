set dotenv-load
set shell := ["zsh", "-cu"]

SUCCESS_MESSAGE := "[\\u001b[32mSUCCESS\\u001b[0m]"
INFO_MESSAGE := "[\\u001b[36mINFO\\u001b[0m]"

DATABASE_URL := "postgresql://$DATABASE_USER:$DATABASE_PASSWORD@$DATABASE_HOST:$DATABASE_PORT/$DATABASE_NAME?x-multi-statement=true&sslmode=disable"
MIGRATION_PATH := "./migrations/"

COMPOSE_FILE := "./docker-compose.yml"

# list all available commands of just
default:
    @just --list --unsorted

# start the entire stack
@up *ARGS:
    docker-compose -f {{COMPOSE_FILE}} up {{ARGS}}
    just migrate-up

# stop everything gracefully
@down *ARGS:
    docker-compose -f {{COMPOSE_FILE}} down {{ARGS}}

# start existing containers
@start:
    docker-compose -f {{COMPOSE_FILE}} start

# stop running containers
@stop:
    docker-compose -f {{COMPOSE_FILE}} stop

# restart running containers
@restart:
    docker-compose -f {{COMPOSE_FILE}} restart

# rebuild with force
@rebuild:
    just up --build --force-recreate -d

# run swagger-ui container on port 8080
@run-swagger:
    docker run -d --name=swagger-ui --network=host -e SWAGGER_JSON=/api/openapi.yaml -v ./api/openapi.yaml:/api/openapi.yaml swaggerapi/swagger-ui

# start swagger-ui container
@start-swagger:
    docker start swagger-ui

# stop swagger-ui container
@stop-swagger:
    docker stop swagger-ui

# kill swagger-ui container
@kill-swagger:
    docker stop swagger-ui

# remove swagger-ui container
@remove-swagger:
    just stop-swagger
    docker rm swagger-ui

# restart swagger-ui container
@restart-swagger:
    docker restart swagger-ui

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
