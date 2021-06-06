import React from 'react';
import './ComingSoon.css';

export default class ComingSoon extends React.Component {
    render(){
        let content = "This feature is being worked on now. Coming soon!";
        if (this.props.err==="true") {
            content = "Sorry! We're having some issues with a third party tool that was being used. This feature will be back as soon as it gets handled!"
        } 
        return(
            <div className="body">
                <p className="copy">{content}</p>
            </div>
        );
    }
}