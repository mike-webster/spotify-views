import React from 'react';
import Result from './Result.js';

export default class TopGenres extends React.Component {
    constructor(props){
        super(props);
        this.state = {
            items: [],
            state: "loading",
            sort: "short_term"
        };

        this.changeSort = this.changeSort.bind(this);
        this.fetchGenres = this.fetchGenres.bind(this);
    };

    componentDidMount() {
        this.fetchGenres();
    };

    fetchGenres = () => {
        let url = process.env.REACT_APP_API_BASE_URL;
        url += "/genres?time_range=" + this.state.sort;
        let totals = {};
        fetch(url, {credentials: 'include'})
        .then(res => res.json())
        .then(
            (result) => {
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
                    error: error
                });
                console.log(error);
            }
        );
    };

    changeSort = (e) => {
        this.setState({sort: e.target.value}, ()=>{
            // do this in the callback to make sure we wait
            // for the state to change
            this.fetchGenres();
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
                key={"top-genre-" + i}
                artist={item} 
            />
        });

        // TODO: why am I doing this?
        let recs = []
        for (var i = 0; i < items.length; i++) {
            recs.push(items[i]);
        }

        return (
            <div key="react-body" className="body">
                <select value={this.state.sort} onChange={this.changeSort}>
                    <option value="Recent">Recent</option>
                    <option value="In Between">In Between</option>
                    <option value="Going Way Back">Going Way Back</option>
                </select>
                <div key="tops-data" className="flex-table">
                    {recs}
                </div>
            </div>
        );
    }
}