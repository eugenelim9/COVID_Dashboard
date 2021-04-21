import React, { useRef, useEffect, useState } from 'react';
import * as d3 from 'd3';
import { Input, Switch, Menu, Dropdown } from 'antd';
import * as topojson from "topojson-client";
import api from '../../../../Constants/APIEndpoints/APIEndpoints';

/* Component */
const MyD3Component = ({ dashInfo, stateParam, state1Param, filterParam, privateParam, owner }) => {
    /* The useRef Hook creates a variable that "holds on" to a value across rendering
    passes. In this case it will hold our component's SVG DOM element. It's
    initialized null and React will assign it later (see the return statement) */
    const [state, setState] = useState(stateParam);
    const [state1, setState1] = useState(state1Param);
    const [covid, setCovid] = useState([]);
    const [filter, setFilter] = useState(filterParam)
    const [priv, setPriv] = useState(privateParam)

    const Switchy = ({}) => {
        function onChange(checked) {
            setPriv(!priv);
        }
        return(
            !priv ? <Switch defaultChecked onChange={onChange} /> : <Switch onChange={onChange} />
        );
    }

    const Dropfilter = ({}) => {
        var menu = (
            <Menu onClick={onClick}>
                <Menu.Item>
                    loading...
                </Menu.Item>
            </Menu>
        );

        const onClick = ({ key }) => {
            setFilter(key);
        }

        if (covid.length > 0) {
            menu = (
                <Menu onClick={onClick}>
                    {covid.columns.map((item) => {
                        return (
                            <Menu.Item key={item}>
                                {item}
                            </Menu.Item>
                        )
                    })}
                </Menu>
            )
        }

        return(
            <Dropdown overlay={menu}>
                <a style={{display:'block'}} className="ant-dropdown-link">
                    Click To Select Filter
                </a>
            </Dropdown>
        );
    }
    
    async function getData() {
        d3.csv("/renamed/COVID19_state.csv")
            .then((data) => {
                setCovid(data);
                makeMap(data);
                graph(data);
            })
            .catch((error) => {
                console.log(error);
            });
    }

    async function makeMap(data) {
        const height = 700
        const width = 1000

        let infected = new Map(data.map(d => [d.State, +d[filter]]));

        let shapefile = await d3.json("https://cdn.jsdelivr.net/npm/us-atlas@3/states-albers-10m.json");
        let path = d3.geoPath();
        let max_value = d3.max(data, d => Math.abs(d[filter]));
        let min_value = d3.min(data, d => Math.abs(d[filter]));
        let color = d3.scaleSequential([min_value, max_value], d3.interpolateReds);

        // Create the svg
        const svg = d3.select("#map")
            .attr("width", width)
            .attr("height", height)
            .style("text-anchor", "middle")
            .style('font-family', 'sans-serif')
            .style("font-size", "12px")

        // Add paths for each state
        const paths = svg.selectAll("path")
            .data(topojson.feature(shapefile, shapefile.objects.states).features)
            .join("path")
            .attr("d", path)
            .attr("fill", d => {
                return color(infected.get(d.properties.name))
            })

        // Add text for each state
        const labels = svg.selectAll("text")
            .data(topojson.feature(shapefile, shapefile.objects.states).features)
            .join("text")
            .text(d => infected.get(d.properties.name).toFixed(1))
                // .style("fill", "red")
            .attr("y", d => path.centroid(d)[1])
            .attr("x", d => path.centroid(d)[0])

        // Add a title
        const title = svg.append("text")
            .attr('x', width/2)
            .attr("y", 20)
            .text("Total COVID-19 " + filter + " Cases By State")
            .style("font-size", "20px")
    }

    function graph(data) {
        let infected = new Map(data.map(d => [d.State, +d.Infected]))

        for (let k of infected.keys()) {
            if (k !== state && k !== state1) {
                infected.delete(k);
            }
        }

        const margin = ({top: 30, right: 0, bottom: 80, left: 80})
        const height = 400
        const width = 800
        const color = "blue"

        const svg = d3.select("#chart")
            .attr("viewBox", [0, 0, width, height]);
        
        const y = d3.scaleLinear()
            .domain([0, d3.max(data, d => d.Infected)]).nice()
            .range([height - margin.bottom, margin.top])
        
        const x = d3.scaleBand()
            .domain(d3.range([...infected.keys()].length))
            .range([margin.left, width - margin.right])
            .padding(0.1)
        
        const xAxis = g => g
            .attr("transform", `translate(0,${height - margin.bottom})`)
            .call(d3.axisBottom(x).tickFormat(i => [...infected.keys()][i]).tickSizeOuter(0))

        const yAxis = g => g
            .attr("transform", `translate(${margin.left},0)`)
            .call(d3.axisLeft(y).ticks(null, data.format))
            .call(g => g.select(".domain").remove())
            .call(g => g.append("text")
                .attr("x", -margin.left)
                .attr("y", 10)
                .attr("fill", "currentColor")
                .attr("text-anchor", "start")
                .text(data.Infected))
        
        svg.selectAll("g").remove();
        svg.append("g")
            .attr("fill", color)
            .selectAll("rect")
            .data(infected)
            .join("rect")
            .attr("x", (d, i) => x(i))
            .attr("y", d => y(d[1]))
            .attr("height", d => y(0) - y(d[1]))
            .attr("width", x.bandwidth());
        
        svg.append("g")
            .attr("class", "x axis")
            .attr("transform", "translate(0," + height + ")")
            .call(xAxis)
            .selectAll("text")  
            .style("text-anchor", "end")
            .attr("dx", "-.8em")
            .attr("dy", ".15em")
            .attr("transform", "rotate(-65)" );
        
        svg.append("g")
            .call(yAxis);

        const title = svg.append("text")
            .attr('x', width/3)
            .attr("y", 20)
            .text("Comparing COVID Infected Cases")
            .style("font-size", "20px")
    }

    useEffect(() => {
        getData()
    }, [state, state1, filter]);

    const authToken = localStorage.getItem("Authorization");

    // use PATCH /v1/dashboards/{dashboardsID}
    let submitForm = async (e) => {
        e.preventDefault();
        let patchJson = {
            params: {
                state: state, 
                state1: state1,
                filter: filter
            },
            private: priv
        }
        const response = await fetch(api.base + api.handlers.dashboards, {
            method: "PATCH",
            body: JSON.stringify(patchJson),
            headers: new Headers({
                "Content-Type": "application/json",
                "Authorization": authToken
            })
        });
        if (response.status >= 300) {
            const error = await response.text();
            return;
        }
    }

    return (
        <div className="switchy">
            <div className="master">
                <div className="item1">
                {owner && <div className="item1-title">Compare States</div>}
                {owner && <Input className="state-filter"
                    type="text"
                    value={state}
                    onChange={event => setState(event.target.value)}
                />}
                {owner && <Input className="state-filter"
                    type="text"
                    value={state1}
                    onChange={event => setState1(event.target.value)}
                />}
                {owner && <div className="switch">
                    Public: &nbsp;
                    <Switchy />
                </div>}                 
                </div>                 
                <div className="item1">
                    {owner && <div className="item2-title">Filter Map</div>}
                    {owner && <Input className="state-filter"
                        type="text"
                        value={filter}
                        onChange={event => setFilter(event.target.value)}
                    />}
                    {owner && <Dropfilter />}
                </div>
                <div className="dash item2">  
                    <div className="userTitle">
                        {dashInfo.creator.firstName}'s Dash!!
                    </div>        
                    <svg
                        className="d3-component"
                        id="chart"
                    />   
                    <svg
                        className="d3-component"
                        id="map"
                    /> 
                </div>
                {owner && <button onClick={submitForm}><h3 className="save-button">Save Filters</h3></button>}
            </div>
        </div>
    );
}

export default MyD3Component;