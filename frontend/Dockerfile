#nginx/Dockerfile
FROM node:13.12.0-alpine as build
WORKDIR /app
ENV PATH /app/node_modules/.bin:$PATH
COPY package.json ./
COPY package-lock.json ./
RUN npm ci
RUN npm install react-scripts@3.4.1 -g
COPY . ./
RUN npm run build

FROM nginx:latest
COPY --from=build /app/build /usr/share/nginx/html
# COPY ./nginx/mime.types /etc/nginx/mime.types
EXPOSE 80
# COPY nginx/nginx.conf /etc/nginx/nginx.conf
COPY nginx/ /etc/nginx/
CMD ["nginx", "-g", "daemon off;"]