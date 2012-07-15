# Vision

```Go
	ctx, _ := script.NewContext("/home/git/repositories")
	ctx.Environment.Add(script.DefaultEnvironment)
	for _, file := range ctx.Glob("*") {
		ctx, _ = ctx.Cd(file.Path())
		_, e := ctx.Run("git clone . /tmp/%s/goroot/src/%s", file.ShellFullname()
		if e != nil {
			// Handling
		}
		ctx = ctx.Cd("/tmp/%s/goroot/src/%s", file.ShellFullname(), file.Shellname())
		_, e = ctx.Run("go get .")
		if e != nil {
			// Handling
		}
		for _, file := range ctx.Glob("bin/*") {
			e := ctx.Copy("/tmp", file.Path())
			if e != nil {
				// Handling
			}
		}
		e := ctx.Copy("/tmp", ctx.Glob("data")[0].Path())
		if e != nil {
			// Handling
		}
	}
```
