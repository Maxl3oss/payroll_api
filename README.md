# ENV
```yaml
# Stage status to start server:
#   - "dev", for start server without graceful shutdown
#   - "prod", for start server with graceful shutdown
STAGE_STATUS="dev"

# Server settings:
SERVER_HOST="0.0.0.0"
PORT=[Port]
SERVER_READ_TIMEOUT=60

# JWT settings:
JWT_SECRET_KEY=[TokenSecret]
JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT=15
JWT_REFRESH_KEY=[TokenRefesh]
JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT=720
SECRET_KEY=[secret]

# Database settings:
CONNECT_TYPE="mysql"
DB_USER=[User]
DB_PASSWORD=[Password]
DB_HOST=[Host]
DB_PORT=[Port]
DB_NAME=[Name]
DB_SSL_MODE="disable"

# setup sender email
SEND_MAIL_USER=[EmailSender]
SEND_MAIL_PASS=[EmailPassword]
```
