function hideAll() {
    document.querySelectorAll('[data-item]').forEach(function(item){
        item.removeAttribute("data-item-visible");
        item.setAttribute("data-item-invisible", "");
    });
}

function displaySome(start, num) {
    // num = num * 3; // I did this because the data-item used to be a div container but now its on the items
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

document.addEventListener("DOMContentLoaded", function(){
    var rows = document.querySelectorAll('[data-item]');
    var prev = document.getElementById("prev");
    var nxt = document.getElementById("next");
    var page = parseInt(document.getElementById("page").getAttribute('value'));
    var perPage = 10;
    var rowLength = rows.length /3;
    var pageLimit = parseInt(rowLength / perPage);

    if ((rowLength - 1) % perPage != 0) {
        pageLimit++;
    } 
    nxt.addEventListener("click", function(){
        hideAll();
        var page = parseInt(document.getElementById("page").getAttribute('value'));
        var perPage = 30;

        var starting = perPage * (1 + page);

        var rows = document.querySelectorAll("[data-item]");
        if (starting + perPage > rowLength) {
            // next page will run out of room
            displaySome(starting, rowLength - starting);
        } else {
            displaySome(starting, perPage);
        }

        document.getElementById("page").setAttribute("value", page + 1);
    });
    prev.addEventListener("click", function(){
        hideAll();
        var page = parseInt(document.getElementById("page").getAttribute('value'));
        var perPage = 30;

        var starting = perPage * (page-1);

        var rows = document.querySelectorAll("[data-item]");
        if (page == 0) {
            // next page will run out of room
            displaySome(0, perPage);
            document.getElementById("page").setAttribute("value", page);
            return;
        } else {
            displaySome(starting, perPage);
        }

        document.getElementById("page").setAttribute("value", page - 1);
    });

    hideAll();
    displaySome(0,30);

    var toggle = document.getElementById("toggle");
    var qs = window.location.search.substring(1);
    if (qs.includes("sort=asc")) {
        toggle.setAttribute("href", "/library/tempo");
        toggle.innerHTML = "switch to DESC";
    } else {
        toggle.setAttribute("href","/library/tempo?sort=asc");
        toggle.innerHTML = "switch to ASC";
    }
});