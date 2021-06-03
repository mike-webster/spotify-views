import React from 'react';
import './Nav.css';
import './Shared.css';
import img from './logo.png';

export default class Nav extends React.Component {
    constructor(){
        super();
        this.expand = this.expand.bind(this);
    };
    expand = () => {
        var nav = document.getElementById("nav");
        if (nav == null) {
            console.log("error");
            return;
        }

        let btn = document.getElementsByClassName("menu-btn")[0];
        let action = "";
        let state = nav.getAttribute("data-expanded");
        if (state == null || state == "false") {
            action = "expand";
        } else {
            action = "collapse";
        }

        // iterating throuhg each li in the nav
        for (let index = 1; index < nav.children[0].children.length; index++) {
            const element = nav.children[0].children[index];
            if (window.getComputedStyle(element).display === "none") {
                nav.children[0].children[index].style.display = "flex";
            } else {
                nav.children[0].children[index].style.display = "none";
            }
        }

        let body = document.getElementsByClassName("body")[0];
        let header = document.getElementsByClassName("nav")[0];

        if (action == "expand") {
            body.style.height = "80vh"
            header.style.height = "15vh"
            nav.setAttribute("data-expanded", "true")
        } else {
            body.style.height = "88vh"
            header.style.height = "7vh"
            nav.setAttribute("data-expanded", "false");
        }
    };
    render(){
        return(
                <div className="nav" id="nav" data-expanded="false">
                    <ul className="links">
                        <li className="navLink menu-btn">
                            <a href="javascript:void(0);" onClick={this.expand} ><i className="fa fa-bars"></i></a>
                        </li>
                        <li className="navLink">
                            <a href="/">Word Cloud</a>
                        </li>
                        <li className="navLink">
                            <a href="/discover">Tops</a>
                        </li>
                        <li className="navLink">
                            <a href="/">Library</a>
                        </li>
                    </ul>
                    <div className="img" style={{ backgroundImage: `url(${img})` }}></div>

                </div>
        );
    };
}