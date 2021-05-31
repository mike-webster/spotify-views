import React from 'react';
import './ActionSelect.css';

export default class ActionSelect extends React.Component {
    constructor(props){
        super(props);
        this.state = {
            action: "recommendations"
        }

        this.changeAction = this.changeAction.bind(this);
    };

    changeAction = (e) => {
        this.setState({action: e.target.value});
    };

    render(){
        return <React.Fragment>
            <div className="action-select">
                <p>What would you like to see?</p>
                <select value={this.state.action} onChange={this.changeAction} id="action">
                    <option value="recommendations">Recommendations</option>
                    <option value="tracks">Top Tracks</option>
                    <option value="artists">Top Artists</option>
                    <option value="genres">Top Genres</option>
                </select>
            </div>
        </React.Fragment>
    };

}