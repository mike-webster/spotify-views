import React from 'react';
import Result from './Result.js';

export default class TopTracks extends React.Component {
    constructor(props){
        super(props);
        this.state = {
            items: [],
            state: "loading",
            sort: "short_term"
        };

        this.changeSort = this.changeSort.bind(this);
        this.fetchTopTracks = this.fetchTopTracks.bind(this);
    };

    componentDidMount(){
        this.fetchTopTracks();
    };

    fetchTopTracks = () => {
        let url = process.env.REACT_APP_API_BASE_URL;
        url += "/tracks/top?time_range=" + this.state.sort;
        fetch(url, {
            credentials: 'include'
        })
        .then(res => res.json())
        .then(
            (result) => {
                // add the results to the state as 'items'
                let tmp = [];
                for (var i = 0; i < result.length; i++) {
                    tmp.push(result[i])
                }
                this.setState({
                    state: "success",
                    items: tmp,
                });
                
            },
            (error) => {
                // TODO: something in this error state
                this.setState({
                    state: "error",
                    error: error
                });
                console.log(error);
                document.cookie = "sv-authed=0;path=/";
                window.location.href = "/";
            }
        );
    };

    changeSort = (e) => {
        this.setState({sort: e.target.value}, ()=>{
            // do this in the callback to make sure we wait
            // for the state to change
            this.fetchTopTracks();
        });
    };

    render(){
        // show the state of the page and footer while we're loading
        if (this.state.items.length < 1) {
            // TODO: make this better
            return <div className="body">
                <div>state: {this.state.state}</div>
            </div>
        }

        // iterate through items received and 
        const items = this.state.items.map((item, i) => {
            return <Result 
                key={"top-tracks-" + i}
                url={item.external_urls.spotify} 
                image={item.album.images[0].url} 
                artist={item.artists[0].name} 
                name={item.name} 
            />
        });

        // TODO: why am I doing this?
        let recs = []
        for (var i = 0; i < items.length; i++) {
            recs.push(items[i]);
        }

        return <div key="react-body" className="body">
            <select value={this.state.sort} onChange={this.changeSort}>
                <option value="Recent">Recent</option>
                <option value="In Between">In Between</option>
                <option value="Going Way Back">Going Way Back</option>
            </select>
            <div key="tops-data" className="flex-table">
                {recs}
            </div>
        </div>
    }
}