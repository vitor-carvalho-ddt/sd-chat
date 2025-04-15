let conn;

window.onload = function () {
  const loginDiv = document.getElementById("login");
  const chatDiv = document.getElementById("chat");
  const loginInput = document.getElementById("login-username");
  const loginBtn = document.getElementById("login-btn");

  const usernameField = document.getElementById("username");
  const to = document.getElementById("to");
  const msg = document.getElementById("msg");
  const log = document.getElementById("log");

  function appendLog(item) {
    const doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
    log.appendChild(item);
    if (doScroll) {
      log.scrollTop = log.scrollHeight - log.clientHeight;
    }
  }

  loginBtn.onclick = function () {
    const username = loginInput.value.trim();
    if (!username) return;

    // Show chat UI and hide login
    loginDiv.style.display = "none";
    chatDiv.style.display = "block";
    usernameField.value = username;

    // Open WebSocket
    conn = new WebSocket("ws://" + document.location.host + "/ws");

    conn.onopen = function () {
      const initPayload = {
        username: username,
        to: "server",
        message: "register"
      };
      conn.send(JSON.stringify(initPayload));
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

    conn.onclose = function () {
      const item = document.createElement("div");
      item.innerHTML = "<b>Connection closed.</b>";
      appendLog(item);
    };
  };

  document.getElementById("form").onsubmit = function () {
    if (!conn || !msg.value) return false;

    const payload = {
      username: usernameField.value,
      to: to.value || "all",
      message: msg.value
    };

    conn.send(JSON.stringify(payload));
    msg.value = "";
    return false;
  };
};

