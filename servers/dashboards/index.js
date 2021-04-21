const mongoose = require('mongoose');
const express = require('express');
const axios = require('axios').default;
// global.fetch = require("node-fetch");

const csv = require('csv-parser');
const fs = require('fs');

const { dashboardSchema, countriesCovidSchema } = require('./schema');
const { 
    allDashHandler, 
    specificDashGetHandler,
    allDashPostHandler,
    specificDashPatchHandler,
    dataGetHandler,
    dataDelHandler
} = require('./handlers')
const mongoEndpoint = process.env.MONGO_ENDPOINT
// TODO: port and endpoint in mongo
const port = process.env.DASHBOARDPORT;

const CountriesCovid = mongoose.model("CountriesCovid", countriesCovidSchema)
const Dashboard = mongoose.model("Dashboard", dashboardSchema)


const app = express();
app.use(express.json());

const connect = () => {
    mongoose.connect(mongoEndpoint, { useFindAndModify: false });
};

// connecting mongo db
connect();
// open main on connection
mongoose.connection.on('error', console.error)
    .on('disconnected', connect)
    .once('open', main);


// wrapper for handlers, ensures X-User is valid
const RequestWrapper = (handler, SchemeAndDBForwarder) => {
    return (req, res) => {
        let user;

        try {
            user = JSON.parse(req.get("X-User"))
        } catch (err) {
            res.status(401).send("Unauthorized")
            return
        }   
        handler(req, res, SchemeAndDBForwarder, user);
    }
};

const SimpleWrapper = (handler, SchemeAndDBForwarder) => {
    return (req, res) => {
        handler(req, res, SchemeAndDBForwarder)
    }
}

const methodNotAllowedHandler = (req, res) => {
    res.status(405).send("method not allowed")
}

// /v1/dashboards
app.route("/v1/dashboards")
    // get - get all dashboards that are not private and respond to user with dashboard info as JSON array
    .get(SimpleWrapper(allDashHandler, { Dashboard }))
    // post - create new chart with JSON body
    // JSON body the following info - title,creator, description, chartParams
    // needs title and params if not, status bad request
    // respond with the chart as JSON and 201
    .post(RequestWrapper(allDashPostHandler, { Dashboard }))
    .all(methodNotAllowedHandler)


// /v1/dashboards/{dashboardsID}
app.route("/v1/dashboards/:dashID")
    // get - respond with info if user is creator of the dashboard or if not private
    // else respond 403(forbidden)
    .get(RequestWrapper(specificDashGetHandler, { Dashboard }))
    // patch - update dashboard with JSON body that includes the information to be updated.
    // user needs to be creator of the dashboard else 403
    // if dashbaordsIDs don't exist, respond 400
    // respond with updated dashboard as JSON
    .patch(RequestWrapper(specificDashPatchHandler, { Dashboard }))
    .all(methodNotAllowedHandler)

// Get the data from the database and send it in the request for the frontend to use it for visualization
app.route("/v1/data")
    .get(SimpleWrapper(dataGetHandler, { CountriesCovid }))
    .delete(SimpleWrapper(dataDelHandler, { CountriesCovid }))
    .all(methodNotAllowedHandler)



async function populateDB() {
    try {
        // if db is empty, if not, populate
        const exists  = await CountriesCovid.exists({state: "Alabama"})
        if (exists === false) {
            const results = []
            fs.createReadStream('COVID19_state.csv')
                .pipe(csv({}))
                .on('data', (data) => {
                    results.push(data)
                })
                .on('end', async () => { 
                    console.log("parsing done")
                    // console.log(results)
                    for (let i = 0; i < results.length; i++) {
                        let data = results[i]
                        let newEntry = {
                            state: data.State,
                            tested: data.Tested,
                            infected: data.Infected,
                            deaths: data.Deaths,
                            population: data.Population,
                            popDensity: data["Pop Density"],
                            gini: data.Gini,
                            icuBeds: data["ICU Beds"],
                            income: data.Income,
                            gdp: data["GDP"],
                            unemployment: data.Unemployment,
                            sexRatio: data["Sex Ratio"],
                            smokingRate: data["Smoking Rate"],
                            fluDeaths: data["Flu Deaths"],
                            respDeaths: data["Respiratory Deaths"],
                            physicians: data.Physicians,
                            hospitals: data.Hospitals,
                            healthSpending: data["Health Spending"],
                            pollution: data.Pollution,
                            medLargeAirports: data["Med-Large Airports"],
                            temperature: data.Temperature,
                            urban: data.Urban,
                            schoolClosureDate: data["School Closure Date"]
                        }
                        const entryData = new CountriesCovid(newEntry);
                        await entryData.save()
                    }
                    
                    console.log("end data alsdflsdg")
                })
        }
    } catch(e) {
        console.log(e)
    }
}

async function main() {
    await populateDB()
    app.listen(port, "", () => {
        console.log(`server listening at ${port}`);
    })
    
}
