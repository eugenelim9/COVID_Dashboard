import React, { Component } from 'react';
import PropTypes from 'prop-types';
import SignForm from '../SignForm/SignForm';
import api from '../../../../Constants/APIEndpoints/APIEndpoints';
import Errors from '../../../Errors/Errors';
import PageTypes from '../../../../Constants/PageTypes/PageTypes';
import { Button, Layout, Menu, Breadcrumb, Modal, Input } from 'antd';

const { Header, Content, Footer } = Layout;

class ModalX extends React.Component {
  state = {
    modal2Visible: false,
  };

  setModal2Visible(modal2Visible) {
    this.setState({ modal2Visible });
  }

  render() {
    return (
      <>
        <Button onClick={() => this.setModal2Visible(true)}>
          Sign In
        </Button>
        <Modal
          title=" "
          centered
          visible={this.state.modal2Visible}
          onCancel={() => this.setModal2Visible(false)}
          footer={[
            <Button key="back" onClick={() => this.setModal2Visible(false)}>
              Cancel
            </Button>
          ]}
        >
        <form onSubmit={this.props.submitForm}>
            {this.props.fields.map(d => {
                const { key, name } = d;
                return <div key={key} className="home-signform">
                    <Input
                        placeholder={name}
                        value={this.props.values[key]}
                        name={key}
                        onChange={this.props.setField}
                        type={key === "password" || key === "passwordConf" ? "password" : ''}
                    />
                </div>
            })}
            <Input type="submit" value="Submit" />
        </form>
        </Modal>
      </>
    );
  }
}

/**
 * @class
 * @classdesc SignIn handles the sign in component
 */
class SignInNew extends Component {
    static propTypes = {
        setPage: PropTypes.func,
        setAuthToken: PropTypes.func
    }

    constructor(props) {
        super(props);

        this.state = {
            email: "",
            password: "",
            error: ""
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

    /**
     * @description submitForm handles the form submission
     */
    submitForm = async (e) => {
        e.preventDefault();
        const { email, password } = this.state;
        const sendData = { email, password };
        const response = await fetch(api.base + api.handlers.sessions, {
            method: "POST",
            body: JSON.stringify(sendData),
            headers: new Headers({
                "Content-Type": "application/json"
            })
        });
        if (response.status >= 300) {
            const error = await response.text();
            this.setError(error);
            return;
        }
        
        const authToken = response.headers.get("Authorization")
        localStorage.setItem("Authorization", authToken);
        this.setError("");
        this.props.setAuthToken(authToken);
        const user = await response.json();
        this.props.setUser(user);
    }

    render() {
        const values = this.state;
        const { error } = this.state;

        return <>
            <div className="home-div">
            <Header className="home-header" style={{ position: 'relative', zIndex: 1, width: '100%' }}>
                <div className="big-daddo">
                    <h1 className="home-title lit"> Dashy-19 </h1>
                    <div className="home-signin" style={{marginLeft: 200, marginRight: 0}}>
                        <ModalX submitForm={this.submitForm} setField={this.setField} values={values} fields={this.fields}>
                            <SignForm
                                setField={this.setField} 
                                submitForm={this.submitForm}
                                values={values}
                                fields={this.fields} />
                        </ModalX>
                    </div>
                    <div className="sign-forgot home-title" style={{marginLeft: 200, marginRight: 0}}>
                        <Button className="home-signup" onClick={(e) => this.props.setPage(e, PageTypes.signUp)}>Sign Up</Button>
                        <Button className="home-forgot" onClick={(e) => this.props.setPage(e, PageTypes.signIn)}>Home</Button>                    
                    </div>
                </div>
            </Header>
            <Errors className="error" error={error} setError={this.setError} />
            <Content className="site-layout content" style={{ padding: '0 50px'}}>
            </Content>
            <Footer className="le-foot" style={{ textAlign: 'center' }}> By Albert | Eugene | Gavin | Tow </Footer>
            </div>
        </>
    }
}

export default SignInNew;