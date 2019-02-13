// Static header content

package structeditor

import (
	"strings"
)

const STATIC_HEADER = `
<html>
  <head>
    <title>Struct Editor</title>
    <script language="javascript">
      function update(path, fieldName) {
        let newValue = document.getElementById(fieldName).value;
        let data = new FormData();
        data.append("path", path);
        data.append("value", newValue);

        let req = new XMLHttpRequest();
        req.addEventListener("load", function() {
          location.reload();
        });
        // todo: listeners for errors
        req.open("post", "${MUTATE_URL}");
        req.send(data);
      }
    </script>
  </head>
  <body>
`

const STATIC_FOOTER = `
  </body>
</html>
`

func wrapContent(content, mutateUrl string) string {
	return strings.Replace(
		STATIC_HEADER,
		"${MUTATE_URL}",
		mutateUrl, -1) + content + STATIC_FOOTER
}

// TODO:
// adjust renderer to carry path state
// adjust renderer to use wrapContent
// adjust render tests to verify wrapped content
// adjust renderer to render update buttons on paths to content
// add mutator.go (change value of a struct element based on path)
