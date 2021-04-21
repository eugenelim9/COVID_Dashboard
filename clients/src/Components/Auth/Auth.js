import React from 'react';
import PropTypes from 'prop-types';
import PageTypes from '../../Constants/PageTypes/PageTypes';
import SignUp from './Components/SignUp/SignUp';
import SignIn from './Components/SignIn/SignIn';
import ForgotPassword from './Components/ForgotPassword/ForgotPassword';
import Clean from './Components/Clean';
import SignInNew from './Components/SignInNew/SignInNew';

/**
 * @class Auth
 * @description This is an auth object that controls what page
 * is loaded based on sign up or sign in state
 */
const Auth = ({ page, setPage, setAuthToken, setUser }) => {
    switch (page) {
        case PageTypes.signUp:
            return <SignUp setPage={setPage} setAuthToken={setAuthToken} setUser={setUser} />
        case PageTypes.signIn:
            return <SignIn setPage={setPage} setAuthToken={setAuthToken} setUser={setUser} />
        case PageTypes.forgotPassword:
            return <ForgotPassword setPage={setPage} />
        case PageTypes.clean:
            return <Clean setPage={setPage} />
        case PageTypes.signInNew:
            return <SignInNew setPage={setPage} setAuthToken={setAuthToken} setUser={setUser} />;
        default:
            return <>Error, invalid path reached</>
    }
}

Auth.propTypes = {
    page: PropTypes.string.isRequired,
    setPage: PropTypes.func.isRequired,
    setAuthToken: PropTypes.func.isRequired,
    setUser: PropTypes.func.isRequired
}

export default Auth;