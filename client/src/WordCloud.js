
    // <script src="/web/css/sortable.min.js"></script>

import React from "react"

export default class WordCloud extends React.Component {
    constructor(){
        super();

        this.hideAll = this.hideAll.bind(this);
        this.displaySome = this.displaySome.bind(this);
        this.loadCloud = this.loadCloud.bind(this);
        this.addWord = this.addWord.bind(this);
        this.state = {words: []}
    }

    hideAll() {
        document.querySelectorAll('[data-item]').forEach(function(item){
            item.removeAttribute("data-item-visible");
            item.setAttribute("data-item-invisible", "");
        });
    }

    displaySome(start, num){
        console.log("displaying some: " + start + " - " + num);
        var rows = document.querySelectorAll('[data-item]');
        for(let i = start; i < start + num; i++) {
            try {
                rows[i].setAttribute("data-item-visible", "");
                rows[i].removeAttribute("data-item-invisible");
                console.log("showing: ", i);
            } catch(ex) {
                console.log("couldn't display item");
                console.log(ex);
            }
        }
    }

    loadCloud(data){
        data.maps.forEach(this.addWord);
        WordCloud(document.getElementById("newcloud"), { 
            list: this.state.words,
            gridSize: Math.round(16 * document.getElementById('newcloud').offsetWidth / 1024),
        });
        document.getElementById("loader").style.display = "none";
    }

    addWord(value) {
        let words = this.state.words;
        words.push([value.Key, value.Value])
        this.setState({words: words});
    }

    getData() {
        return fetch(process.env.REACT_APP_API_BASE_URL + '/wordcloud/data', {
            method: "GET",
            credentials: "include"
        })
        .then(
            (response) => {
                console.log(response.body);
                return response;
            }
        )
        .then(response => response.json())
        .then(
            (body) => {
                console.log(body);
                return body;
            },
            (error) => {
                console.log(error);
            }
        )
        .then(
            (result) => {
                console.log(result);
                this.loadCloud(result);
            },
            (error) => {
                console.log(error);
            }
        )
        .catch(err => console.log(err));
    }

    componentDidMount = () => {
        this.getData();
    }

    render(){
        return(
            <React.Fragment>
                <script src="/web/css/wordcloud2.js"></script>
                <div className="content">
                    <div style={{width:"100%", textAlign:"center"}}>
                        <h1>Word Cloud</h1>
                        <p>We use the lyrics from your most popular songs and find the 50 most common words to generate this for you.</p>
                        <div id="loader" className="loader"></div>
                        <br />
                        <div id="newcloud" style={{width: "500px", height: "500px", display: "inline-block"}}></div>
                    </div>
                </div> 
            </React.Fragment>
        )
    };
}