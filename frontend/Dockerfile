# stage1 as builder
FROM node:latest as builder

WORKDIR /vue-ui

# Copy the package.json and install dependencies
COPY package*.json ./
RUN npm install

# Copy rest of the files
COPY . .

# Build the project
RUN npm run build

FROM nginx:alpine as production-build
COPY ./.nginx/nginx.conf /etc/nginx/nginx.conf

## Remove default nginx index page
RUN rm -rf /usr/share/nginx/html/*

# Copy from the stahg 1
COPY --from=builder /vue-ui/dist /usr/share/nginx/html

EXPOSE 8080
ENTRYPOINT ["nginx", "-g", "daemon off;"]

# First stage
# FROM node as builder

# # Set the working directory
# WORKDIR /frontend

# # Copy the application files to the working directory
# COPY . /frontend/

# # Install dependencies
# RUN npm install

# # Build the application
# RUN npm run serve

# # Second stage
# FROM nginx:alpine

# # Copy the built application files from the first stage to the nginx directory
# COPY --from=builder /frontend/dist /usr/share/nginx/html

# # Expose port 8080
# EXPOSE 8080

# # Start the nginx server
# CMD ["nginx", "-g", "daemon off;"]
