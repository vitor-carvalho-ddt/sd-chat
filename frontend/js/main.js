window.onload = function () {
  var conn;
  var username = document.getElementById("username");
  var msg = document.getElementById("msg");
  var log = document.getElementById("log");

  function appendLog(item) {
    var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
    log.appendChild(item);
    if (doScroll) {
      log.scrollTop = log.scrollHeight - log.clientHeight;
    }
  }

  document.getElementById("form").onsubmit = function () {
    if (!conn || !username.value || !msg.value) {
      return false;
    }

    var payload = {
      username: username.value,
      message: msg.value
    };

    conn.send(JSON.stringify(payload));
    msg.value = "";
    return false;
  };

  if (window["WebSocket"]) {
    conn = new WebSocket("ws://" + document.location.host + "/ws");

    conn.onclose = function (evt) {
      var item = document.createElement("div");
      item.innerHTML = "<b>Connection closed.</b>";
      appendLog(item);
    };

    conn.onmessage = function (evt) {
      try {
        const data = JSON.parse(evt.data);
        const formatted = `${data.username}: ${data.message}`;

        const item = document.createElement("div");
        item.innerText = formatted;
        appendLog(item);
      } catch (err) {
        const item = document.createElement("div");
        item.innerText = evt.data;
        appendLog(item);
      }
    };

  } else {
    var item = document.createElement("div");
    item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
    appendLog(item);
  }
};

