# Dockerfile for frontend application
FROM node:19.6.0 as frontend-build

WORKDIR /app

COPY frontend/package*.json ./
RUN npm install

COPY frontend .
RUN npm run build

# Dockerfile for backend application
FROM golang:1.20.3 as backend-build

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend .
RUN go build -o main .

# Dockerfile for nginx
FROM nginx:1.21.0

COPY --from=frontend-build /app/dist /usr/share/nginx/html
COPY nginx/nginx.conf /etc/nginx/conf.d/default.conf

# Expose ports
EXPOSE 3000
EXPOSE 8080

# Start nginx and backend server
CMD service nginx start && ./backend/main
