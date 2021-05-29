import React from 'react';
import './Button.css';

export default class Button extends React.Component {
    handleClick = () => {
        window.location.href=process.env.REACT_APP_API_BASE_URL + this.props.path;
    }
    render(){
        return(
            <div onClick={this.handleClick} className={this.props.css}>{this.props.text}</div>
        );
    };
}