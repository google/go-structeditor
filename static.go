// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package structeditor

import (
	"strings"
)

// Static header content
const STATIC_HEADER = `
<html>
  <head>
    <title>Struct Editor</title>
    <script language="javascript">
      function sendCommand(operator, path, extraArgs) {
        let urlParams = "?operator=" + operator +
            "&path=" + encodeURIComponent(path);
        if (extraArgs) {
          urlParams += extraArgs;
        }
        let req = new XMLHttpRequest();
        req.addEventListener("load", function() {
          // location.reload();
        });
        // todo: listeners for errors
        req.open("post", "${MUTATE_URL}" + urlParams);
        req.send("");
      }

      function update(path, fieldName) {
        let newValue = document.getElementById(fieldName).value;
        sendCommand("set", path, "&value=" + encodeURIComponent(newValue));
      }

      function grow(path) {
        sendCommand("grow", path);
      }

      function shrink(path) {
        sendCommand("shrink", path);
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
