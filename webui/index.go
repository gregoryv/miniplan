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
<form method="POST">
<input type=hidden name=submit value=add>
<input type=hidden name=priority value="{{.InsertPrio}}" >
<input type=submit value="+" class="insert" title=insert>
</form>
<a href="#{{.Ref}}" class="idref">#</a>
<a name="{{.Ref}}">{{.Ref}}</a>
</div>

<div class="mid">
<form method="POST" class="change">
<input type=hidden name="uuid" value="{{.Ref}}">
<input type=text name="title" value="{{.Title}}" /><br>
<textarea rows="{{.LineHeight}}" name="description">{{.Description}}</textarea>

<table>
<tr>
<td>
<input type=number min=1 name=priority value="{{.Priority}}"/>
<input type=hidden name=submit value=update>
<input type=submit value=Save class="save btn" title=save>
</form>
</td>

<td>
<form method="POST">
<input type=hidden name="uuid" value="{{.Ref}}">
<input type=hidden name=submit value=delete>
<input type=submit value=Delete class="delete btn">
</form>
</td>
</tr>
</table>


</div>


</div>
{{end}}


<div class="row">

<div class="left">
<form method="POST">
<input type=hidden name=submit value=add>
<input type=hidden name=priority value="{{.LastPriority}}" >
<input type=submit value="+" class="insert" title=insert>
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
