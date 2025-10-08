# Stage 1: Build backend
FROM golang:alpine AS backend-builder

WORKDIR /app
COPY woerterbuch-backend/src /app/src
COPY woerterbuch-backend/go.mod /app
COPY woerterbuch-backend/main.go /app

RUN go mod tidy && go build -o woerterbuch

# Stage 2: Build frontend
FROM node:18-alpine AS frontend-builder

WORKDIR /app
COPY woerterbuch-frontend/package.json woerterbuch-frontend/package-lock.json ./
RUN npm install

COPY woerterbuch-frontend/public public
COPY woerterbuch-frontend/src src
RUN npm run build

# Stage 3: Final image
FROM nginx:alpine

COPY nginx/config/nginx.conf /etc/nginx/
COPY --from=backend-builder /app/woerterbuch /usr/bin/woerterbuch
COPY --from=frontend-builder /app/build /usr/share/nginx/html

RUN addgroup -S wb_user_group && \
    adduser -S wb_user -G wb_user_group && \
    chown -R wb_user:wb_user_group /usr/bin/woerterbuch
USER wb_user

EXPOSE 80

ENTRYPOINT ["sh", "-c", "/app/woerterbuch & nginx -g 'daemon off;'"]
