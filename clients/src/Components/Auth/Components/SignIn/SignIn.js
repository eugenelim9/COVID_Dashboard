import React, { Component } from 'react';
import PropTypes from 'prop-types';
import SignForm from '../SignForm/SignForm';
import api from '../../../../Constants/APIEndpoints/APIEndpoints';
import Errors from '../../../Errors/Errors';
import PageTypes from '../../../../Constants/PageTypes/PageTypes';
import Card from './viz_cards';
import { Button, Layout, Menu, Breadcrumb, Modal, Input } from 'antd';

const { Header, Content, Footer } = Layout;

/**
 * @class
 * @classdesc SignIn handles the sign in component
 */
class SignIn extends Component {
    static propTypes = {
        setPage: PropTypes.func,
        setAuthToken: PropTypes.func
    }

    constructor(props) {
        super(props);

        this.state = {
            email: "",
            password: "",
            error: "",
            cards: []
        };

        this.fields = [
            {
                name: "Email",
                key: "email"
            },
            {
                name: "Password",
                key: "password"
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

    componentDidMount() {
        let currentComponent = this;
        var cardsInfo = [];
        var plz;
        async function loadCards() {
            const response = await fetch(api.base + api.handlers.dashboard, {
            method: "GET",
            headers: new Headers({
                "Content-Type": "application/json"
            })
            });
            if (response.status >= 300) {
                const error = response.text();
                return;
            }
            cardsInfo = await response.json();
            return cardsInfo;
        }
        loadCards().then(function(result) {
            currentComponent.setState({
                cards: result
            })
        })
    }

    render() {
        const values = this.state;
        const { error } = this.state;
        var cardsInfo = [];
        let plz;
        return <>
        <div className="home-div">
        <Header className="home-header" style={{ position: 'relative', zIndex: 1, width: '100%' }}>
            <div className="big-daddo">
                <h1 className="home-title lit"> Dashy-19 </h1>
                <div className="sign-forgot home-title">
                    <Button className="home-signin" onClick={(e) => this.props.setPage(e, PageTypes.signInNew)}>Sign In</Button>
                    <Button className="home-signup" onClick={(e) => this.props.setPage(e, PageTypes.signUp)}>Sign Up</Button>
                </div>
            </div>
        </Header>
        <Errors className="error" error={error} setError={this.setError} />
        <Content className="site-layout content" style={{ padding: '0 50px'}}>
            <div className="cards">
            {
                (this.state.cards).map((card) => {
                    return (
                        <Card
                            cardInfo={card}
                            setPage={this.props.setPage}
                        />
                    )
                })
            }
            </div>
        </Content>
        <Footer className="le-foot" style={{ textAlign: 'center', position: 'relative', bottom: 0 }}> By Albert | Eugene | Gavin | Tow </Footer>
        </div>
    </>
    }
}

export default SignIn;
