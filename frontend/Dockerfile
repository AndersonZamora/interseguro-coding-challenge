# Este Dockerfile es solo para pruebas locales con docker-compose. En Render
# el frontend se despliega como Static Site nativo (build: npm run build,
# publishDir: dist), sin Docker — ver render.yaml.
FROM node:22-alpine AS build
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
ARG VITE_GO_API_URL=http://localhost:8080
ARG VITE_NODE_API_URL=http://localhost:3000
ENV VITE_GO_API_URL=$VITE_GO_API_URL
ENV VITE_NODE_API_URL=$VITE_NODE_API_URL
RUN npm run build

FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80
