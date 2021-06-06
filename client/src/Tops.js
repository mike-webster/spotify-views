import React from 'react';
import './Shared.css';
import TopTracks from './TopTracks.js';
import TopArtists from './TopArtists.js';
import TopGenres from './TopGenres.js';

export default class Tops extends React.Component {
    constructor(props){
        super(props);
        this.state = {
            state: "loading"
        };

    };

    render(){
        if (this.props.focus === "tracks") {
            console.log("loading tracks");
            return <TopTracks key="top-tracks" />
        } else if (this.props.focus === "genres") {
            console.log("loading genres");
            return <TopGenres key="top-genres" />
        } else if (this.props.focus === "artists") {
            console.log("loading artists");
            return <TopArtists key="top-artists" />
        }
    }
}