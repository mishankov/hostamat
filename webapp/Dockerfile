FROM node:latest AS build
WORKDIR /app
COPY . .
RUN npm ci && npm run build
CMD ["BODY_SIZE_LIMIT=Infinity", "node", "build"]
