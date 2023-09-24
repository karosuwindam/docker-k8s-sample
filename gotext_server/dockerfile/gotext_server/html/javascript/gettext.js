var HOSTURL = "";

var Yearlist = [];
var Quartlist = {};
var Titlelist = {};

var Nowyear = "";
var Nowquart = "";

function getdata(outid) {
  var req = new XMLHttpRequest(); // XMLHttpRequest オブジェクトを生成する
  req.onreadystatechange = function () {
    // XMLHttpRequest オブジェクトの状態が変化した際に呼び出されるイベントハンドラ
    if (req.readyState == 4 && req.status == 200) {
      // サーバーからのレスポンスが完了し、かつ、通信が正常に終了した場合
      var data = req.responseText;
      var jdata = JSON.parse(data);
      console.log(jdata); // 取得した JSON ファイルの中身を表示
      tmpdata(jdata);
      document.getElementById(outid).innerHTML = viewdata();
    } else if (req.readyState == 4 && req.status != 200) {
      var data = req.responseText;
      var jata = JSON.parse(data);
      console.log(jata); // 取得した JSON ファイルの中身を表示
    }
  };
  var url = HOSTURL + "/v1/text";
  req.open("GET", url, true); // HTTPメソッドとアクセスするサーバーの　URL　を指定
  req.send(null); // 実際にサーバーへリクエストを送信
}

function tmpdata(jdata) {
  for (var i = 0; i < jdata.length; i++) {
    var tmp = jdata[i];
    var flag = true;

    for (var j = 0; j < Yearlist.length; j++) {
      if (Yearlist[j] == tmp.Year) {
        flag = false;
        break;
      }
    }
    if (flag) {
      Yearlist.push(tmp.Year);
    }
    if (Quartlist[tmp.Year] == undefined) {
      Quartlist[tmp.Year] = [];
    }
    flag = true;
    for (var j = 0; j < Quartlist.length; j++) {
      if (Quartlist[j] == tmp.Quart) {
        flag = false;
        break;
      }
    }
    if (flag) {
      Quartlist[tmp.Year].push(tmp.Quart);
    }
    if (Titlelist[tmp.Year] == undefined) {
      Titlelist[tmp.Year] = {};
    }
    Titlelist[tmp.Year][tmp.Quart] = tmp.Title;
  }
}
function viewdata() {
  var output = "";
  output += createList1();
  output += " " + createList2();
  output +=
    " " +
    '<input type="text" name="search" id="search" onkeyup=\'searchTxt(this.value)\'>';
  output += "<br>" + createList3();
  return output;
}

function createList1() {
  var output = "";
  output += '<select name="year" id="year" onchange="changeList2(this.value)">';
  for (var i = 0; i < Yearlist.length; i++) {
    output +=
      '<option value="' + Yearlist[i] + '">' + Yearlist[i] + "</option>";
  }
  Nowyear = Yearlist[0];
  output += "</select>";
  return output;
}

function createList2() {
  var output = "";
  output +=
    '<select name="quart" id="quart" onchange="clearText();changeList3(this.value,\'\')">';
  for (var i = 0; i < Quartlist[Nowyear].length; i++) {
    output +=
      '<option value="' +
      Quartlist[Nowyear][i] +
      '">' +
      Quartlist[Nowyear][i] +
      "</option>";
  }
  Nowquart = Quartlist[Nowyear][0];
  output += "</select>";
  return output;
}

function changeList2(v) {
  Nowyear = v;
  var data = document.getElementById("quart");
  while (0 < data.childNodes.length) {
    data.removeChild(data.childNodes[0]);
  }

  for (var i = 0; i < Quartlist[Nowyear].length; i++) {
    const option1 = document.createElement("option");
    option1.value = Quartlist[Nowyear][i];
    option1.textContent = Quartlist[Nowyear][i];
    data.appendChild(option1);
  }
  changeList3(Quartlist[Nowyear][0], "");
}

function changeList3(v, keyword) {
  Nowquart = v;
  var data = document.getElementById("title");
  while (0 < data.childNodes.length) {
    data.removeChild(data.childNodes[0]);
  }
  var nowcount = -1;
  for (var i = 0; i < Titlelist[Nowyear][Nowquart].length; i++) {
    const option1 = document.createElement("option");
    var tmpText = Titlelist[Nowyear][Nowquart][i];
    if (
      keyword == "" ||
      tmpText.toUpperCase().indexOf(keyword.toUpperCase()) != -1
    ) {
      if (nowcount < 0) {
        nowcount = i;
      }
      option1.value = tmpText;
      option1.textContent = tmpText;
      data.appendChild(option1);
      continue;
    }
    var kana = roma2hiragana(keyword);
    if (
      keyword == "" ||
      tmpText.toUpperCase().indexOf(kana.toUpperCase()) != -1
    ) {
      if (nowcount < 0) {
        nowcount = i;
      }
      option1.value = tmpText;
      option1.textContent = tmpText;
      data.appendChild(option1);
      continue;
    }
    var katakana = hiragana2katakana(kana);
    if (
      keyword == "" ||
      tmpText.toUpperCase().indexOf(katakana.toUpperCase()) != -1
    ) {
      if (nowcount < 0) {
        nowcount = i;
      }
      option1.value = tmpText;
      option1.textContent = tmpText;
      data.appendChild(option1);
      continue;
    }
  }
  if (nowcount < 0) {
    nowcount = 0;
  }
  titleOut(Titlelist[Nowyear][Nowquart][nowcount]);
}
function createList3() {
  var output = "";
  output += '<select name="title" id="title" onchange="titleOut(this.value)">';
  for (var i = 0; i < Titlelist[Nowyear][Nowquart].length; i++) {
    output +=
      '<option value="' +
      Titlelist[Nowyear][Nowquart][i] +
      '">' +
      Titlelist[Nowyear][Nowquart][i] +
      "</option>";
  }
  output += "</select>";
  titleOut(Titlelist[Nowyear][Nowquart][0]);
  return output;
}
function titleOut(str) {
  var data = document.getElementById("titledata");
  data.value = str;
}

function clearText() {
  document.getElementById("search").value = "";
}

function searchTxt(keyword) {
  changeList3(Nowquart, keyword);
  // console.log(keyword)
}

//ローマ字をひらがなに変換する
function roma2hiragana(str) {
  var kana;
  var reples_data = [
    { from: "kk", to: "っk" },
    { from: "ss", to: "っs" },
    { from: "tt", to: "っt" },
    { from: "pp", to: "っp" },
    { from: "mm", to: "っm" },
    { from: "yy", to: "っy" },
    { from: "rr", to: "っr" },
    { from: "ww", to: "っw" },
    { from: "gg", to: "っg" },
    { from: "zz", to: "っz" },
    { from: "dd", to: "っd" },
    { from: "bb", to: "っb" },
    { from: "ff", to: "っf" },
    { from: "kk", to: "っk" },
    { from: "hh", to: "っh" },
    { from: "jj", to: "っj" },
    { from: "ll", to: "っl" },
    { from: "vv", to: "っv" },
    { from: "xx", to: "っx" },
    { from: "kya", to: "きゃ" },
    { from: "kyi", to: "きぃ" },
    { from: "kyu", to: "きゅ" },
    { from: "kye", to: "きぇ" },
    { from: "kyo", to: "きょ" },
    { from: `sha`, to: "しゃ" },
    { from: "shi", to: "し" },
    { from: "shu", to: "しゅ" },
    { from: "she", to: "しぇ" },
    { from: "sho", to: "しょ" },
    { from: "cha", to: "ちゃ" },
    { from: "chi", to: "ち" },
    { from: "chu", to: "ちゅ" },
    { from: "che", to: "ちぇ" },
    { from: "cho", to: "ちょ" },
    { from: "nya", to: "にゃ" },
    { from: "nyi", to: "にぃ" },
    { from: "nyu", to: "にゅ" },
    { from: "nye", to: "にぇ" },
    { from: "nyo", to: "にょ" },
    { from: "hya", to: "ひゃ" },
    { from: "hyi", to: "ひぃ" },
    { from: "hyu", to: "ひゅ" },
    { from: "hye", to: "ひぇ" },
    { from: "hyo", to: "ひょ" },
    { form: "pya", to: "ぴゃ" },
    { from: "pyi", to: "ぴぃ" },
    { from: "pyu", to: "ぴゅ" },
    { from: "pye", to: "ぴぇ" },
    { from: "pyo", to: "ぴょ" },
    { from: "mya", to: "みゃ" },
    { from: "myi", to: "みぃ" },
    { from: "myu", to: "みゅ" },
    { from: "mye", to: "みぇ" },
    { from: "myo", to: "みょ" },
    { from: "rya", to: "りゃ" },
    { from: "ryi", to: "りぃ" },
    { from: "ryu", to: "りゅ" },
    { from: "rye", to: "りぇ" },
    { from: "ryo", to: "りょ" },
    { from: "gya", to: "ぎゃ" },
    { from: "gyi", to: "ぎぃ" },
    { from: "gyu", to: "ぎゅ" },
    { from: "gye", to: "ぎぇ" },
    { from: "gyo", to: "ぎょ" },
    { from: "ja", to: "じゃ" },
    { from: "ji", to: "じ" },
    { from: "ju", to: "じゅ" },
    { from: "je", to: "じぇ" },
    { from: "jo", to: "じょ" },
    { from: "bya", to: "びゃ" },
    { from: "byi", to: "びぃ" },
    { from: "byu", to: "びゅ" },
    { from: "bye", to: "びぇ" },
    { from: "byo", to: "びょ" },
    { from: "pya", to: "ぴゃ" },
    { from: "pyi", to: "ぴぃ" },
    { from: "pyu", to: "ぴゅ" },
    { from: "pye", to: "ぴぇ" },
    { from: "pyo", to: "ぴょ" },
    { from: "fa", to: "ふぁ" },
    { from: "fi", to: "ふぃ" },
    { from: "fu", to: "ふ" },
    { from: "fe", to: "ふぇ" },
    { from: "fo", to: "ふぉ" },
    { from: "va", to: "ヴぁ" },
    { from: "vi", to: "ヴぃ" },
    { from: "vu", to: "ヴ" },
    { from: "ve", to: "ヴぇ" },
    { from: "vo", to: "ヴぉ" },
    { from: "tsa", to: "つぁ" },
    { from: "tsi", to: "つぃ" },
    { from: "tsu", to: "つ" },
    { from: "tse", to: "つぇ" },
    { from: "tso", to: "つぉ" },
    { from: "xtu", to: "っ" },
    { form: "xa", to: "ぁ" },
    { from: "xi", to: "ぃ" },
    { from: "xu", to: "ぅ" },
    { from: "xe", to: "ぇ" },
    { from: "xo", to: "ぉ" },
    { from: "xya", to: "ゃ" },
    { from: "xyi", to: "ぃ" },
    { from: "xyu", to: "ゅ" },
    { from: "xye", to: "ぇ" },
    { from: "xyo", to: "ょ" },
    { from: "ka", to: "か" },
    { from: "ki", to: "き" },
    { from: "ku", to: "く" },
    { from: "ke", to: "け" },
    { from: "ko", to: "こ" },
    { from: "sa", to: "さ" },
    { from: "si", to: "し" },
    { from: "su", to: "す" },
    { from: "se", to: "せ" },
    { from: "so", to: "そ" },
    { from: "ta", to: "た" },
    { from: "ti", to: "ち" },
    { from: "tu", to: "つ" },
    { from: "te", to: "て" },
    { from: "to", to: "と" },
    { from: "na", to: "な" },
    { from: "ni", to: "に" },
    { from: "nu", to: "ぬ" },
    { from: "ne", to: "ね" },
    { from: "no", to: "の" },
    { from: "ha", to: "は" },
    { from: "hi", to: "ひ" },
    { from: "hu", to: "ふ" },
    { from: "he", to: "へ" },
    { from: "ho", to: "ほ" },
    { from: "ma", to: "ま" },
    { from: "mi", to: "み" },
    { from: "mu", to: "む" },
    { from: "me", to: "め" },
    { from: "mo", to: "も" },
    { from: "ya", to: "や" },
    { from: "yi", to: "い" },
    { from: "yu", to: "ゆ" },
    { from: "ye", to: "いぇ" },
    { from: "yo", to: "よ" },
    { from: "ra", to: "ら" },
    { from: "ri", to: "り" },
    { from: "ru", to: "る" },
    { from: "re", to: "れ" },
    { from: "ro", to: "ろ" },
    { from: "wa", to: "わ" },
    { from: "wi", to: "うぃ" },
    { from: "wu", to: "う" },
    { from: "we", to: "うぇ" },
    { from: "wo", to: "を" },
    { from: "nn", to: "ん" },
    { from: "ga", to: "が" },
    { from: "gi", to: "ぎ" },
    { from: "gu", to: "ぐ" },
    { from: "ge", to: "げ" },
    { from: "go", to: "ご" },
    { from: "za", to: "ざ" },
    { from: "zi", to: "じ" },
    { from: "zu", to: "ず" },
    { from: "ze", to: "ぜ" },
    { from: "zo", to: "ぞ" },
    { from: "da", to: "だ" },
    { from: "di", to: "ぢ" },
    { from: "du", to: "づ" },
    { from: "de", to: "で" },
    { from: "do", to: "ど" },
    { from: "ba", to: "ば" },
    { from: "bi", to: "び" },
    { from: "bu", to: "ぶ" },
    { from: "be", to: "べ" },
    { from: "bo", to: "ぼ" },
    { from: "pa", to: "ぱ" },
    { from: "pi", to: "ぴ" },
    { from: "pu", to: "ぷ" },
    { from: "pe", to: "ぺ" },
    { from: "po", to: "ぽ" },
    { from: "va", to: "ヴぁ" },
    { from: "vi", to: "ヴぃ" },
    { from: "vu", to: "ヴ" },
    { from: "ve", to: "ヴぇ" },
    { from: "vo", to: "ヴぉ" },

    { from: "a", to: "あ" },
    { from: "i", to: "い" },
    { from: "u", to: "う" },
    { from: "e", to: "え" },
    { from: "o", to: "お" },
    { from: "-", to: "ー"}
  ];
  for (var i = 0; i < reples_data.length; i++) {
    if (str.indexOf(reples_data[i].from) != -1) {
      str = str.replace(
        new RegExp(reples_data[i].from, "g"),
        reples_data[i].to
      );
    }
  }
  kana = str;
  return kana;
}

// ひらがなをカタカナに変換する
function hiragana2katakana(str) {
  var kana;
  var reples_data = [
    { from: "あ", to: "ア" },
    { from: "い", to: "イ" },
    { from: "う", to: "ウ" },
    { from: "え", to: "エ" },
    { from: "お", to: "オ" },
    { from: "か", to: "カ" },
    { from: "き", to: "キ" },
    { from: "く", to: "ク" },
    { from: "け", to: "ケ" },
    { from: "こ", to: "コ" },
    { from: "さ", to: "サ" },
    { from: "し", to: "シ" },
    { from: "す", to: "ス" },
    { from: "せ", to: "セ" },
    { from: "そ", to: "ソ" },
    { from: "た", to: "タ" },
    { from: "ち", to: "チ" },
    { from: "つ", to: "ツ" },
    { from: "て", to: "テ" },
    { from: "と", to: "ト" },
    { from: "な", to: "ナ" },
    { from: "に", to: "ニ" },
    { from: "ぬ", to: "ヌ" },
    { from: "ね", to: "ネ" },
    { from: "の", to: "ノ" },
    { from: "は", to: "ハ" },
    { from: "ひ", to: "ヒ" },
    { from: "ふ", to: "フ" },
    { from: "へ", to: "ヘ" },
    { from: "ほ", to: "ホ" },
    { from: "ま", to: "マ" },
    { from: "み", to: "ミ" },
    { from: "む", to: "ム" },
    { from: "め", to: "メ" },
    { from: "も", to: "モ" },
    { from: "や", to: "ヤ" },
    { from: "いぇ", to: "イェ" },
    { from: "ゆ", to: "ユ" },
    { from: "い", to: "イ" },
    { from: "よ", to: "ヨ" },
    { from: "ら", to: "ラ" },
    { from: "り", to: "リ" },
    { from: "る", to: "ル" },
    { from: "れ", to: "レ" },
    { from: "ろ", to: "ロ" },
    { from: "わ", to: "ワ" },
    { from: "うぃ", to: "ウィ" },
    { from: "う", to: "ウ" },
    { from: "うぇ", to: "ウェ" },
    { from: "を", to: "ヲ" },
    { from: "ん", to: "ン" },
    { from: "が", to: "ガ" },
    { from: "ぎ", to: "ギ" },
    { from: "ぐ", to: "グ" },
    { from: "げ", to: "ゲ" },
    { from: "ご", to: "ゴ" },
    { from: "ざ", to: "ザ" },
    { from: "じ", to: "ジ" },
    { from: "ず", to: "ズ" },
    { from: "ぜ", to: "ゼ" },
    { from: "ぞ", to: "ゾ" },
    { from: "だ", to: "ダ" },
    { from: "ぢ", to: "ヂ" },
    { from: "づ", to: "ヅ" },
    { from: "で", to: "デ" },
    { from: "ど", to: "ド" },
    { from: "ば", to: "バ" },
    { from: "び", to: "ビ" },
    { from: "ぶ", to: "ブ" },
    { from: "べ", to: "ベ" },
    { from: "ぼ", to: "ボ" },
    { from: "ぱ", to: "パ" },
    { from: "ぴ", to: "ピ" },
    { from: "ぷ", to: "プ" },
    { from: "ぺ", to: "ペ" },
    { from: "ぽ", to: "ポ" },
    { from: "ぁ", to: "ァ" },
    { from: "ぃ", to: "ィ" },
    { from: "ぅ", to: "ゥ" },
    { from: "ぇ", to: "ェ" },
    { from: "ぉ", to: "ォ" },
    { from: "ゃ", to: "ャ" },
    { from: "ゅ", to: "ュ" },
    { from: "ょ", to: "ョ" },
  ];
  for (var i = 0; i < reples_data.length; i++) {
    str = str.replace(new RegExp(reples_data[i].from, "g"), reples_data[i].to);
  }
  kana = str;
  return kana;
}
