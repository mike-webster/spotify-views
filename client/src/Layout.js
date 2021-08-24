import React from 'react';
import Nav from './Nav.js';
import Footer from './Footer.js';
import './Layout.css';

export default class Layout extends React.Component {
    render(){
        return(
            <React.Fragment>
                {this.props.nav === "true" ? <Nav /> : ""}
                {this.props.children}
                <Footer />
            </React.Fragment>
        );
    }
}