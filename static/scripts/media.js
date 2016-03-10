function selectText(element) {
  var text = document.getElementById(element),
    range,
    selection;
if (document.body.createTextRange) {
    range = document.body.createTextRange();
    range.moveToElementText(text);
    range.select();
} else if (window.getSelection) {
    selection = window.getSelection();
    range = document.createRange();
    range.selectNodeContents(text);
    selection.removeAllRanges();
    selection.addRange(range);
  }
}

function confirmDelete(event) {
  event.preventDefault();
  event.stopPropagation();
  if (!confirm("Are you sure?")) return false;
  post(this.action);
}
var deleteForms = document.getElementsByClassName("delete-form");
for (var i = 0; i < deleteForms.length; i++) {
  deleteForms.item(i).addEventListener("submit", confirmDelete, false);
}

function selectEl() {
  selectText(this.id);
}
var imgLinks = document.getElementsByClassName("img-link");
for (var i = 0; i < imgLinks.length; i++) {
  imgLinks.item(i).addEventListener("click", selectEl, false);
}
