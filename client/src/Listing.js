import React from 'react';
import './Listing.css';
import Footer from './Footer.js';
import Nav from './Nav.js';
import Recommendations from './Recommendations.js'

export default class Listing extends React.Component {
    constructor(props){
        super(props);
        this.state = {
            items: [],
            action: "recommendations"
        };

        this.changeAction = this.changeAction.bind(this);
    };

    changeAction = (e) => {
        this.setState({action: e.target.value});
    };

    componentDidMount(){
        fetch(process.env.REACT_APP_API_BASE_URL + "/tracks/recommendations", {
            credentials: 'include'
        })
        .then(res => res.json())
        .then(
            (result) => {
                console.log(result.tracks)
                let tmp = [];
                for (var i = 0; i < result.tracks.length; i++) {
                    tmp.push(result.tracks[i])
                }
                this.setState({
                    state: "success",
                    items: tmp
                });
            },
            (error) => {
                // TODO: something in this error state
                this.setState({
                    state: "error",
                    error
                });
                console.log(error);
                console.log("redirecting");
                window.location.href = "/";
            }
        )
    };

    render(){
        return <React.Fragment>
                <Nav />
                <div className="body">
                    <p>What would you like to see?</p>
                    <select value={this.state.action} onChange={this.changeAction} id="action">
                        <option value="recommendations">Recommendations</option>
                        <option value="tracks">Top Tracks</option>
                        <option value="artists">Top Artists</option>
                        <option value="genres">Top Genres</option>
                    </select>
                    <Recommendations />
                </div>
                <Footer />
            </React.Fragment>
    };
};