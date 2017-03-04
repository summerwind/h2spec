package reporter

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/summerwind/h2spec/config"
	"github.com/summerwind/h2spec/log"
	"github.com/summerwind/h2spec/spec"
)

const homeTemplate string = `
<html>
<head><title>h2specd Report</title></head>
<body>
<script>
function runTests() {
	var links = document.getElementsByTagName("a");
	for (var i = 0; i < links.length; i++) {
		var link = links[i].getAttribute("href");
		console.log(link);
		try {
			var xhr = new XMLHttpRequest();
			xhr.open("GET", link, false);
			xhr.send();
		} catch(err) {
		}
	}

	var xhr = new XMLHttpRequest();
	xhr.open("GET", "/report", false);
	xhr.send();
	document.getElementById("report").innerHTML = xhr.responseText;
}
</script>
<button type="button" onclick="runTests()">Run Tests</button>
<div id="report">
  %s
</div>
</body>
</html>
`

type WebReportServer struct {
	http.Server

	config *config.Config
	spec   *spec.ClientTestGroup
}

func NewWebReportServer(config *config.Config, spec *spec.ClientTestGroup) *WebReportServer {
	server := &WebReportServer{
		Server: http.Server{Addr: config.Addr()},
		config: config,
		spec:   spec,
	}

	handler := http.NewServeMux()
	handler.HandleFunc("/", server.home)
	handler.HandleFunc("/report", server.report)

	server.Handler = handler

	return server
}

func (server *WebReportServer) RunForever() error {
	log.Println(fmt.Sprintf("Report server is listened at http://%s", server.config.Addr()))
	return server.ListenAndServe()
}

func (server *WebReportServer) home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, homeTemplate, htmlReport(server.spec, server.config))
}

func (server *WebReportServer) report(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, htmlReport(server.spec, server.config))
}

func htmlReport(tg *spec.ClientTestGroup, c *config.Config) string {
	var buffer bytes.Buffer

	passed := tg.PassedCount
	failed := tg.FailedCount
	skipped := tg.SkippedCount

	total := passed + failed + skipped
	tmp := "<div>%d tests, %d passed, %d skipped, %d failed</div>"
	buffer.WriteString(fmt.Sprintf(tmp, total, passed, skipped, failed))

	buffer.WriteString(htmlReportForTestGroup(tg, c))
	return buffer.String()
}

func htmlReportForTestGroup(tg *spec.ClientTestGroup, c *config.Config) string {
	mode := c.RunMode(tg.ID())
	if mode == config.RunModeNone {
		return ""
	}

	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("<div>%s</div>", tg.Title()))

	for _, tc := range tg.Tests {
		buffer.WriteString(htmlReportForTestCase(tc, c))
	}

	for _, g := range tg.Groups {
		buffer.WriteString(htmlReportForTestGroup(g, c))
	}

	buffer.WriteString("<br>")
	return buffer.String()
}

func htmlReportForTestCase(tc *spec.ClientTestCase, c *config.Config) string {
	formatter := "<div>%s<a href=\"%s\" target=\"_blank\">%s</a>%s</div>"

	tr := tc.Result

	if tr == nil {
		resultLabel := "<span style=\"color: red;\">&nbsp;&nbsp;</span>"
		return fmt.Sprintf(formatter, resultLabel, tc.FullPath(c), tc.FullPath(c), tc.Desc)
	}

	if !tr.Failed {
		resultLabel := "<span style=\"color: green;\">✔</span>"
		return fmt.Sprintf(formatter, resultLabel, tc.FullPath(c), tc.FullPath(c), tc.Desc)
	}

	var buffer bytes.Buffer

	resultLabel := "<span style=\"color: red;\">✖</span>"
	buffer.WriteString(fmt.Sprintf(formatter, resultLabel, tc.FullPath(c), tc.FullPath(c), tc.Desc))

	err, ok := tr.Error.(*spec.TestError)
	formatter = "<div style=\"padding-left: %dpx; color: %s\">%s</div>"
	if ok {
		msg := fmt.Sprintf("-> %s", tc.Requirement)
		buffer.WriteString(fmt.Sprintf(formatter, 20, "red", msg))

		label := "Expected:"
		for i, ex := range err.Expected {
			if i != 0 {
				label = strings.Repeat("&nbsp;", len(label))
			}
			msg = fmt.Sprintf("%s&nbsp;%s", label, ex)
			buffer.WriteString(fmt.Sprintf(formatter, 30, "yellow", msg))
		}
		msg = fmt.Sprintf("&nbsp;&nbsp;Actual:&nbsp;%s", err.Actual)
		buffer.WriteString(fmt.Sprintf(formatter, 30, "green", msg))

	} else if err != nil {
		errMsg := fmt.Sprintf("Error: %v", err)
		buffer.WriteString(fmt.Sprintf(formatter, 20, "red", errMsg))
	} else {
		errMsg := fmt.Sprintf("Error: %v", tr.Error.Error())
		buffer.WriteString(fmt.Sprintf(formatter, 20, "red", errMsg))
	}

	return buffer.String()
}
