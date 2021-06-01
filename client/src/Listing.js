import React from 'react';
import './Listing.css';
import './Shared.css';
import Footer from './Footer.js';
import Nav from './Nav.js';
import Recommendations from './Recommendations.js';
import ActionSelect from './ActionSelect.js';
export default class Listing extends React.Component {
    constructor(props){
        super(props);
    };

    render(){
        return <React.Fragment>
                <Nav />
                <div className="body">
                    <ActionSelect />
                </div>
                <Footer />
            </React.Fragment>
    };
};