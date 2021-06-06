import React from 'react';
import './Shared.css';
import Result from './Result.js';

export default class Tops extends React.Component {
    constructor(props){
        super(props);
        this.state = {
            items: [],
            state: "loading",
            sort: "short_term"
        };

        this.changeSort = this.changeSort.bind(this);
        this.fetchData = this.fetchData.bind(this);
        this.fetchGenreData = this.fetchGenreData.bind(this);
        this.getGenreRow = this.getGenreRow.bind(this);
        this.getArtistRow = this.getArtistRow.bind(this);
        this.getTracksRow = this.getTracksRow.bind();
    };

    componentDidMount(){
        console.log("rendering")
        this.fetchData();
    };

    fetchData = () => {
        console.log("fetching data for " + this.props.focus);

        if (this.props.focus === "genres") {
            this.fetchGenreData();
            return
        }

        let url = process.env.REACT_APP_API_BASE_URL;
        url += "/" + this.props.focus + "/top?time_range=" + this.state.sort;
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
                    items: tmp
                }, console.log("items: ", this.state.items));
                
            },
            (error) => {
                // TODO: something in this error state
                this.setState({
                    state: "error",
                    error: error,
                    items: []
                });
                console.log(error);
                console.log("redirecting");
                //window.location.href = "/";
            }
        );
    };

    fetchGenreData = () => {
        let url = process.env.REACT_APP_API_BASE_URL;
        url += "/genres?time_range=" + this.state.sort;
        let totals = {};

        fetch(url, {credentials: 'include'})
        .then(res => res.json())
        .then(
            (result) => {
                console.log(result);
                // add the results to the state as 'items'
                let tmp = [];
                for (var i = 0; i < result.length; i++) {
                    tmp.push(result[i].Key);
                    if (totals[result[i].Key] == null) {
                        // new key
                        totals[result[i].Key] = result[i].Value
                    } else {
                        // add to existing key
                        totals[result[i].Key] += result[i].Value
                    }
                }

                this.setState({
                    state: "success",
                    items: tmp
                });
            },
            (error) => {
                this.setState({
                    state: "error",
                    error: error,
                    items: []
                });
                console.log(error);
            }
        );
    };

    changeSort = (e) => {
        this.setState({sort: e.target.value}, ()=>{
            // do this in the callback to make sure we wait
            // for the state to change
            this.fetchData();
        });
    };

    getTracksRow(i) {
        return <Result 
            url={i.external_urls.spotify} 
            image={i.album.images[0].url} 
            artist={i.artists[0].name} 
            name={i.name} 
        />; 
    };

    getArtistRow(i) {
        return <Result 
            url={i.external_urls.spotify} 
            image={ (i.images != null) ? i.images[0].url : i.album.images[0].url } 
            artist={i.name} 
        />
    };

    getGenreRow(i) {
        <Result 
            artist={i} 
        />
    };

    render(){
        // show the state of the page and footer while we're loading
        if (this.state.items.length > 0) {
            // TODO: make this better
            return <div className="body">
                <div>state: {this.state.state}</div>
            </div>
        }

        // iterate through items received and 
        const items = this.state.items.map((i) => {
            if (this.props.focus === "tracks") {
                return this.getTracksRow(i);
            } else if (this.props.focus === "artists") {
                return this.getArtistRow(i);
            } else if (this.props.focus === "genres") {
                return this.getGenreRow(i);
            }
        });

        // TODO: why am I doing this?
        let recs = []
        for (var i = 0; i < items.length; i++) {
            recs.push(items[i]);
        }

        return <React.Fragment>
            <div className="body">
                <select value={this.state.sort} onChange={this.changeSort}>
                    <option value="Recent">Recent</option>
                    <option value="In Between">In Between</option>
                    <option value="Going Way Back">Going Way Back</option>
                </select>
                <div className="flex-table">
                    {recs}
                </div>
            </div>
        </React.Fragment>
    }
}