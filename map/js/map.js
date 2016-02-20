'use strict';

/** @type {?Two} */
var two;

/** @type {!Array<!System>} */
var systems = [];

function initTwo() {
  // Make a two.js instance and place it on the page.
  var elem = document.getElementById('draw-shapes');
  var params = { width: 400, height: 400 };
  two = new Two(params).appendTo(elem);
  loadMapData();
}
document.addEventListener('DOMContentLoaded', initTwo);

function loadMapData() {
  var req = new XMLHttpRequest();
  req.addEventListener('load', handleMapData);
  req.open('GET', '/static/systems_sample_large.json');
  req.send();
}

function handleMapData() {
  var data = JSON.parse(this.responseText);
  for (var i in data) {
    systems.push(System.fromJSON(data[i]));
  }
  var bounds = findBounds();
}

function findBounds() {
  var bounds = {
    min_x: null,
    max_x: null,
    min_y: null,
    max_y: null,
  };
  for (var i in systems) {
    var system = systems[i];
    if (bounds.min_x == null || bounds.min_x > system.x) {
      bounds.min_x = system.x;
    }
    if (bounds.max_x == null || bounds.max_x < system.x) {
      bounds.max_x = system.x;
    }
    if (bounds.min_y == null || bounds.min_y > system.y) {
      bounds.min_y = system.y;
    }
    if (bounds.max_y == null || bounds.max_y < system.y) {
      bounds.max_y = system.y;
    }
  }
  return bounds;
}


/**
 * @constructor
 */
System = function() {
  this.id = null;
  this.name = null;
  this.x = null;
  this.y = null;
  this.z = null;
  this.faction = null;
  this.population = null;
  this.government = null;
  this.allegiance = null;
  this.state = null;
  this.security = null;
  this.primary_economy = null;
  this.power = null;
  this.power_state = null;
  this.needs_permit = null;
  this.simbad_ref = null;
};

/**
 * Parse a JSON object and create a system.
 * @param {!Object} obj
 * @return {System}
 */
System.fromJSON = function(obj) {
  var s = new System();
  s.id = obj.id;
  s.name = obj.name;
  s.x = obj.x;
  s.y = obj.y;
  s.z = obj.z;
  s.faction = obj.faction;
  s.population = obj.population;
  s.government = obj.government;
  s.allegiance = obj.allegiance;
  s.state = obj.state;
  s.security = obj.security;
  s.primary_economy = obj.primary_economy;
  s.power = obj.power;
  s.power_state = obj.power_state;
  s.needs_permit = obj.needs_permit;
  s.simbad_ref = obj.simbad_ref;

  return s;
};
