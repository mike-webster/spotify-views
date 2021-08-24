import React from 'react';
import './Home.css';
import './Shared.css';
import Button from './Button.js';
import CookieBanner from './CookieBanner.js';
import Layout from './Layout.js';
export default class Home extends React.Component {
    constructor(){
        super();
    };

    componentDidMount(){
        // check for ?authed
        if (window.location.search.includes("?svauth")) {
            try{
                let val = window.location.search.split("?svauth=")[1];
                let d = new Date();
                let numHours = 2;
                d.setTime(d.getTime() + (numHours*60*60*1000));

                let expires = "expires="+ d.toUTCString();
                document.cookie = "svauth=" + val + ";" + expires + ";path=/";
                window.location.href = "/discover";
            }catch (error) {console.log(error);}
        }
    };

    render(){
        return(<p>error</p>);
    };
}