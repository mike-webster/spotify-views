import React from 'react';
import './Shared.css';
import './Recommendations.css';
import Result from './Result.js';

export default class Recommendations extends React.Component {
    constructor(props){
        super(props);
        this.state = {
            items: [],
            state: "loading"
        };
    };

    componentDidMount(){
        let body  = document.getElementById('recommendations');
        if (body) body.innerHTML = '<div id="loader" class="loader"></div>';
        
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

                let body  = document.getElementById('recommendations');
                body.innerHTML = '';

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
                document.cookie = "sv-authed=0;path=/";
                // window.location.href = "/";
            }
        )
    };

    render(){
        // show the state of the page and footer while we're loading
        if (!this.state.items.length) {
            // TODO: make this better
            return <div id="recommendations" className="body"></div>
        }

        // iterate through items received and 
        const items = this.state.items.map((item, i) => {
            return <Result 
                key={i}
                url={item.external_urls.spotify} 
                image={item.album.images[0].url} 
                artist={item.artists[0].name} 
                name={item.name} 
            />;
        });

        // TODO: why am I doing this?
        let recs = []
        for (var i = 0; i < items.length; i++) {
            recs.push(items[i]);
        }

        return <div data-check="false" className="body">
            <div  className="flex-table">
                {recs}
            </div>
        </div>
    }
}