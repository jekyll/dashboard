package dashboard

import (
	"html/template"
)

type templateInfo struct {
	Projects []*Project
}

var (
	indexTmpl = template.Must(template.New("index.html").Parse(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
  <title>Dashboard</title>
</head>
<body>

<table>
  <caption>Jekyll At-a-Glance Dashboard</caption>
  <thead>
    <tr>
      <th>Repo</th>
      <th>Gem</th>
      <th>Travis</th>
      <th>Downloads</th>
      <th>Commits</th>
      <th>Pull Requests</th>
      <th>Issues</th>
      <th>Unreleased commits</th>
    </tr>
  </thead>
  <tbody>
{{range .Projects}}
<tr>
  <td>
	<a href="https://github.com/{{.Nwo}}" title="{{.Name}} on GitHub">{{.Name}}</a>
  </td>
  <td>{{if .Gem}}
    <a href="https://rubygems.org/gems/{{.Gem.Name}}" title="{{.Gem.Name}} homepage">v{{.Gem.Version}}</a>
  {{end}}</td>
  <td>{{if .Travis}}
    <a href="https://travis-ci.org/{{.Travis.Nwo}}/builds/{{.Travis.Branch.Id}}">{{.Travis.Branch.State}}</a>
  {{end}}</td>
  <td>{{if .Gem}}{{.Gem.Downloads}}{{else}}no info{{end}}</td>
  <td>{{if ge .GitHub.CommitsThisWeek 0}}{{.GitHub.CommitsThisWeek}}{{else}}no info{{end}}</td>
  <td>
    {{if gt .GitHub.OpenPRs 0}}
      <a href="https://github.com/{{.Nwo}}/pulls">{{.GitHub.OpenPRs}}<a/>
    {{else if ge .GitHub.OpenPRs 0}}
      {{.GitHub.OpenPRs}}
	{{else}}
      no info
    {{end}}
  </td>
  <td>
    {{if gt .GitHub.OpenIssues 0}}
      <a href="https://github.com/{{.Nwo}}/issues">{{.GitHub.OpenIssues}}<a/>
    {{else if ge .GitHub.OpenIssues 0}}
      {{.GitHub.OpenIssues}}
    {{else}}
      no info
    {{end}}
  </td>
  <td>
    {{if gt .GitHub.CommitsSinceLatestRelease 0}}
      <a href="https://github.com/{{.Nwo}}/compare/{{.GitHub.LatestReleaseTag}}...master">{{.GitHub.CommitsSinceLatestRelease}}<a/>
    {{else if ge .GitHub.CommitsSinceLatestRelease 0}}
      {{.GitHub.CommitsSinceLatestRelease}}
    {{else}}
      no info
    {{end}}
  </td>
</tr>
{{end}}
  </tbody>
</table>

<div>
	<strong>Commits are as of this week. Issues and pull requests are total open.</strong>
	<a href="https://github.com/jekyll/dashboard">Source Code</a>.
</div>

<p>
	Look wrong? <form action="/reset.json" method="post"><input type="Submit" value="Reset the cache."></form>
</p>

</body>
</html>
`))
)
