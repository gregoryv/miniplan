package webui

import (
	_ "embed"
	"html/template"
	"net/http"
)

var index = `
<!doctype html>

<html>
<head><title>miniplan</title>

<link rel="stylesheet" type="text/css" href="/static/theme.css" />
<script src="/static/tools.js"></script>

</head>
<body id="body">

<h1>Miniplan</h1>

{{range .Changes}}

<div class="row entry">

<div class="left">
<a href="#{{.Ref}}" class="idref">#</a>
<a name="{{.Ref}}">{{.Ref}}</a>
</div>

<div class="mid">
<form method="POST" class="task">
<input type=hidden name="uuid" value="{{.Ref}}">
<input type=text name="title" value="{{.Title}}" /><br>
<textarea rows="{{.LineHeight}}" name="description">{{.Description}}</textarea>
<input name=priority value="{{.Priority}}"/>
<input type=hidden name=submit value=update>
<input type=submit value=S class="one-letter">
</form>
</div>


<div class="right">
<form method="POST" class="one">
<input type=hidden name=submit value=add>
<input type=submit value=I>
<input type=hidden name=priority value="{{.InsertPrio}}" >
</form>
<form method="POST" class="one">
<input type=hidden name="uuid" value="{{.Ref}}">
<input type=hidden name=submit value=delete>
<input type=submit value=D>
</form>
</div>

</div>
{{end}}


<div class="row">

<div class="left"></div>
<div class="mid"></div>
<div class="right">

<div class="right">
<form method="POST" class="one">
<input type=hidden name=submit value=add>
<input type=submit value=I>
<input type=hidden name=priority value="{{.LastPriority}}" >
</form>

</div>
</div>
</body>
</html>
`

var tpl = template.Must(template.New("").Parse(index))

// static assets

func serveTheme(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/css")
	w.Write(theme)
}

//go:embed assets/theme.css
var theme []byte

func serveTools(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/javascript")
	w.Write(tools)
}

//go:embed assets/tools.js
var tools []byte
