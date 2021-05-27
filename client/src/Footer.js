import React from 'react';
import './Footer.css';

export default class Footer extends React.Component {
    render(){
        return (
            <footer>
                <div class="footer">
                    <a href="mailto:spotify-views@webstercode.com">Contact</a>
                    <div class="blank"></div>
                    <a href="https://www.spotify.com">Spotify</a>
                </div>
            </footer>
        );
    };
}