import React, { Component } from 'react';
import PropTypes from 'prop-types';
import SignForm from '../SignForm/SignForm';
import api from '../../../../Constants/APIEndpoints/APIEndpoints';
import Errors from '../../../Errors/Errors';
import PageTypes from '../../../../Constants/PageTypes/PageTypes';
import { Button, Layout, Menu, Breadcrumb, Modal, Input } from 'antd';

const { Header, Content, Footer } = Layout;

/**
 * @class
 * @classdesc SignUp handles the sign up component
 */
class SignUp extends Component {
    static propTypes = {
        setPage: PropTypes.func,
        setAuthToken: PropTypes.func
    }

    constructor(props) {
        super(props);

        this.state = {
            email: "",
            userName: "",
            firstName: "",
            lastName: "",
            password: "",
            passwordConf: "",
            error: ""
        };

        this.fields = [
            {
                name: "Email",
                key: "email"
            },
            {
                name: "Username",
                key: "userName"
            },
            {
                name: "First name",
                key: "firstName"
            },
            {
                name: "Last name",
                key: "lastName"
            },
            {
                name: "Password",
                key: "password"
            },
            {
                name: "Password Confirmation",
                key: "passwordConf"
            }];
    }

    /**
     * @description setField will set the field for the provided argument
     */
    setField = (e) => {
        this.setState({ [e.target.name]: e.target.value });
    }

    /**
     * @description setError sets the error message
     */
    setError = (error) => {
        this.setState({ error })
    }

    /**
     * @description submitForm handles the form submission
     */
    submitForm = async (e) => {
        e.preventDefault();
        const { email,
            userName,
            firstName,
            lastName,
            password,
            passwordConf } = this.state;
        const sendData = {
            email,
            userName,
            firstName,
            lastName,
            password,
            passwordConf
        };

        let token;

        const f1 = await fetch(api.base + api.handlers.users, {
            method: "POST",
            body: JSON.stringify(sendData),
            headers: new Headers({
                "Content-Type": "application/json"
            })
        }).then((response) => {
            const authToken = response.headers.get("Authorization")
            localStorage.setItem("Authorization", authToken);
            this.setError("");
            this.props.setAuthToken(authToken);
            let user;
            const answer = response.json().then((data) => {
                user = data;
                this.props.setUser(user)
            });
            token = authToken;
        })

        var postBody = {
            title:"test3 dash",
            description:"testing",
            params: {
                state: "Washington",
                state1: "Ohio",
                filter: "Tested"
            },
            private:"false"
        }

        const f2 = await fetch(api.base + api.handlers.dashboard, {
            method: "POST",
            body: JSON.stringify(postBody),
            headers: new Headers({
                "Content-Type": "application/json",
                "Authorization": token
            })
        }).then((response) => {
            const authToken = localStorage.getItem("Authorization");
            this.setError("");
            this.props.setAuthToken(authToken);
            response.json()
        }) 

        async function fetchURLs() {
            try {
                var data = await Promise.all([
                    f1,
                    f2
                ]);
            } catch (error) {
                console.log(error);
            }
        }
    }

    render() {
        const values = this.state;
        const { error } = this.state;
        return <>
            <div className="home-div">
            <Header className="home-header" style={{ position: 'relative', zIndex: 1, width: '100%'}}>
                <div className="big-daddo">
                    <h1 className="home-title lit"> Dashy-19 </h1>
                    <div className="sign-forgot">
                        <Button className="home-signup" onClick={(e) => this.props.setPage(e, PageTypes.signInNew)}>Sign In Instead</Button>
                    </div>
                </div>
            </Header>
            <Errors error={error} setError={this.setError} />
            <Content className="site-layout content signupfoot" style={{ padding: '0 50px'}}>
                <SignForm
                    setField={this.setField}
                    submitForm={this.submitForm}
                    values={values}
                    fields={this.fields} />
            </Content>
            <Footer className="le-foot" style={{ textAlign: 'center' }}> By Albert | Eugene | Gavin | Tow </Footer>
            </div>
        </>
    }
}

export default SignUp;