import React from 'react';
import './Home.css';
import Button from './Button.js';
import CookieBanner from './CookieBanner.js';

export default class Home extends React.Component {
    componentDidMount(){
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
                        <Button path="/login" text="LOG IN WITH SPOTIFY" css="btn half-width center-text" />
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