import React from 'react';
import './ActionSelect.css';
import Nav from './Nav.js';
import Footer from './Footer.js';
import Recommendations from './Recommendations.js';
import Tops from './Tops.js';

export default class ActionSelect extends React.Component {
    constructor(props){
        super(props);
        this.state = {
            action: "recommendations",
            body: "88vh"
        }

        this.changeAction = this.changeAction.bind(this);
    };

    changeAction = (e) => {
        this.setState({action: e.target.value});
    };

    render(){
        let content;
        if (this.state.action === "recommendations") {
            content = <Recommendations />;
        } else if (this.state.action === "tracks") {
            content = <Tops focus="tracks"/>;
        } else if (this.state.action === "artists") {
            content = <Tops focus="artists"/>;
        }else if (this.state.action === "genres") {
            content = <Tops focus="genres"/>;
        }

        return <React.Fragment>
            <Nav />
            <div className="body">
                <div className="action-select">
                    <p>What would you like to see?</p>
                    <select value={this.state.action} onChange={this.changeAction} id="action">
                        <option value="recommendations">Recommendations</option>
                        <option value="tracks">Top Tracks</option>
                        <option value="artists">Top Artists</option>
                        <option value="genres">Top Genres</option>
                    </select>
                </div>

                {content}
            </div>
            <Footer />
        </React.Fragment>
    };

}