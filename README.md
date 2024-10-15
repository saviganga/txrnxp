# txrnxp

## Setup

### Prerequisites

To run the development environment, you need to have the following installed:

- **GO**
- **Docker**
- **Docker Compose**
- **Environment Variables**

### Environment Variables

A `.env.example` file has been added to the repository. You should create a `.env` file and fill in the required fields with your values to configure your environment. The `VERSION` field has been prefilled to suit the configurations on the application.


## Run the project

1. Clone the repository
```bash
git clone git@github.com:saviganga/txrnxp.git
```

2. Set up your `.env` file
```bash
cd txrnxp/
cp .env.example .env
```
Fill the .env file with your values

3. Navigate to the folder ./initialisers/db.go and uncomment the migrations code to run the database migrations

4. Build the project using `docker-compose`
```bash
docker-compose build
```

5. After the build is completed, start the project with `docker-compose`
```bash
docker-compose up
```
