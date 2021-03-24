// something I found for document.ready
Document.prototype.ready = function(callback) {
    if(callback && typeof callback === 'function') {
        document.addEventListener("DOMContentLoaded", function() {
        if(document.readyState === "interactive" || document.readyState === "complete") {
            return callback();
        }
        });
    }
};

function setCookie(cname, cvalue, exdays) {
    var d = new Date();
    d.setTime(d.getTime() + (exdays*24*60*60*1000));
    var expires = "expires="+ d.toUTCString();
    document.cookie = cname + "=" + cvalue + ";" + expires + ";path=/";
}

function getCookie(cname) {
    var name = cname + "=";
    var decodedCookie = decodeURIComponent(document.cookie);
    var ca = decodedCookie.split(';');
    for(var i = 0; i <ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length);
        }
    }
    return "";
}

document.ready(function(){
    ck = getCookie("cookie-banner")
    if (ck != "") {
        return
    }

    btn = document.getElementById("cookie-close")
    btn.onclick = function(){
        // make the close button work
        document.getElementById("cookie-banner").style.display = "none";
        setCookie("cookie-banner", "1", 30);
    };

    // if we don't have a cookie set for the cookie banner, we'll need
    // to show the thing.
    document.getElementById("cookie-banner").style.display = "block";
});

function topnavClick() {
    var x = document.getElementById("svTopNav");
    if (x.className === "topnav") {
      x.className += " responsive";
    } else {
      x.className = "topnav";
    }
}

/* Toggle between adding and removing the "responsive" class to topnav when the user clicks on the icon */
function myFunction() {
    var x = document.getElementById("myTopnav");
    if (x.className === "topnav") {
      x.className += " responsive";
    } else {
      x.className = "topnav";
    }
  }