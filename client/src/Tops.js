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
    };

    changeSort = (e) => {
        let cur = this.state.sort;
        let upd = e.target.value;
        this.setState({sort: upd});
        this.fetchData();
    };

    fetchData = () => {
        let url = process.env.REACT_APP_API_BASE_URL + "/tracks/top" + "?time_range=" + this.state.sort;
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
                });
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

    componentDidMount(){
        console.log("rendering")
        this.fetchData();
    };

    render(){
        // show the state of the page and footer while we're loading
        if (!this.state.items.length) {
            // TODO: make this better
            return <React.Fragment>
                <div className="body">
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