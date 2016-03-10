function submitBlogPost(event) {
  event.preventDefault();
  event.stopPropagation();
  var params = {};
  var fields = ["title", "body", "postID", "userID", "featuredImage"];
  for (var field in fields) {
    var el = document.getElementById(fields[field] + "-input");
    if (!el) {
      console.log("No " + fields[field]);
      continue;
    };
    params[fields[field]] = el.value;
  }
  var dateValue = document.getElementById("created-input").value;
  var date = new Date(dateValue);
  if (dateValue && date instanceof Date) {
    params["posted"] = date.getTime();
  }
  post("/blog/post/create", params);
}

document.getElementById("create-blog-post").addEventListener("submit", submitBlogPost, false);
