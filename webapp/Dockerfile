FROM oven/bun:1 AS build
WORKDIR /app
COPY . .
RUN bun install && bun run --bun build

FROM oven/bun:1
WORKDIR /app
COPY --from=build /app/node_modules node_modules
COPY --from=build /app/build build
COPY --from=build /app/package.json .
EXPOSE 3000/tcp
CMD ["bun", "run", "./build"]
