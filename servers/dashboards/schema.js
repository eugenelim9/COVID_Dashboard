const { SchemaType } = require('mongoose');

const Schema = require('mongoose').Schema;

const dashboardSchema = new Schema({
    creator: { type: Object, required: true, unique: true },
    title: {type: String, required: true},
    description: {type: String },
    params: {type: Object, required: true},
    createdAt: {type: Date, required: true},
    editedAt: {type: Date},
    private: {type: Boolean, default: true}
})

// const chartsSchema = new Schema({
//     dashboardID: {type: Schema.Types.ObjectId, required: true},
//     title: {type: String, required: true},
//     chartType: {type: String, required: true},
//     params: {type: Schema.Types.Mixed, required:false},
//     creator: {type: Object, required: true},
//     createdAt: {type: Date, required: true},
//     editedAt: {type: Date},
// })

const countriesCovidSchema = new Schema({
    state: {type: String, required: true},
    tested: {type: Number, required: true},
    infected: {type: Number},
    deaths: {type: Number},
    population: {type: Number},
    popDensity: {type: Number},
    gini: {type: Number},
    icuBeds: {type: Number},
    income:{type: Number},
    gdp:{type: Number},
    unemployment:{type: Number},
    sexRatio:{type: Number},
    smokingRate:{type: Number},
    fluDeaths:{type: Number},
    respDeaths: {type: Number},
    physicians: {type: Number},
    hospitals: {type: Number},
    healthSpending: {type: Number},
    pollution: {type:Number},
    medLargeAirports: {type: Number},
    temperature: {type: Number},
    urban: {type: Number},
    schoolClosureDate: {type: Date}
})

/*
  0: "State"
  1: "Tested"
  2: "Infected"
  3: "Deaths"
  4: "Population"
  5: "Pop Density"
  6: "Gini"
  7: "ICU Beds"
  8: "Income"
  9: "GDP"
  10: "Unemployment"
  11: "Sex Ratio"
  12: "Smoking Rate"
  13: "Flu Deaths"
  14: "Respiratory Deaths"
  15: "Physicians"
  16: "Hospitals"
  17: "Health Spending"
  18: "Pollution"
  19: "Med-Large Airports"
  20: "Temperature"
  21: "Urban"
  22: "Age 0-25"
  23: "Age 26-54"
  24: "Age 55+"
  25: "School Closure Date"
*/

module.exports = { dashboardSchema, countriesCovidSchema }