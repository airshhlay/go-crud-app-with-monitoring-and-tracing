# Shopee Favourites
## Introduction
An entry task for MTS Service Governance.

**Timeline**\
Design: 01 Aug 2022 - 02 Aug 2022\
Development: 03 Aug 2022 - 15 Aug 2022

## Project Components
**Frontend**
- Shopee Favourites: A simple web application made with React
- Nginx: A reverse proxy

**Backend**
- Gateway: A HTTP server to the frontend and GRPC client to the different microservices. Also performs authorization and authentication of users.
- User Service: A microservice that handles user login/signup. Makes use of its own MySQL database `userservicedb`
- Item Service: A microservice that handles user authentication 

**Observability**
- Prometheus: Scrapes metrics exposed by each service at the `/metrics` endpoint
- Grafana: Pulls data from Prometheus for visualization
- Jaeger: Collects and displays traces created by the services

## Setup

### Prerequisites
* You have Docker Desktop installed and running.
* You have MySQL locally installed and running.

### Running the entire project

#### Database setup
1. Login to MySQL locally as a root user using `mysql -u root -p` and enter your password
2. Run the command `source <ABSOLUTE PATH TO ROOT FOLDER OF PROJECT>/services/userService/db/schema/mysql.sql` to create the `userservicedb` and the necessary tables.
3. Run the command `source <ABSOLUTE PATH TO ROOT FOLDER OF PROJECT>/services/itemService/db/schema/mysql.sql` to create the `itemservicedb` and the necessary tables.

### Starting Docker
1. Ensure you are at the root folder of the project.
2. Run the command `docker-compose build` to create and build the necessary images
3. Run `docker-compose up` to start the containers
4. Enter `localhost:80` or `localhost` in your browser. You should now be able to see the web application.

### Using the web app
1. In order to use Shopee Favourites, you will have to first create an account. Click on signup and enter a username and password.
2. If you receive a successful signup message, you can now login.
3. Login with your credentials.
4. Go to `shopee.sg` and copy/paste a link from any item you like into the input field at the top of the page.
5. If you provided a valid link, you will now see your item added to your list, along with its price! :)

## Viewing the UI
Jaeger UI: `localhost:16886`\
Grafana UI: `localhost:3000`

To import the Grafana dashboard, click on `Dashboards > Import` in Grafana and copy and paste the contents of `grafana-shopee-favourites.json` 