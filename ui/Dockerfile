FROM node:14-alpine as build-stage
WORKDIR /app
COPY ./package.json ./package-lock.json ./
RUN npm config set depth 0 \
	&& npm cache clean --force \
	&& npm i
COPY . .

# ENTRYPOINT [ "ng", "serve" ]

# ENTRYPOINT [ "ash" ]
RUN ./node_modules/.bin/ng build --prod=true
# RUN ls -lahg ./dist

# ng build my-app -c production

# RUN ./node_modules/.bin/ng build --env=prod

FROM nginx:1.17-alpine
COPY nginx/default.conf /etc/nginx/conf.d/
RUN rm -rf /usr/share/nginx/html/*
COPY --from=build-stage /app/dist/todo-app/ /usr/share/nginx/html
RUN ls -lahg /usr/share/nginx/html
CMD ["nginx", "-g", "daemon off;"]
