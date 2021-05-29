import React from 'react';
import './Recommendations.css';
import Footer from './Footer.js';
import Result from './Result.js';
import Nav from './Nav.js';

export default class Recommendations extends React.Component {
    constructor(props){
        super(props);
        this.state = {
            items: [],
            state: "loading"
        };
    };

    componentDidMount(){
        fetch(process.env.REACT_APP_API_BASE_URL + "/tracks/recommendations", {
            credentials: 'include'
        })
        .then(res => res.json())
        .then(
            (result) => {
                // add the results to the state as 'items'
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
        // show the state of the page and footer while we're loading
        if (!this.state.items.length) {
            // TODO: make this better
            return <React.Fragment>
                <div class="body">
                    <h1>Recommendations</h1>
                    <p>To receive fresh recommendations, refresh the page.</p>
                    <div>state: {this.state.state}</div>
                </div>
            </React.Fragment>
        }

        // iterate through items received and 
        const items = this.state.items.map((i) => {
            return <Result 
                url={i.external_urls.spotify} 
                image={i.album.images[0].url} 
                artist={i.artists[0].name} 
                name={i.name} 
            />;
        });

        // TODO: why am I doing this?
        let recs = []
        for (var i = 0; i < items.length; i++) {
            recs.push(items[i]);
        }

        return <React.Fragment>
                <div className="body">
                    <div className="flex-table">
                        {recs}
                    </div>
                </div>
            </React.Fragment>
    }
}