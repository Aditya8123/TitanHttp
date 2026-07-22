# Builder Stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copy dependency definition
COPY package.json package-lock.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Compile TypeScript and build React app
RUN npm run build

# Production Stage
FROM nginx:alpine

# Copy custom Nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy static assets from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Create runtime directory structures, touch PID file, and grant access to non-root nginx user
RUN touch /var/run/nginx.pid && \
    chown -R nginx:nginx /var/run/nginx.pid /var/cache/nginx /var/log/nginx /usr/share/nginx/html

# Switch to the non-root user
USER nginx

# Expose non-privileged port
EXPOSE 8080

CMD ["nginx", "-g", "daemon off;"]
