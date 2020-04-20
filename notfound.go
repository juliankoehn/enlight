package enlight

const notFoundPage = `
	<!DOCTYPE html>
<html>
<head>
    <title>Page not found - 404</title>
</head>
<body>
 
 
The page your looking for is not available
 
</body>
</html>
`

// NotFoundHandler is the default 404 handler
func NotFoundHandler(c Context) error {
	if c.WantsJSON() {
		return ErrNotFound
	}
	return c.HTML(404, notFoundPage)
}
