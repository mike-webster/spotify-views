import React from 'react';
import Result from './Result.js';

export default class SpotifyOAuth extends React.Component {
    componentDidMount(){
        let qs = window.location.search.replace("?", "");
        let pairs = qs.split("&");
        let code = "";
        for (var i = 0; i < pairs.length; i++) {
            let splits = pairs[i].split("=");
            if (splits[0] == "code") {
                code = splits[1];
            }
        }

        if (code === "") {
            console.log("code not found or the user did not grant access");
            return
        }

        let url = process.env.REACT_APP_API_BASE_URL;
        url += "/spotify/token";
        const data = URLSearchParams();
        data.append("code", code);

        fetch(url)
        .then(res => res.json())
        .then(
            (result) => {
                console.log("success");
                console.log(result);
            },
            (error) => {
                console.log("error");
                console.log(error);
            }
        )
    };

    render(){
        return <p>...contacting spotify for a token...</p>
    };
}