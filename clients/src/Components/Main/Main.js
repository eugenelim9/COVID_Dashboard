import React from 'react';
import PageTypes from '../../Constants/PageTypes/PageTypes';
import MainPageContent from './Content/MainPageContent/MainPageContent';
import SignOutButton from './Components/SignOutButton/SignOutButton';
import UpdateName from './Components/UpdateName/UpdateName';
import UpdateAvatar from './Components/UpdateAvatar/UpdateAvatar';

const Main = ({ page, setPage, setAuthToken, setUser, user }) => {
    let content = <></>
    let contentPage = true;
    switch (page) {
        case PageTypes.signedInMain:
            content = <MainPageContent user={user} setPage={setPage} />;
            break;
        default:
            content = <>Error, invalid path reached</>;
            contentPage = false;
            break;
    }
    return <>
        <SignOutButton setUser={setUser} setAuthToken={setAuthToken} />
        {content}
    </>
}

export default Main;