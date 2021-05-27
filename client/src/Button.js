import React from 'react';
import './Button.css';

export default class Button extends React.Component {
    handleClick = () => {
        window.location.href="/login";
    }
    render(){
        return(
            <div onClick={this.handleClick} className={this.props.css}>{this.props.text}</div>
        );
    };
}