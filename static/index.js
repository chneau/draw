window.onload = () => {
    let kelly = ["#e6194b", "#3cb44b", "#ffe119", "#0082c8", "#f58231", "#911eb4", "#46f0f0", "#f032e6", "#d2f53c", "#fabebe", "#008080", "#e6beff", "#aa6e28", "#fffac8", "#800000", "#aaffc3", "#808000", "#ffd8b1", "#000080", "#808080", "#FFFFFF", "#000000"];
    let lineWidths = [5, 10, 20, 40, 80, 160, 300, 500, 1000, 10000];
    let canvas = document.getElementById('canvas');
    let legend = document.getElementById('legend');
    let ctx = canvas.getContext("2d");
    ctx.canvas.width = window.innerWidth;
    ctx.canvas.height = window.innerHeight;
    let socket = new WebSocket(`ws://${location.host}/ws`);
    let color = 0;
    let indexLineWidth = 1;
    var previous = null;


    socket.onmessage = (event) => {
        let msg = JSON.parse(event.data);
        ctx.strokeStyle = kelly[msg.c];
        ctx.lineWidth = lineWidths[msg.w];
        ctx.lineCap = "round";
        ctx.beginPath();
        ctx.moveTo(msg.s[0], msg.s[1]);
        ctx.lineTo(msg.s[2], msg.s[3]);
        ctx.stroke();
    };


    canvas.onmousemove = (event) => {
        if (event.which == 1) {
            if (previous != null) {
                socket.send(JSON.stringify({ s: [previous.pageX, previous.pageY, event.pageX, event.pageY], c: color, w: indexLineWidth }));
            }
            previous = event;
            return;
        }
        previous = null;
    };


    canvas.onmouseup = (event) => {
        if (event.which == 3) {
            ++indexLineWidth;
        }
        if (event.which == 2) {
            --indexLineWidth;
            if (indexLineWidth < 0) {
                indexLineWidth = lineWidths.length - 1;
            }
        }
        indexLineWidth = indexLineWidth % lineWidths.length;
        legend.innerHTML = lineWidths[indexLineWidth];
        legend.style = `color: ${kelly[color]};`;
    };


    canvas.oncontextmenu = (event) => {
        event.preventDefault();
    };


    canvas.onmousewheel = (event) => {
        event.preventDefault();
        if (event.deltaY < 0) {
            ++color;
        } else {
            --color;
            if (color < 0) {
                color = kelly.length - 1;
            }
        }
        color = color % kelly.length;
        legend.innerHTML = lineWidths[indexLineWidth];
        legend.style = `color: ${kelly[color]};`;
    };

    document.onkeydown = function (evt) {
        switch (evt.keyCode) {
            case 87: // w
            case 38: // up
                ++color;
                break;
            case 83: // s
            case 40: // down
                --color;
                if (color < 0) {
                    color = kelly.length - 1;
                }
                break;
            case 65: // a
            case 37: // left
                --indexLineWidth;
                if (indexLineWidth < 0) {
                    indexLineWidth = lineWidths.length - 1;
                }
                break;
            case 68: // d
            case 39: // right
                ++indexLineWidth
                break;
            default:
                break;
        }
        indexLineWidth = indexLineWidth % lineWidths.length;
        color = color % kelly.length;
        legend.innerHTML = lineWidths[indexLineWidth];
        legend.style = `color: ${kelly[color]};`;
    };

    legend.innerHTML = "use your mouse to draw";
    legend.style = `color: ${kelly[color]};`;
}
