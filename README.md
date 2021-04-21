# Project: Dashy-19

## Project Description:
Dashy-19 is for people who care about COVID-19 trends and are interested in creating their own dashboards, looking at metrics that matter to other people, and to promote awareness. The platform enables users to find current and customizable information by viewing other people’s dashboards as well. They can use this tool to easily understand the virus better, thus properly preparing themselves for the virus. As developers, we want to make an app that spreads awareness for the COVID-19 pandemic.  With all the misinformation spreading about COVID-19, we want to demystify the pandemic by providing a network of readily available visualizations. 

The platform will have a home page where everybody lands. Without sign in, people can view dashboards published from other people. Dashboards are visualizations of the current COVID-19 data, and there will be many things for creators to play with to really tailor the data to their needs. A logged in/authenticated user can also create, save, and publish their very own dashboard.

## Infastructure
Our server relies on a Microservice infastructure where we have a seperate docker container that all dashboard related requests will be forwarded too. The client side exclusively interacts with our whole domain container which funnels all requests through the gateway and forwards it appropraitely to our microservices and databases.

![Image of Yaktocat](infastructure.jpg)

## Use Cases and Priority

| Priority | User               | Description                                                                           |
|----------|--------------------|---------------------------------------------------------------------------------------|
| P0       | Customized User    | I want to be able to create an account and save my personal dashboard                 |
| P0       | Power User         | I want to be able to view other users’ dashboards                                     |
| P1       | Infrequent Checker | I want to be able to view and edit my dashboard for COVID related data in my location |
| P1       | A Sharer           | I want to be able to make my dashboard public so that other users can view it         |

## Endpoints

**User auth-signin**

`/v1/users/`
- POST 
  - 201: User created
  - 401: Wrong credentials 

`/v1/sessions`
- POST 
  - 201: created a new user session
  - 403: invalid username/email forbidden
  - 415: unsupported media
  - 500: internal server error
- DELETE:
  - 200: Successfully delete session
  - 400: Bad request

**Dashboard**

`/v1/dashboards/`
- GET - Get all dashboards that are not private (all the public ones)
  - 200: Successfully get public dashboards
  - 400: Bad request

`/v1/dashboard/:dashID`
- GET - Return the dashboard of the passed in dashboard id. 
  - 200: Successfully get specific dashboards
  - 400: Bad request
  - 401: Unauthorized - not public / or not the creator: If the user does not own the dashboard and if it’s not public (private = false) 
- PATCH
  - 200: Successfully update data
  - 400: Bad request

**Data**

`/v1/data`
- GET - Get the data from the database and send it in the request for the frontend to use it for visualization
  - 200: Successfully get data


## Models
We will be utilizing both SQL and MongoDB to store relevant information to the database.
We are using MySQL to store user information and MongoDB to store the dashboard and data (for visualization) information.

#### Users (SQL)
The following is the user model that is created when signing up.
```
{
  id: "user_id_value"
  email: "email",
  pass_hash: "password_hash".
  username: "username",
  first_name: "first_name",
  last_name: "last_name",
  photo_url: "urlPhoto"
}
```

#### Dashboard
The dashboard model when a new dashboard is created.
```
{
    creator: "id_of_user_who_created",
    title: "title_of_dashboard",
    description: "description",
    params: "object_of_parameters",
    createdAt: "date_created",
    editedAt: "date_edited",
    private: "boolean_private"
}
```

#### Data Model (mongoDB)
The following is the model to the data that we store and update to use for visualization of the dashboard.
```
{
    state: "state",
    tested: "number_of_tested_cases",
    infected: "number_of_infected_cases",
    deaths: "number_of_deaths",
    population: "number_of_population",
    popDensity: "population_density",
    gini: "GINI_index",
    icuBeds: "number_of_icu_beds",
    income: "total_income",
    gdp: "GDP",
    unemployment: "unemployment",
    sexRatio: "sex_ratio",
    smokingRate: "smoking_rate",
    fluDeaths: "number_of_death_from_flu",
    respDeaths: "bumber_of_respiratory_deaths",
    physicians: "number_of_physicians",
    hospitals: "number_of_hospitals",
    healthSpending: "total_health_spending",
    pollution: "number_of_population",
    medLargeAirports: "medium_to_large_airports",
    temperature: "temperature",
    urban: "urban_percentage",
    schoolClosureDate: "school_closure_date"
}
```
