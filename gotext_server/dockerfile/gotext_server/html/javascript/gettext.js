var HOSTURL = ""

var Yearlist = []
var Quartlist = {}
var Titlelist = {}

var Nowyear = ""
var Nowquart = ""

function getdata(outid) {
    var req = new XMLHttpRequest();		  // XMLHttpRequest オブジェクトを生成する
    req.onreadystatechange = function() {		  // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
      if(req.readyState == 4 && req.status == 200){ // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
          var data = req.responseText;
          var jdata = JSON.parse(data);
          console.log(jdata);		          // 取得した JSON ファイルの中身を表示
          tmpdata(jdata)
          document.getElementById(outid).innerHTML = viewdata()
      }else if (req.readyState == 4 && req.status != 200){ 
          var data = req.responseText;
          var jata = JSON.parse(data);
          console.log(jata);		          // 取得した JSON ファイルの中身を表示
      }
    };
    var url = HOSTURL + "/v1/text"
    req.open("GET", url, true); // HTTPメソッドとアクセスするサーバーの　URL　を指定
    req.send(null);					    // 実際にサーバーへリクエストを送信
}

function tmpdata(jdata) {
    for (var i=0;i<jdata.length;i++) {
        var tmp = jdata[i];
        var flag = true

        for (var j=0;j<Yearlist.length;j++) {
            if (Yearlist[j] == tmp.Year) {
                flag = false
                break
            }
        }
        if (flag){
            Yearlist.push(tmp.Year)
        }
        if (Quartlist[tmp.Year] == undefined) {
            Quartlist[tmp.Year] = []
        }
        flag = true
        for (var j=0;j<Quartlist.length;j++) {
            if (Quartlist[j] == tmp.Quart) {
                flag = false
                break
            }
        }
        if (flag){
            Quartlist[tmp.Year].push(tmp.Quart)
        }
        if (Titlelist[tmp.Year] == undefined) {
            Titlelist[tmp.Year] = {}
        }
        Titlelist[tmp.Year][tmp.Quart] = tmp.Title
    }
}
function viewdata() {
    var output = ""
    output += createList1()
    output += " "+createList2()
    output += "<br>"+createList3()
    return output
}

function createList1() {
    var output = ""
    output += "<select name=\"year\" id=\"year\" onchange=\"changeList2(this.value)\">"
    for (var i=0;i<Yearlist.length;i++) {
        output += "<option value=\""+Yearlist[i]+"\">"+Yearlist[i]+"</option>"
    }
    Nowyear = Yearlist[0]
    output += "</select>"
    return output
}

function createList2(){
    var output = ""
    output += "<select name=\"quart\" id=\"quart\" onchange=\"changeList3(this.value)\">"
    for (var i=0;i<Quartlist[Nowyear].length;i++) {
        output += "<option value=\""+Quartlist[Nowyear][i]+"\">"+Quartlist[Nowyear][i]+"</option>"
    }
    Nowquart = Quartlist[Nowyear][0]
    output += "</select>"
    return output

}

function changeList2(v) {
    Nowyear = v
    var data = document.getElementById("quart")
    while (0 < data.childNodes.length) {
        data.removeChild(data.childNodes[0]);
    }
    
    for (var i=0;i<Quartlist[Nowyear].length;i++) {
        const option1 = document.createElement('option');
        option1.value = Quartlist[Nowyear][i]
        option1.textContent = Quartlist[Nowyear][i]
        data.appendChild(option1)
    }
    changeList3(Quartlist[Nowyear][0])

}

function changeList3(v) {
    Nowquart = v
    var data = document.getElementById("title")
    while (0 < data.childNodes.length) {
        data.removeChild(data.childNodes[0]);
    }
    
    for (var i=0;i<Titlelist[Nowyear][Nowquart].length;i++) {
        const option1 = document.createElement('option');
        option1.value = Titlelist[Nowyear][Nowquart][i]
        option1.textContent = Titlelist[Nowyear][Nowquart][i]
        data.appendChild(option1)
    }
    titleOut(Titlelist[Nowyear][Nowquart][0])
}
function createList3() {
    var output = ""
    output += "<select name=\"title\" id=\"title\" onchange=\"titleOut(this.value)\">"
    for (var i=0;i<Titlelist[Nowyear][Nowquart].length;i++) {
        output += "<option value=\""+Titlelist[Nowyear][Nowquart][i]+"\">"+Titlelist[Nowyear][Nowquart][i]+"</option>"
    }
    output += "</select>"
    titleOut(Titlelist[Nowyear][Nowquart][0])
    return output
}
function titleOut(str) {
    var data = document.getElementById("titledata")
    data.value = str
}
