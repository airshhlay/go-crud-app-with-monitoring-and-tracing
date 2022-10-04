# Shopee Favourites
A favouriting system for the Shopee web page.

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

**Load Testing**
- k6: Executes scripts that call the various endpoints
- influxdb: Stores the load testing results from k6 and used as a datasource for Grafana

## Setup

### Prerequisites
* You have Docker Desktop installed and running.
* You have MySQL locally installed and running.

### Running the entire project

#### Database setup
1. Login to MySQL locally as a root user using `mysql -u root -p` and enter your password
2. Run the command `source <ABSOLUTE PATH TO ROOT FOLDER OF PROJECT>/services/userService/db/schema/mysql.sql` to create the `userservicedb` and the necessary tables.
3. Run the command `source <ABSOLUTE PATH TO ROOT FOLDER OF PROJECT>/services/itemService/db/schema/mysql.sql` to create the `itemservicedb` and the necessary tables.

### Running the services
1. Ensure you are at the root folder of the project.
2. Run the command `docker-compose build` to create and build the necessary images
3. Run `docker-compose up` to start the containers
4. Enter `localhost:80` or `localhost` in your browser. You should now be able to see the web application.

### Running monitoring
1. Open a separate terminal window in the root folder of the project and `cd` into the monitoring folder with `cd monitoring`
2. Run `docker-compose up` to start Prometheus, Grafana, Jaeger and InfluxDB
3. Run the load testing script with `docker-compose run k6 run /scripts/script.js --vus <number of virtual users> --duration <duration to run the script>`, or to run with a default of 10 VUs, just use the command ``docker-compose run k6 run /scripts/script.js` 
4. You can now view the various UIs below.

### Using the web app
1. In order to use Shopee Favourites, you will have to first create an account. Click on signup and enter a username and password.
2. If you receive a successful signup message, you can now login.
3. Login with your credentials.
4. Go to `shopee.sg` and copy/paste a link from any item you like into the input field at the top of the page.
5. If you provided a valid link, you will now see your item added to your list, along with its price! :)

## Viewing the UIs
Jaeger UI: `localhost:16886`\
Grafana UI: `localhost:3000`

The dashboards for Shopee Favourites and K6 will be automatically imported and inside the Grafana UI. Alternatively, you can import the dashboards directly by going to the `monitoring/dashboards` folder and copy-pasting the contents after clicking `Dashboards > Import` in the Grafana UI.\
Import the Node Exporter dashboard with ID: `1860`\
Import the MySQL Exporter dashboard with ID: `14057`
