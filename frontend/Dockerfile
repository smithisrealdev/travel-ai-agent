# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Production stage
FROM node:20-alpine

WORKDIR /app

# Copy built application from builder
COPY --from=builder /app/.output /app/.output
COPY --from=builder /app/package*.json ./

# Install production dependencies only
RUN npm ci --only=production

# Expose port
EXPOSE 3000

# Set environment variable
ENV NODE_ENV=production
ENV HOST=0.0.0.0

# Start the application
CMD ["node", ".output/server/index.mjs"]
