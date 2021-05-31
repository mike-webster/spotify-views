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
        this.state = {
            items: [],
        };

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
                
                alert("There was an error with the request, your session has probably expired.");
                document.cookie = "sv-authed=";
                window.location.href = "/";
            }
        )
    };

    render(){
        return <React.Fragment>
                <Nav />
                <div className="body">
                    <ActionSelect />
                    <Recommendations />
                </div>
                <Footer />
            </React.Fragment>
    };
};