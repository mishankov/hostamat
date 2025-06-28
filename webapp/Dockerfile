# Stage 1: Build the SvelteKit application
FROM node:latest AS build

WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build

# Stage 2: Create the production image
FROM node:latest AS final

WORKDIR /app

# Copy only necessary files from the build stage
COPY --from=build /app/package.json ./
COPY --from=build /app/node_modules/ ./
COPY --from=build /app/build/ ./

EXPOSE 3000

CMD ["node", "build"]