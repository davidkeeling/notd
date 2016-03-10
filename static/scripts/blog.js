var buttons = document.querySelectorAll(".fbshare");
for (var i = 0; i < buttons.length; i++) {
  var button = buttons.item(i);
  fbShare(button);
}

function fbShare(button) {
  button.addEventListener("click", function(event) {
    FB.ui({
      method: "share",
      href: this.dataset.url,
    })
  })
}
