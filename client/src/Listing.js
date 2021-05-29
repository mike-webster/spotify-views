import React from 'react';
import './Listing.css';
import Footer from './Footer.js';
import Result from './Result.js';

export default class Listing extends React.Component {
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
            }
        )
    };

    render(){
        let ret = [];
        
        // adding page title copy
        ret.push(<h1>Recommendations</h1>)
        ret.push(<p>To receive fresh recommendations, refresh the page.</p>)

        // show the state of the page and footer while we're loading
        if (!this.state.items.length) {
            ret.push(<div>{this.state.state}</div>);
            ret.push(<Footer />);
            return (ret)
        }

        // iterate through items received and 
        const items = this.state.items.map((i) => {
            return (<Result url={i.external_urls.spotify} image={i.album.images[0].url} artist={i.artists[0].name} name={i.name} />);
        });

        let recs = []
        for (var i = 0; i < items.length; i++) {
            recs.push(items[i]);
        }
        {console.log("adding table")}
        ret.push(<React.Fragment><div className="flex-table">{recs}</div></React.Fragment>)

        ret.push(<Footer />);

        return ret
    };
};