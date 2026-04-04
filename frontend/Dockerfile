FROM node:22-alpine AS builder
ARG GIT_COMMIT="unknown"
ENV VITE_GIT_COMMIT=${GIT_COMMIT}
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
# Prepare directories and pid for non-root execution
RUN mkdir -p /var/cache/nginx /var/log/nginx && \
    chown -R 10001:10001 /var/cache/nginx /var/log/nginx /etc/nginx/conf.d && \
    sed -i 's|pid\s*/run/nginx.pid;|pid /tmp/nginx.pid;|' /etc/nginx/nginx.conf
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 8080
