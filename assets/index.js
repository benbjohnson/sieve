var Sieve = {};

(function() {

// The database containing all data streamed from the server.
var db = {rows:[]}

// A flag stating if the browser is currently connected via SSE.
var connected = false;

var svg = d3.select(".chart").append("svg");
var g = {
  line: svg.append("g"),
  axes: {
    x: svg.append("g").attr("class", "x axis"),
    y: svg.append("g").attr("class", "y axis"),
  }
};

var path = g.line.append("path");

var margin = {top: 20, right: 20, bottom: 30, left: 50};

// Subscribes to the data stream from the server using SSE.
Sieve.subscribe = function() {
  if (!window.EventSource) {
    alert("Please use a browser that supports SSE.")
    return;
  }

  // Open the connection and stream in data.
  var source = new EventSource('/subscribe');
  source.addEventListener('message', function(e) {
    // Parse message into JSON and convert timestamp.
    var data = JSON.parse(e.data)
    data.timestamp = moment(data.timestamp).toDate();

    // Add to database.
    db.rows.push(data)

    // Update visualization.
    Sieve.update();
  }, false);

  source.addEventListener('open', function(e) {
    connected = true;
  }, false);

  source.addEventListener('error', function(e) {
    if (e.readyState == EventSource.CLOSED) {
      connected = false;
    }
  }, false);
};

// Updates the chart view.
Sieve.update = function() {
  var chart = d3.select(".chart");
  var width = $(".chart").width() - margin.left - margin.right;
  var height = 300 - margin.top - margin.bottom;

  var x = d3.time.scale()
    .range([0, width])
    .domain(d3.extent(db.rows, function(d) { return d.timestamp; }));
  var y = d3.scale.linear()
    .range([height, 0])
    .domain(d3.extent(db.rows, function(d) { return d.value; }));

  var xAxis = d3.svg.axis().scale(x).orient("bottom");
  var yAxis = d3.svg.axis().scale(y).orient("left");

  var line = d3.svg.line()
      .x(function(d) { return x(d.timestamp); })
      .y(function(d) { return y(d.value); });

  // Setup dimensions on SVG and G.
  svg.attr("width", width + margin.left + margin.right)
    .attr("height", height + margin.top + margin.bottom)
  g.line.attr("transform", "translate(" + margin.left + "," + margin.top + ")");
  g.axes.x.attr("transform", "translate(" + margin.left + "," + (height+margin.top) + ")");
  g.axes.y.attr("transform", "translate(" + margin.left + "," + margin.top + ")");

  g.axes.x.call(xAxis);
  g.axes.y.call(yAxis)

  path
      .datum(db.rows)
      .attr("class", "line")
      .attr("d", line);
};

})();
