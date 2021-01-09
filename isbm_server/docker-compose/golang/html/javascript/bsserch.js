var camflag = false;
function camch(data){
    if (camflag){
        Quagga.stop();
        document.getElementById("interactive").style.display = "none"
        camflag=false;
        data.value = "start"
    }
    else{
        Quagga.init({
            inputStream: { type : 'LiveStream' },
            decoder: {
                readers: [{
                    format: 'ean_reader',
                    config: {}
                }]
            }
        }, (err) => {
        
            if(!err) {
        
                Quagga.start();
                document.getElementById("interactive").style.display = ""
                camflag = true;
                data.value = "stop"

        
            }
        
        });
    }
}
Quagga.init({
    inputStream: { type : 'LiveStream' },
    decoder: {
        readers: [{
            format: 'ean_reader',
            config: {}
        }]
    }
}, (err) => {

    if(!err) {

        Quagga.start();
        camflag = true;

    }

});
Quagga.onProcessed(function(result) {
    var drawingCtx = Quagga.canvas.ctx.overlay,
        drawingCanvas = Quagga.canvas.dom.overlay;

    if (result) {
        if (result.boxes) {
            drawingCtx.clearRect(0, 0, parseInt(drawingCanvas.getAttribute("width")), parseInt(drawingCanvas.getAttribute("height")));
            result.boxes.filter(function (box) {
                return box !== result.box;
            }).forEach(function (box) {
                Quagga.ImageDebug.drawPath(box, {x: 0, y: 1}, drawingCtx, {color: "green", lineWidth: 2});
            });
        }

        if (result.box) {
            Quagga.ImageDebug.drawPath(result.box, {x: 0, y: 1}, drawingCtx, {color: "#00F", lineWidth: 2});
        }

        if (result.codeResult && result.codeResult.code) {
            Quagga.ImageDebug.drawPath(result.line, {x: 'x', y: 'y'}, drawingCtx, {color: 'red', lineWidth: 3});
        }
    }
});

Quagga.onDetected((result) => {

    var code = result.codeResult.code;
    console.log(code)

    if ((code.substr(0,3)=="978")||(code.substr(0,3)=="491")||(code.substr(0,3)=="456")){
        if (document.getElementById("keyword").value != code){
            document.getElementById("keyword").value = code
            getJSON('output',true)
        }
    }

});