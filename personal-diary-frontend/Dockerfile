# Stage 1: Build the app
FROM node:20 AS build

WORKDIR /app
COPY . .

RUN npm install
RUN npm run build

# Stage 2: Serve using Nginx
FROM nginx:alpine

# Copy built files to Nginx's html directory
COPY --from=build /app/dist /usr/share/nginx/html

# Optional: custom nginx config (see below)
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
