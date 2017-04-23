window.onload = main;

var ws;

function main() {
    var oBtn = document.getElementById('btn1');
    oBtn.onclick = OnButton1;
}

function OnButton1() {
    ws = new WebSocket('ws://localhost/ws');

    ws.onopen = OnOpen;
    ws.onmessage = OnMessage;
}

function OnOpen(event) {
    ws.send('hello websocket');
}

function OnMessage(event) {
    var oTxt = document.getElementById('txt1');
    oTxt.value = event.data;
};