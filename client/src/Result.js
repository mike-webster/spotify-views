import React from 'react';
import './Result.css';

export default class Result extends React.Component {
    render(){
        return(
            <a className="flex-table-item" href={this.props.url}>
                <div className="album" style={{backgroundImage: `url(${this.props.image})`}}>
                    <span className="blurry-back-white album-info">{this.props.artist}</span>
                    <span className="blurry-back-white album-info">{this.props.name}</span>
                </div>
            </a>
        );
    }
};