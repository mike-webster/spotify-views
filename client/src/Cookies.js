import React from 'react';

export default class Cookies extends React.Component {
    render() {
        return(
            <div id="cookie-banner" class="cookie-banner">
                By using this website, you agree to our use of cookies. We use cookies to improve the user experience. <span id="cookie-close" class="cookie-close">X</span>
            </div>
        );
    };
}