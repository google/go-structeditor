# Struct Editor

User interface library for editing arbitrary structures using a
UI rendered in HTML.

## Background

Struct Editor is intended as a debugging tool for observing and modifying state
inside a Go server accessible via HTTP. It was originally written to simplify
debugging of a client-server turn-based game engine, but can be applied to any
circumstance where global data is represented as a struct.

## Usage

An example of configuring the editor is provided in the `viewer` example in
the `examples` subdirectory.

```go
	editor := structeditor.ServeEditor(demoData, "/", http.DefaultServeMux)

	log.Fatal(http.ListenAndServe(":8000", nil))
```

Once the server is running, you can view the demoData structure at
http://localhost:8000/. Making edits to the structure will modify the structure
on the server.

## Known Issues / Future Work

* Mutation of the struct in the editor is not synchronized or protected against
  multithreaded access
* Private members of structs cannot be mutated
* Several Go types cannot be rendered
    * complex
    * interface
    * map
* General UI usability cleanups
    * Errors are not reported
    * The UI does not notify the user when a change is committed
    * The UI is not reloaded upon change
    * Boolean data types are exposed as string fields, not dropdowns or checkboxes
    * Newline and comma misplacement
* Structure cannot be modified (slices cannot have elements added / removed;
  pointers cannot be cleared)
* Extremely large structs can bog down the UI

## Security Notice

Use of this library exposes internal state of your server directly to an
insecure HTTP endpoint. If you do not control access to the server, you should
wrap access to the editor's ViewHandler and MutateHandler in an authentication /
authorization solution.

## Disclaimer

This is not an officially supported Google product.
