package tmpl

const HtmlFile = `<html>
<head>
	 <meta charset="utf-8">
    <script src="https://cdn.jsdelivr.net/npm/json-schema-view-js@2.0.1/dist/bundle.min.js"></script>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/json-schema-view-js@2.0.1/dist/style.min.css">
    <style>
        .json-schema-view-dark {
            background: #222;
        }
        h3 + div > div {
            padding: 1em;
        }
        pre {
            background: #333;
            color: white;
            padding: 10px;
            font-size: 1.2em;
            line-height: 1.3;
        }
        /*
         * Page layout
        */
        body {
            max-width: 80vw;
            margin: 30px auto;
            font-family: sans-serif;
        }
        header {
            overflow: hidden;
        }
        header h1 {
            float: left;
        }
        header > div {
            text-align: right;
            padding: 10px;
        }
        header > div * {
            vertical-align: middle;
        }
        .code-in-gh {
            background: steelblue;
            color: white;
            padding: .5em;
            border-radius: 5px;
            text-decoration: none;
        }
    </style>
</head>
<body>
<div class="results"></div>
<script>
    var schema = [
        {{ . }}
    ];
    var results = document.querySelector('.results');
    document.addEventListener('DOMContentLoaded', function () {
        schema.forEach(function (example) {
            var title = document.createElement('h3');
            var formatter = new JSONSchemaView(example.schema, 1, { theme: example.theme });
            title.innerText = example.title;
            results.appendChild(title)
            results.appendChild(formatter.render());
        });
    });
</script>
</body>
</html>
`
