FROM node:18-alpine

WORKDIR /app

COPY package.json package-lock.json ./

RUN npm install

COPY public public

EXPOSE 3000

ENTRYPOINT ["npm", "start"]