var onHardMode = false;
document.getElementById("easy").addEventListener("click", function(event) {
  event.preventDefault();
  if (onHardMode) {
    stopHardMode();
    // alert("changed mode to easy");
  }
}, false);

document.getElementById("hard").addEventListener("click", function(event) {
  event.preventDefault();
  if (!onHardMode) {
    startHardMode();
    // alert("changed mode to hard");
  }
}, false);

function startHardMode() {
  getMovableItems();
  window.addEventListener("mousemove", runaway, false);
  onHardMode = true;
  document.getElementsByTagName("html")[0].className = "hardmode";
}

function stopHardMode() {
  getMovableItems("reset");
  window.removeEventListener("mousemove", runaway, false);
  onHardMode = false;
  document.getElementsByTagName("html")[0].className = "easymode";
}

var movableItems;
var movableItemShadows;
function getMovableItems(mode) {
  if (mode === "reset") {
    for (var i = 0; i < movableItems.length; i++) {
      if (movableItems[i].dataset.top) {
        movableItems[i].style.top = movableItems[i].dataset.top;
      } else {
        movableItems[i].style.top = "";
      }
      if (movableItems[i].dataset.left) {
        movableItems[i].style.left = movableItems[i].dataset.left;
      } else {
        movableItems[i].style.left = "";
      }
      var node = document.getElementById("movableItemShadow" + String(i));
      if (node && node.parentNode) {
        node.parentNode.removeChild(node);
      } else {
        console.log("ERROR");
      }
    }
    return;
  }
  var movables = [
    document.getElementsByTagName("p"),
    document.getElementsByTagName("a"),
    document.getElementsByTagName("li"),
    document.getElementsByTagName("h1"),
    document.getElementsByTagName("h2"),
    document.getElementsByTagName("h3"),
    document.getElementsByTagName("h4"),
    document.getElementsByTagName("h5"),
    document.getElementsByTagName("h6"),
    document.getElementsByClassName("movable"),
  ];

  movableItems = [];
  movableItemShadows = [];
  for (var h = 0; h < movables.length; h++) {
    var mm = movables[h];
    for (var i = 0; i < mm.length; i++) {
      var el = mm.item(i);
      movableItems.push(el);
      var box = el.getBoundingClientRect();
      if (el.style.top) {
        el.dataset.top = el.style.top;
      } else {
        el.style.top = box.top + "px";
      }
      if (el.style.left) {
        el.dataset.left = el.style.left;
      } else {
        el.style.left = box.left + "px";
      }
      el.style.width = box.width + "px";
      el.style.height = box.height + "px";
      el.dataset.width = box.width;
      el.dataset.height = box.height;

      var shadow = document.createElement(el.tagName);
      shadow.id = "movableItemShadow" + movableItemShadows.length;
      shadow.className = "shadow";
      shadow.style.width = box.width + "px";
      shadow.style.height = box.height + "px";
      if (el.dataset.top) {
        shadow.style.position = "absolute";
        shadow.style.top = el.dataset.top;
      }
      if (el.dataset.left) {
        shadow.style.position = "absolute";
        shadow.style.left = el.dataset.left;
        shadow.style.backgroundColor = "white";
      }
      movableItemShadows.push(shadow);
    }
  }
  for (var i = 0; i < movableItemShadows.length; i++) {
    movableItems[i].parentElement.insertBefore(movableItemShadows[i], movableItems[i]);
  }
}

var movingHandle;

function runaway(event) {
  var mouseLocation = {
    left: event.clientX,
    right: event.clientX,
    top: event.clientY,
    bottom: event.clientY
  };
  var maxTop = document.getElementsByTagName("html").item(0).getBoundingClientRect().height - 100;
  var maxLeft = document.documentElement.clientWidth - 100;
  for (var i = 0; i < movableItems.length; i++) {
      var obj = movableItems[i];
      if (obj.dataset.moving) {
        return;
      }
      var rect = obj.getBoundingClientRect();
      var relPos = relativePosition(mouseLocation, rect);
      if (relPos.distance < 100) {
        var mvtLeft = 0;
        if (relPos.isOnLeft) {
          mvtLeft = 2 * (100-relPos.distance);
        } else if (relPos.isOnRight) {
          mvtLeft = -2 * (100-relPos.distance);
        }
        var mvtTp = 0;
        if (relPos.isOnTop) {
          mvtTp = 2 * (100-relPos.distance);
        } else if (relPos.isOnBottom) {
          mvtTp = -2 * (100-relPos.distance);
        }
        if (mvtTp === 0 && mvtLeft === 0) return;
        var lft = Number(obj.style.left.slice(0, -2));
        var tp = Number(obj.style.top.slice(0, -2));
        var newCoords = [lft + mvtLeft, tp + mvtTp];
        if (newCoords[0] < 100) {
          newCoords[0] = 100;
        } else if (newCoords[0] + obj.dataset.width > maxLeft) {
          newCoords[0] = maxLeft - obj.dataset.width;
        }
        if (newCoords[1] < 100) {
          newCoords[1] = 100;
        } else if (newCoords[1] + obj.dataset.height > maxTop) {
          newCoords[1] = maxTop - obj.dataset.height;
        }
        obj.style.left = String(newCoords[0]) + "px";
        obj.style.top = String(newCoords[1]) + "px";

        obj.dataset.moving = "true";
        movingHandle = window.setTimeout(function (event) {
          for (var j = 0; j < movableItems.length; j++) {
            movableItems[j].dataset.moving = "";
          }
        }, 500);
      }

  }
}

function relativePosition(el, rel) {
  var horiz, vert;
  var isOnLeft = false;
  var isOnTop = false;
  var isOnRight = false;
  var isOnBottom = false;
  if (el.left <= rel.left + 10) {
    if (el.right >= rel.left + 10) {
      //on left, but right side is not
      horiz = 0;
    } else {
      //fully on left
      horiz = rel.left - (el.right + 10);
      isOnLeft = true;
    }
  } else {
    if (el.left <= rel.right - 10) {
      horiz = 0;
    } else {
      horiz = el.left - (rel.right - 10);
      isOnRight = true;
    }
  }
  if (el.top <= rel.top + 10) {
    if (el.bottom >= rel.top + 10) {
      vert = 0;
    } else {
      vert = rel.top - (el.bottom + 10);
      isOnTop = true;
    }
  } else {
    if (el.top <= rel.bottom + 10) {
      vert = 0;
    } else {
      vert = el.top - (rel.bottom + 10);
      isOnBottom = true;
    }
  }
  return {
    distance: Math.sqrt(horiz * horiz + vert * vert),
    isOnLeft: isOnLeft,
    isOnRight: isOnRight,
    isOnTop: isOnTop,
    isOnBottom: isOnBottom
  };
}
