window.setInterval(update,2000)

let rows
const nums = 15

window.onload = function () {
    rows = document.getElementsByClassName("rows")[0]
    for(let i=0;i<nums;i++) {
        addrow(i)
    }
}

function update() {
    let a = ""
    for(let i=0;i<nums;i++) {

    }
    ajaxRequest()
}

function ajaxRequest() {
    httpRequest = new XMLHttpRequest();
    if (!httpRequest) {
        alert('Giving up :( Cannot create an XMLHTTP instance');
        return false;
    }
    httpRequest.timeout = 2000
    httpRequest.ontimeout = function (e) {
        document.getElementsByTagName("footer")[0].innerText="Нет соединения с сервером"
    }
    httpRequest.onreadystatechange = ajaxUpdate;

    let params = "?"
    let separator = ""
    for(let i=0;i<nums;i++) {
        let id = separator+"dtype"+i
        separator = "&"
        params += id + "=" + document.getElementById("dtype0").getElementsByTagName("input")[0].value.slice(2,4)
        id = "ipaddr"+i
        params += "&" + id + "=" + document.getElementById(id).getElementsByTagName("input")[0].value
        id = "rs485addr"+i
        params += "&" + id + "=" + document.getElementById(id).getElementsByTagName("input")[0].value
    }
    httpRequest.open('GET', 'ajax_update'+params);
    httpRequest.send();
}

function ajaxUpdate() {
    if (httpRequest.readyState === XMLHttpRequest.DONE) {
        if (httpRequest.status === 200) {
            document.getElementsByTagName("footer")[0].innerText=httpRequest.responseText
            let params = JSON.parse(httpRequest.responseText)
            for(let i=0;i<nums;i++){
                if(params[i].ready === true) {
                    document.getElementById("data" + i).innerText = params[i].data
                } else {
                    document.getElementById("data" + i).innerText = "-"
                }
            }
        } else {
            document.getElementsByTagName("footer")[0].innerText="Нет соединения с сервером"
        }
    }
}

function addrow(i) {
    // name
    let d = document.createElement("div");
    d.setAttribute('id', 'name'+i)
    d.innerHTML = "<input type=\"text\" maxlength=\"25\" size=\"25\">"
    rows.appendChild(d)
    // value
    d = document.createElement("div");
    d.setAttribute('id', 'data'+i)
    d.innerHTML = "-2"
    rows.appendChild(d)
    // type
    d = document.createElement("div");
    d.setAttribute('id', 'dtype'+i)
    d.className+='dropdown'
    d.innerHTML = "<input type=\"text\" />\n" +
                  "<select  onchange=\"this.previousElementSibling.value=this.value; this.previousElementSibling.focus()\">\n" +
                  "<option>0x?? Свое значение</option>\n" +
                  "<option>0x60 Накопление</option>\n" +
                  "<option>0x5d Производительность v2</option>\n" +
                  "<option>0x3f Производительность v1</option>\n" +
                  "</select>\n"
    rows.appendChild(d)
    // ip address
    d = document.createElement("div");
    d.setAttribute('id', 'ipaddr'+i)
    d.innerHTML = "<input type=\"text\" minlength=\"7\" maxlength=\"15\" size=\"15\" pattern=\"^((\\d{1,2}|1\\d\\d|2[0-4]\\d|25[0-5])\\.){3}(\\d{1,2}|1\\d\\d|2[0-4]\\d|25[0-5])$\">"
    rows.appendChild(d)
    // rs-485 address
    d = document.createElement("div");
    d.setAttribute('id', 'rs485addr'+i)
    d.innerHTML = "<input type=\"text\" minlength=\"1\" maxlength=\"3\" size=\"3\" pattern=\"^\\d{1,3}$\">"
    rows.appendChild(d)
}
