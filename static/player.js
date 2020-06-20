document.addEventListener("DOMContentLoaded", function(){
    // when page is ready

    // TODO there is a bug in here somewhere.  the `substr` method is being called on nil
    // it doesn't seem to be breaking anything so it's just noise in the console I'm ignoring.

    var artists = document.getElementsByTagName("p");
    for (let i = 0; i < artists.length; i++) {
        var player = document.getElementById("i"+artists[i].getAttribute('id').substr(1))
        players = document.getElementsByTagName("iframe");
        artists[i].addEventListener("click", function() {
            // hide all
            for (let i = 0; i < players.length; i++) {
                players[i].setAttribute("style", "display:none");
            }

            // show me
            document.getElementById("i"+this.getAttribute('id').substr(1)).setAttribute("style", "");
        });
    }

    // show first one
    document.getElementsByTagName("iframe")[0].click();
});