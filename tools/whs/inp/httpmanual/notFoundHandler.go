package httpmanual


func notFoundHandler(request *Request, resp *responseWriter) {
	Error404(request, resp)
}