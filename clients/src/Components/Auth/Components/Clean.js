import React, { Component } from 'react';
import PropTypes from 'prop-types';
import Dash from '../../Main/Content/MainPageContent/dash';
import PageTypes from '../../../Constants/PageTypes/PageTypes';
import { Button } from 'antd';

/**
 * @class
 * @classdesc SignUp handles the sign up component
 */
class Clean extends Component {
    constructor(props) {
        super(props);
    }

    render() {
        const dash = JSON.parse(localStorage.getItem('dash'))
        return <>
            <Button className="goHome" onClick={(e) => this.props.setPage(e, PageTypes.signIn)}>
                Home
            </Button>
            <h1 className="others-dash">Dashy-19</h1>
            <Dash 
                dashInfo={{
                    dashboardID: dash._id,
                    title: dash.title,
                    description: dash.description,
                    chartType  : "",
                    creator: dash.creator,
                    createdAt: dash.createdAt,
                    editedAt: dash.editedAt
                }} 

                stateParam={dash.params.state}
                state1Param={dash.params.state1}
                filterParam={dash.params.filter}
                privateParam={dash.private}
                owner={false} 
            />
        </>
    }
}

export default Clean;