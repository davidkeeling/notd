function hideFeaturedImages() {
  document.getElementById("featuredImageWrapper").style.display = "none";
}

function showFeaturedImages() {
  document.getElementById("featuredImageWrapper").style.display = "block";
}

function selectFeaturedImage() {
  hideFeaturedImages();
  document.getElementById("featuredImagePreview").src = this.src;
  document.getElementById("featuredImage-input").value = this.dataset.mediaId;
}

var showFeaturedImagesButton = document.getElementById("showFeaturedImages");

showFeaturedImagesButton.addEventListener("click", showFeaturedImages, false);

var featuredImages = document.getElementsByClassName("featuredImage");
for (var i = 0; i < featuredImages.length; i++) {
  featuredImages.item(i).addEventListener("click", selectFeaturedImage, false);
}
