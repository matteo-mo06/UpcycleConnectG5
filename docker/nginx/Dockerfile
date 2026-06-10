FROM node:20-alpine AS frontend
WORKDIR /app
COPY frontend/upcycleconnect/package*.json ./
RUN npm ci
COPY frontend/upcycleconnect/ .
ARG VITE_API_URL
ARG VITE_STRIPE_PUBLISHABLE_KEY
ARG VITE_ONESIGNAL_APP_ID
ENV VITE_API_URL=$VITE_API_URL
ENV VITE_STRIPE_PUBLISHABLE_KEY=$VITE_STRIPE_PUBLISHABLE_KEY
ENV VITE_ONESIGNAL_APP_ID=$VITE_ONESIGNAL_APP_ID
RUN npm run build

FROM nginx:alpine
COPY --from=frontend /app/dist /usr/share/nginx/html
COPY docker/nginx/nginx.conf /etc/nginx/conf.d/default.conf
