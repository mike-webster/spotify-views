import React from 'react';
import './Footer.css';

export default class Footer extends React.Component {
    render(){
        return (
            <div className="footer">
                <a href="mailto:spotify-views@webstercode.com">Contact</a>
                {/* <div className="blank"></div> */}
                <a href="https://www.mikewebster.io">Mike Webster</a>
                <a href="https://www.spotify.com">Spotify</a>
            </div>
        );
    };
}