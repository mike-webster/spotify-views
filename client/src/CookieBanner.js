import React from 'react';
import './CookieBanner.css';

export default class CookieBanner extends React.Component {
    constructor(props){
        super(props);
        this.setCookie = this.setCookie.bind(this);
    };

    setCookie = () => {
        // we're setting a cookie so we know we don't need to
        // display the banner again for this session.
        let d = new Date();
        let exdays = 30;
        d.setTime(d.getTime() + (exdays*24*60*60*1000));

        let expires = "expires="+ d.toUTCString();

        document.cookie = "cookie-banner=1;" + expires + ";path=/";
        document.getElementById("cookie-back").style.display="none";
    };

    handleCloseClick = (e) => {
        e.stopPropagation();
        this.setCookie();
    };

    handleBannerClick = (e) => {
        e.stopPropagation();
    };

    handleBckgrndClick = () =>  {
        document.getElementById("cookie-back").style.display="none";
    };

    componentDidMount() {
        let cookies = document.cookie.split(';');
        let found = false;
        
        // iterate cookies to see if we have set one for the
        // cookie banner (meaning we don't need to show it again)
        for (var i = 0; i < cookies.length; i++) {
            if (cookies[i].includes("cookie-banner=")) {
                found = true;
                break;
            }
        }

        if (!found) {
            // if we didn't find the banner, it means it hasn't been
            // dismissed yet, so we should show the user the banner.
            document.getElementById("cookie-back").style.display="flex";
        }
    };

    render() {
        return(
            <div style={{display: "none"}} onClick={this.handleBckgrndClick} id="cookie-back" className="cookie-back">
                <div id="cookie-banner" className="cookie-banner">
                    <div onClick={this.handleBannerClick} className="banner-content">By using this website, you agree to our use of cookies. We use cookies to improve the user experience. </div>
                    <div onClick={this.handleCloseClick} className="cookie-close">X</div>
                </div>
            </div>
        );
    };
}