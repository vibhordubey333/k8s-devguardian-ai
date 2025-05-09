package output

import (
	"bytes"
	"html/template"
	"strings"
	"time"
)

// HTMLFormatter implements the Formatter interface for HTML output
type HTMLFormatter struct{}

// NewHTMLFormatter creates a new HTML formatter
func NewHTMLFormatter() *HTMLFormatter {
	return &HTMLFormatter{}
}

// Format formats the audit result as HTML
func (f *HTMLFormatter) Format(result AuditResult) ([]byte, error) {
	var buf bytes.Buffer

	tmpl, err := template.New("report").Funcs(template.FuncMap{
		"severityClass": func(severity string) string {
			switch strings.ToLower(severity) {
			case "critical":
				return "critical"
			case "high":
				return "high"
			case "medium":
				return "medium"
			case "low":
				return "low"
			default:
				return "info"
			}
		},
		"severityIcon": func(severity string) string {
			switch strings.ToLower(severity) {
			case "critical":
				return "🔴"
			case "high":
				return "🟠"
			case "medium":
				return "🟡"
			case "low":
				return "🟢"
			default:
				return "⚪"
			}
		},
	}).Parse(htmlTemplate)

	if err != nil {
		return nil, err
	}

	data := struct {
		Result    AuditResult
		Timestamp string
	}{
		Result:    result,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// HTML template for the report
const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Kubernetes Security Audit Report</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
        }
        header {
            background-color: #1a73e8;
            color: white;
            padding: 20px;
            border-radius: 5px;
            margin-bottom: 20px;
        }
        h1, h2, h3 {
            margin-top: 0;
        }
        .summary {
            display: flex;
            justify-content: space-between;
            margin-bottom: 30px;
        }
        .summary-box {
            background-color: #f5f5f5;
            border-radius: 5px;
            padding: 15px;
            flex: 1;
            margin-right: 15px;
        }
        .summary-box:last-child {
            margin-right: 0;
        }
        .finding {
            background-color: #fff;
            border: 1px solid #ddd;
            border-radius: 5px;
            padding: 20px;
            margin-bottom: 20px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .finding-header {
            display: flex;
            justify-content: space-between;
            border-bottom: 1px solid #eee;
            padding-bottom: 10px;
            margin-bottom: 15px;
        }
        .severity {
            font-weight: bold;
            padding: 5px 10px;
            border-radius: 3px;
            color: white;
        }
        .critical {
            background-color: #d32f2f;
        }
        .high {
            background-color: #f57c00;
        }
        .medium {
            background-color: #fbc02d;
            color: #333;
        }
        .low {
            background-color: #388e3c;
        }
        .info {
            background-color: #0288d1;
        }
        .section {
            margin-bottom: 15px;
        }
        .section-title {
            font-weight: bold;
            margin-bottom: 5px;
        }
        .references {
            list-style-type: none;
            padding-left: 0;
        }
        .references li {
            margin-bottom: 5px;
        }
        .references a {
            color: #1a73e8;
            text-decoration: none;
        }
        .references a:hover {
            text-decoration: underline;
        }
        footer {
            margin-top: 30px;
            text-align: center;
            color: #666;
            font-size: 0.9em;
        }
    </style>
</head>
<body>
    <header>
        <h1>Kubernetes Security Audit Report</h1>
        <p>Generated on: {{.Timestamp}}</p>
    </header>

    <section class="summary">
        <div class="summary-box">
            <h2>Summary</h2>
            <p>Total Findings: <strong>{{.Result.Summary.TotalFindings}}</strong></p>
            <h3>By Severity</h3>
            <ul>
                {{range $severity, $count := .Result.Summary.BySeverity}}
                <li>{{severityIcon $severity}} {{$severity}}: {{$count}}</li>
                {{end}}
            </ul>
        </div>
        <div class="summary-box">
            <h2>Resource Types</h2>
            <ul>
                {{range $resource, $count := .Result.Summary.ByResource}}
                <li>{{$resource}}: {{$count}}</li>
                {{end}}
            </ul>
        </div>
    </section>

    <h2>Detailed Findings</h2>

    {{range $index, $explanation := .Result.Explanations}}
    <div class="finding">
        <div class="finding-header">
            <h3>Finding #{{add $index 1}}: {{$explanation.Finding.Resource}}/{{$explanation.Finding.Namespace}}/{{$explanation.Finding.Name}}</h3>
            <span class="severity {{severityClass $explanation.Finding.Severity}}">{{severityIcon $explanation.Finding.Severity}} {{$explanation.Finding.Severity}}</span>
        </div>

        <div class="section">
            <div class="section-title">Issue:</div>
            <p>{{$explanation.Finding.Reason}}</p>
        </div>

        <div class="section">
            <div class="section-title">Explanation:</div>
            <p>{{$explanation.Explanation}}</p>
        </div>

        <div class="section">
            <div class="section-title">Remediation:</div>
            <p>{{$explanation.Remediation}}</p>
        </div>

        <div class="section">
            <div class="section-title">References:</div>
            <ul class="references">
                {{range $ref := $explanation.References}}
                <li><a href="{{$ref}}" target="_blank">{{$ref}}</a></li>
                {{end}}
            </ul>
        </div>
    </div>
    {{end}}

    <footer>
        <p>Generated by K8s DevGuardian AI - An AI-powered Kubernetes security auditing tool</p>
    </footer>
</body>
</html>`
