package inference

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type Request struct {
	Sentence string `json:"sentence"`
}

// Runner sends each rune of the provided string to the channel.
func Runner(runeChan chan<- string) {
	runeChan <- `event: message_start`
	data := "👋🏼 Hi, I'm an AI here! How may I help you?"
	for _, d := range data {
		runeChan <- string(d)
	}
	runeChan <- "event: message_end"
	close(runeChan)
}

func StreamText(c echo.Context) error {
	req := new(Request)
	if err := c.Bind(req); err != nil {
		return err
	}

	w := c.Response().Writer
	flusher, ok := w.(http.Flusher)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Streaming unsupported")
	}

	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().WriteHeader(http.StatusOK)

	// Send initial message
	_, streamErr := fmt.Fprintf(w, "event: start_stream\n\n")
	if streamErr != nil {
		return streamErr
	}
	flusher.Flush()

	cliveResponseChan := make(chan string)
	go Runner(cliveResponseChan)

	for ch := range cliveResponseChan {
		_, err := fmt.Fprintf(w, "data: %v\n\n", ch)
		if err != nil {
			return err
		}
		flusher.Flush()
	}

	// Send end message
	_, endStreamErr := fmt.Fprintf(w, "event: end_stream\n\n")
	if endStreamErr != nil {
		return endStreamErr
	}
	flusher.Flush()

	return nil
}
