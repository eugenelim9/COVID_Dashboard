import React, { Component } from 'react';
import PageTypes from '../../../../Constants/PageTypes/PageTypes';
import './Styles/MainPageContent.css';
import api from '../../../../Constants/APIEndpoints/APIEndpoints';
import MyD3Component from './dash.js';
import PropTypes from 'prop-types';

/**
 * @class
 * @classdesc SignUp handles the sign up component
 */
class MainPageContent extends Component {
    constructor(props) {
        super(props);

        this.state = {
            loading: 'initial',
            data: '',
            dashInfo: {},
            state: "",
            state1: "",
            filter: "",
            private: ""
        };
    }

    async loadData() {
        const response = await fetch(api.base + api.handlers.dashboards, {
            method: "GET",
            headers: new Headers({
                "Content-Type": "application/json",
                "Authorization": localStorage.getItem("Authorization")
            })
        });
        if (response.status >= 300) {
            const error = response.text();
            return;
        }
        var dashInfo = await response.json();
        return dashInfo[0];
    }

    componentDidMount() {
        this.setState({ loading: 'true' });
        this.loadData()
            .then((result) => {
                this.setState({
                    data: result,
                    loading: 'false',
                    dashInfo: result,
                    state: result.params.state,
                    state1: result.params.state1,
                    filter: result.params.filter,
                    private: result.private
                });
        });        
    }

    render() {
        if (this.state.loading === 'initial') {
        return <h2>Intializing...</h2>;
        }
        
        if (this.state.loading === 'true') {
        return <h2>Loading...</h2>;
        }

        return <> 
            <MyD3Component 
                dashInfo={{
                    dashboardID: this.state.dashInfo._id,
                    title: this.state.dashInfo.title,
                    description: this.state.dashInfo.description,
                    chartType  : "",
                    creator: this.state.dashInfo.creator,
                    createdAt: this.state.dashInfo.createdAt,
                    editedAt: this.state.dashInfo.editedAt
                }} 

                stateParam={this.state.state}
                state1Param={this.state.state1}
                filterParam={this.state.filter}
                privateParam={this.state.private}
                
                owner={true}
            />
        </>
    }
}

export default MainPageContent;