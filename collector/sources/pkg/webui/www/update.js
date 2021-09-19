window.setTimeout(update,1000)

let rows

window.onload = function () {
    rows = document.getElementsByClassName("rows")[0]
    for(let i=0;i<15;i++) {
        addrow(i)
    }
}

function update() {
}

function addrow(num) {
    // name
    let d = document.createElement("div");
    d.setAttribute('id', 'name'+num)
    d.innerHTML = "<input type=\"text\" maxlength=\"25\" size=\"25\">"
    rows.appendChild(d)
    // value
    d = document.createElement("div");
    d.setAttribute('id', 'data'+num)
    d.innerHTML = "-2"
    rows.appendChild(d)
    // type
    d = document.createElement("div");
    d.setAttribute('id', 'dtype'+num)
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
    d.setAttribute('id', 'ipaddr'+num)
    d.innerHTML = "<input type=\"text\" minlength=\"7\" maxlength=\"15\" size=\"15\" pattern=\"^((\\d{1,2}|1\\d\\d|2[0-4]\\d|25[0-5])\\.){3}(\\d{1,2}|1\\d\\d|2[0-4]\\d|25[0-5])$\">"
    rows.appendChild(d)
    // rs-485 address
    d = document.createElement("div");
    d.setAttribute('id', 'rs485addr'+num)
    d.innerHTML = "<input type=\"text\" minlength=\"1\" maxlength=\"3\" size=\"3\" pattern=\"^\\d{1,3}$\">"
    rows.appendChild(d)
}
