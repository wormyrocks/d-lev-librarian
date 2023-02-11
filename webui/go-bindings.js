window.onload = function() {
    let num = 0
    searchButton = document.getElementById("search_button");
    responseBox = document.getElementById("response_box");
    responseBox.innerText = "blank"
    searchButton.addEventListener("click", () => {
    addTwo(num).then(result => {
        num = result;
        console.log(num);
        responseBox.innerText = num
    })});
};