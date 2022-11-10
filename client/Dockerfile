FROM node:13.12.0-alpine
WORKDIR /app

ENV PATH /app/node_modules/.bin:$PATH

ENV S3_BUCKET_URL https://spotify-views.s3.us-east-2.amazonaws.com/

# install app dependencies
COPY package.json ./
#COPY package-lock.json ./
RUN npm install --silent
RUN npm install react-scripts@3.4.1 -g --silent
#idk why this isn't working
RUN npm install axios --save


COPY . ./

CMD ["npm", "start"]