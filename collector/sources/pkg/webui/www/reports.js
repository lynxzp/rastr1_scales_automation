function refreshReports() {
    let t = new Date()
    formatTime(t)
    let todayStart = t.setHours(0,0,0,0)
    let todayFinish = t.setHours( 23, 59, 59, 999)
    let monthStart = new Date(t.setDate(1))
    let monthFinish = new Date(monthStart.getFullYear(), monthStart.getMonth() + 1, 0)
    monthFinish = t.setDate(monthFinish.getDate())
    let yearstart  = new Date(monthStart.setMonth(0))
    let yearfinish = new Date(monthStart.setMonth(11))
    let params = [];
    params.push({start:formatTime(todayStart), end:formatTime(todayFinish), shift: 1, column: "todaycolshift1"})
    params.push({start:formatTime(todayStart), end:formatTime(todayFinish), shift: 2, column: "todaycolshift2"})
    params.push({start:formatTime(todayStart), end:formatTime(todayFinish), shift: 3, column: "todaycolshift3"})
    params.push({start:formatTime(todayStart), end:formatTime(todayFinish), shift: 0, column: "todaycol"})
    params.push({start:formatTime(monthStart), end:formatTime(monthFinish), shift: 0, column: "monthcol"})
    params.push({start:formatTime(yearstart), end:formatTime(yearfinish), shift: 0, column: "yearcol"})
    reportRequest(params)
    // todo:
    // GetReport(todayStart, todayFinish, 0, "customcol")
}

function formatTime(t) {
    t = new Date(t)
    return pad(t.getDate()) + "." + pad((t.getMonth()+1)) + "." + t.getFullYear() + " " + pad(t.getHours()) + ":" + pad(t.getMinutes()) + ":" + pad(t.getSeconds())
}

function pad(num, size) {
    if (size === undefined) {
        size = 2
    }
    num = num.toString()
    while (num.length < size) {
        num = "0" + num
    }
    return num
}

let httpReportRequest
function reportRequest(struct) {
    httpReportRequest = new XMLHttpRequest();
    if (!httpReportRequest) {
        alert('Giving up :( Cannot create an XMLHTTP instance');
        return false;
    }
    httpReportRequest.timeout = 2000
    httpReportRequest.ontimeout = function (e) {
        document.getElementsByTagName("footer")[0].innerText="Таймаут соединения с сервером"
    }
    httpReportRequest.onreadystatechange = reportResponse;
    params = "?params=" + encodeURI(JSON.stringify(struct))

    httpReportRequest.open('GET', 'report'+params);
    httpReportRequest.send();
}

let scfr=['0', '00*40', '00*70', '020*40', '020*70', '15*10', '210*20', '25*20', '3песок(0*5)', '45*10', '510*20', '65*20', '7песок(0-5)']

function reportResponse() {
    if (httpReportRequest.readyState === XMLHttpRequest.DONE) {
        if (httpReportRequest.status === 200) {
            document.getElementsByTagName("footer")[0].innerText=""
            let params = JSON.parse(httpReportRequest.responseText)
            for (let i=0; i<6; i++) {
                for (let j=0; j<scfr.length; j++) {
                    if( params[i].accumulation[scfr[j]] !== undefined) {
                        document.getElementById(params[i].column + "_0").innerText = params[i].accumulation[scfr[j][0]]
                    }
                }
            }
            // console.log(params)
            // todo: something
        } else {
            document.getElementsByTagName("footer")[0].innerText="Нет соединения с сервером (статус: " +
                httpReportRequest.status + ")"
        }
    }
}
