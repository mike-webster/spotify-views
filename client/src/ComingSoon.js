import React from 'react';
import './ComingSoon.css';

export default class ComingSoon extends React.Component {
    render(){
        return(
            <div className="body">
                <p className="copy">Sorry! We're having some issues with a third party tool that was being used.
                    This feature will be back as soon as it gets handled!
                </p>
            </div>
        );
    }
}