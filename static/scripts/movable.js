"use strict";

var dave = document.getElementById("dave");
var ben = document.getElementById("ben");
var grant = document.getElementById("grant");
var taylor = document.getElementById("taylor");
var jack = document.getElementById("jack");
var bandmates = [jack, grant, taylor, ben, dave];
var shadows = [];
for (var i = 0; i < 5; i++) {
  shadows[i] = document.getElementById("shadow-" + bandmates[i].id);
  shadows[i].style.left = bandmates[i].style.left;
  shadows[i].style.top = bandmates[i].style.top;
  shadows[i].style.width = bandmates[i].style.width;
  shadows[i].style.height = bandmates[i].style.height;
}

var map = document.getElementById("tapestry");

function moveThis(event) {
  let grabEvent = event;
  let movingEl = this;
  if (movingEl.id === "the-map") grabEvent.stopPropagation();

  let movingElRect = movingEl.getBoundingClientRect();
  let mapRect = map.getBoundingClientRect();
  let parentRect;

  let mouseOffset = [
    Number(movingEl.style.left.slice(0, -2)) - grabEvent.clientX,
    Number(movingEl.style.top.slice(0, -2)) - grabEvent.clientY,
  ];
  let steps = [];
  let i = 1;
  let stepSize = 2;
  function moveit(event) {
    i++;
    if (i % stepSize !== 0) return;
    let mouseMoveEvent = event;
    // let diff = [mouseMoveEvent.clientX - grabEvent.clientX, mouseMoveEvent.clientY - grabEvent.clientY];
    // let newMapCoords = [mapInit[0] - diff[0], mapInit[1] - diff[1]];
    let newCoords = [mouseMoveEvent.clientX + mouseOffset[0], mouseMoveEvent.clientY + mouseOffset[1]];
    // let newBounds = getBounds(diff, movingElRect);
    let newBounds = {};
    newBounds.top = mouseMoveEvent.clientY - grabEvent.clientY + movingElRect.top;
    newBounds.left = mouseMoveEvent.clientX - grabEvent.clientX + movingElRect.left;
    newBounds.bottom = newBounds.top + movingElRect.height;
    newBounds.right = newBounds.left + movingElRect.width;

    // normalizeWithin(newCoords, newBounds, boxRect);
    // normalizeWithin(newCoords, newBounds, parentRect);
    steps.push(newCoords);
    // movingEl.style.left = "" + steps[i/stepSize-1][0] + "px";
    // movingEl.style.top = "" + steps[i/stepSize-1][1] + "px";
    movingEl.style.left = "" + steps[i/stepSize-1][0] + "px";
    movingEl.style.top = "" + steps[i/stepSize-1][1] + "px";
  }

  function stop(event) {
    event.stopPropagation();
    window.removeEventListener("click", stop, true);
    window.removeEventListener("mousemove", moveit, false);
    console.log("stopped moving " + movingEl.id);
  }

  window.addEventListener("click", stop, true);
  window.addEventListener("mousemove", moveit, false);
  console.log("started moving " + movingEl.id);
}

function getBounds(coords, rect) {
  let newBounds = {};
  newBounds.top = coords[1] + rect.top;
  newBounds.left = [0] + rect.left;
  newBounds.bottom = newBounds.top + rect.height;
  newBounds.right = newBounds.left + rect.width;
  return newBounds;
}

function normalizeWithin(coords, child, parent) {
  if (child.top < parent.top) {
    coords[1] += (parent.top - child.top);
  } else if (child.bottom > parent.bottom) {
    coords[1] -= (child.bottom - parent.bottom);
  }
  if (child.right > parent.right) {
    coords[0] -= (child.right - parent.right);
  } else if (child.left < parent.left) {
    coords[0] += (parent.left - child.left);
  }
  return coords;
}

function normalizeWithinWindow(coords, child, parent) {
  if (child.top > parent.top) {
    coords[1] -= child.top - parent.top;
  } else if (child.bottom < parent.bottom) {
    coords[1] += parent.bottom - child.bottom;
  }
  if (child.right < parent.right) {
    coords[0] += parent.right - child.right;
  } else if (child.left > parent.left) {
    coords[0] -= child.left - parent.left;
  }
  return coords;
}

function relativePosition(el, rel) {
  let horiz, vert;
  let isOnLeft = false;
  let isOnTop = false;
  let isOnRight = false;
  let isOnBottom = false;
  if (el.left <= rel.left) {
    if (el.right >= rel.left) {
      horiz = 0;
    } else {
      horiz = rel.left - el.right;
      isOnLeft = true;
    }
  } else {
    if (el.left <= rel.right) {
      horiz = 0;
    }
    else {
      horiz = el.left - rel.right;
      isOnRight = true;
    }
  }
  if (el.top <= rel.top) {
    if (el.bottom >= rel.top) {
      vert = 0;
    } else {
      vert = rel.top - el.bottom;
      isOnTop = true;
    }
  } else {
    if (el.top <= rel.bottom) {
      vert = 0;
    } else {
      vert = el.top - rel.bottom;
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

var movableItems = bandmates.slice(0);

var gg = 0;
var tapRect = map.getBoundingClientRect();
window.addEventListener("mousemove", runaway, false);

function runaway(event) {
  gg++;
  // if (gg % 10 !== 0) return;
  let mouseLocation = {
    left: event.clientX,
    right: event.clientX,
    top: event.clientY,
    bottom: event.clientY
  };
  for (var i = 0; i < movableItems.length; i++) {
      let obj = movableItems[i];
      if (obj.dataset.moving) {
        return;
      }
      let rect = obj.getBoundingClientRect();
      let relPos = relativePosition(mouseLocation, rect);
      if (relPos.distance < 100) {
        let mvtLeft = 0;
        if (relPos.isOnLeft) {
          mvtLeft = 5 * (100-relPos.distance);
        } else if (relPos.isOnRight) {
          mvtLeft = -5 * (100-relPos.distance);
        }
        let mvtTp = 0;
        if (relPos.isOnTop) {
          mvtTp = 5 * (100-relPos.distance);
        } else if (relPos.isOnBottom) {
          mvtTp = -5 * (100-relPos.distance);
        }
        if (mvtTp === 0 && mvtLeft === 0) return;
        let lft = Number(obj.style.left.slice(0, -2));
        let tp = Number(obj.style.top.slice(0, -2));
        let newCoords = [lft + mvtLeft, tp + mvtTp];
        let newBounds = {};
        newBounds.left = rect.left + mvtLeft;
        newBounds.top = rect.top + mvtTp;
        newBounds.right = rect.right + mvtLeft;
        newBounds.bottom = rect.bottom + mvtTp;
        newCoords = normalizeWithin(newCoords, newBounds, tapRect);
        obj.style.left = String(newCoords[0]) + "px";
        obj.style.top = String(newCoords[1]) + "px";

        obj.dataset.moving = "true";
        window.setTimeout(function (event) {
          obj.dataset.moving = "";
        }, 200);
      }

  }
}
