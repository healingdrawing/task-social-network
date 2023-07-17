# Stage 1: Build Vue application
FROM node:19.6.0 as build-stage

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

RUN npm run build

# Stage 2: Serve Vue application with Nginx
FROM nginx:1.21.0 as production-stage

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build-stage /app/dist /usr/share/nginx/html
COPY --from=build-stage /app /usr/share/nginx

EXPOSE 3000

CMD ["nginx", "-g", "daemon off;"]
