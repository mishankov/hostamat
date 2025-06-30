FROM node:latest AS build
WORKDIR /app
COPY . .
RUN npm ci && npm run build
CMD ["node", "build"]
