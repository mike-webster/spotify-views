import React from 'react';
import './ActionSelect.css';
import Layout from './Layout.js';
import Recommendations from './Recommendations.jsx';
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
            content = <Recommendations key="recs" />;
        } else if (this.state.action === "tracks") {
            content = <Tops key="tracks" focus="tracks"/>;
        } else if (this.state.action === "artists") {
            content = <Tops key="artists" focus="artists"/>;
        }else if (this.state.action === "genres") {
            content = <Tops key="genres" focus="genres"/>;
        }

        return<Layout key="as-table" nav="true">
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
        </Layout>
    };

}