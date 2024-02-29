

APP_NAME := "hyprstfu"

# Build the application.
build:
    @go build -o ./build/{{ APP_NAME }} .

# Run the application after a specified sleep time.
run SLEEP_TIME='2': build
    @echo "waiting {{ SLEEP_TIME }} seconds... please put your cursor in a valid window"
    @sleep {{ SLEEP_TIME }}
    @./build/{{ APP_NAME }}
