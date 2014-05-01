var Sieve = {};

(function() {
// The database containing all data streamed from the server.
var db = {rows:[]}

// A flag stating if the browser is currently connected via SSE.
var connected = false;

// Subscribes to the data stream from the server using SSE.
Sieve.subscribe = function() {
  if (!window.EventSource) {
    alert("Please use a browser that supports SSE.")
    return;
  }

  // Open the connection and stream in data.
  var source = new EventSource('/subscribe');
  source.addEventListener('message', function(e) {
    console.log(e.data);
    db.rows.push(e.data)
  }, false);

  source.addEventListener('open', function(e) {
    connected = true;
  }, false);

  source.addEventListener('error', function(e) {
    if (e.readyState == EventSource.CLOSED) {
      connected = false;
    }
  }, false);
}
})();
