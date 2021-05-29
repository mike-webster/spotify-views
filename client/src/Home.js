import React from 'react';
import './Home.css';
import './Shared.css';
import Button from './Button.js';
import CookieBanner from './CookieBanner.js';
import { Redirect } from "react-router-dom";

export default class Home extends React.Component {
    constructor(){
        super();
        this.checkCookie = this.checkCookie.bind(this);
        this.setCookie = this.setCookie.bind(this);
    };
    setCookie = () => {
        // we're setting a cookie to know we're logged in because
        // we can't see the cookie that was set by the api
        let d = new Date();
        let numHours = 2;
        d.setTime(d.getTime() + (numHours*60*60*1000));

        let expires = "expires="+ d.toUTCString();

        document.cookie = "sv-authed=1;" + expires + ";path=/";
    };
    checkCookie = () => {
        let cookies = document.cookie.split(';');
        for (var i = 0; i < cookies.length; i++) {
            if (cookies[i].includes("sv-authed=")) {
                console.log("found authed: ", cookies[i]);
                // token will still work, no need to auth
                window.location = "/discover";
            }
        }
    };
    componentDidMount(){
        // check for ?authed
        if (window.location.search.includes("?authed")) {
            this.setCookie();
        }

        // if we authed recently, no need to login
        this.checkCookie();
    };

    render(){
        return (
            <div className="container">
                <div className="home-banner">
                    <div className="blurry-back-white">
                        <h1 className="full-width center-text large-pad-top">Spotify Views</h1>
                        <p className="full-width center-text">Find More Of What You Love</p>
                    </div>
                    <div className="fix"></div>
                    <div className="blurry-back-white large-pad-v large-marg-top">
                        <Button path="/login?redirectUrl=?authed" text="LOG IN WITH SPOTIFY" css="btn half-width center-text" />
                    </div>
                </div>
                <p className="white-back small-marg med-pad-v large-pad-h">
                    Using Spotify Views, you can take a dive into your music taste to discover more about the music you love! Log in with your existing
                    Spotify account and learn about your top tracks, artists and genres, get recommendations for new artists, and more!
                </p>
                <CookieBanner />
            </div>
        );
    };
}