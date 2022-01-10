function refreshReports() {
    // if(document.getElementById("tseha").style.display !== "block")
    //     return
    let t = new Date()
    formatTime(t)
    let todayStart = t.setHours(0,0,0,0)
    let todayFinish = t.setHours( 23, 59, 59, 999)
    let params = [];
    params.push({start:formatTime(todayStart), end:formatTime(todayFinish), shift: 1, column: "todaycolshift1"})
    params.push({start:formatTime(todayStart), end:formatTime(todayFinish), shift: 2, column: "todaycolshift2"})
    params.push({start:formatTime(todayStart), end:formatTime(todayFinish), shift: 3, column: "todaycolshift3"})
    params.push({start:formatTime(todayStart), end:formatTime(todayFinish), shift: 0, column: "todaycol"})
    let monthStart = new Date(t.setDate(1))
    let monthFinish = new Date(monthStart.getFullYear(), monthStart.getMonth() + 1, 0)
    monthFinish = t.setDate(monthFinish.getDate())
    params.push({start:formatTime(monthStart), end:formatTime(monthFinish), shift: 0, column: "monthcol"})
    let yearstart  = new Date(monthStart.setMonth(0))
    let yearfinish = new Date(monthStart.setMonth(11))
    params.push({start:formatTime(yearstart), end:formatTime(yearfinish), shift: 0, column: "yearcol"})

    periodStart = document.getElementById("periodFrom").value+":00"
    periodEnd = document.getElementById("periodTo").value+":00"
    params.push({start: periodStart, end: periodEnd, shift:0, column: "customcol"})
    // console.log(params)
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

let scfr=[['0_песок(0-5)', 11], ['1_0*40', 0], ['1_0*70', 1], ['1_20*40', 2], ['1_20*70', 3] , ['2_5*10', 4], ['3_10*20', 5] ,['3_5*20', 6], ['4_песок(0-5)', 7], ['5_5*10', 8], ['6_10*20', 9], ['7_5*20', 10]]

function reportResponse() {
    if (httpReportRequest.readyState === XMLHttpRequest.DONE) {
        if (httpReportRequest.status === 200) {
            document.getElementsByTagName("footer")[0].innerText=""
            try {
                let params = JSON.parse(httpReportRequest.responseText)
                console.log(params)
                for (let i = 0; i < 7; i++) {
                    for (let j = 0; j < scfr.length; j++) {
                        let val = "-"
                        if ((params[i] !== undefined) && (params[i].accumulation !== null) && (params[i].accumulation[scfr[j][0]] !== undefined)) {
                            val = params[i].accumulation[scfr[j][0]]
                        }
                        document.getElementById(params[i].column + "_" + scfr[j][1]).innerText = val
                    }
                }
            } catch (e) {
                console.log(e)
                console.log(httpReportRequest.responseText)
            }
        } else {
            document.getElementsByTagName("footer")[0].innerText="Нет соединения с сервером (статус: " +
                httpReportRequest.status + ")"
        }
    }
}
