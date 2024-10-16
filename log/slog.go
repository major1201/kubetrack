package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
)

func NewExtendedTextHandler(w io.Writer, opts *slog.HandlerOptions, showExtraErrorStack bool) *ExtendedTextHandler {
	textHandler := slog.NewTextHandler(w, opts)
	return &ExtendedTextHandler{
		TextHandler:         textHandler,
		showExtraErrorStack: showExtraErrorStack,
	}
}

type ExtendedTextHandler struct {
	*slog.TextHandler

	showExtraErrorStack bool
}

func (h *ExtendedTextHandler) Handle(ctx context.Context, r slog.Record) error {
	err := h.TextHandler.Handle(ctx, r)
	if err != nil {
		return err
	}

	// print err stack if err exists
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == "err" {
			fmt.Printf("%+v\n", a.Value.Any())
			return false
		}
		return true
	})
	return err
}

func NewExtendedJSONHandler(w io.Writer, opts *slog.HandlerOptions, addErrStackAttr bool) *ExtendedJSONHandler {
	jsonHandler := slog.NewJSONHandler(w, opts)
	return &ExtendedJSONHandler{
		JSONHandler:     jsonHandler,
		addErrStackAttr: addErrStackAttr,
	}
}

type ExtendedJSONHandler struct {
	*slog.JSONHandler

	addErrStackAttr bool
}

func (h *ExtendedJSONHandler) Handle(ctx context.Context, r slog.Record) error {
	var (
		errStack         string
		foundErrStackKey bool
	)

	r.Attrs(func(a slog.Attr) bool {
		if a.Key == "err" {
			errStack = fmt.Sprintf("%+v", a.Value.Any())
			return !foundErrStackKey
		}
		if a.Key == "errstack" {
			foundErrStackKey = true
			return errStack == ""
		}
		return true
	})

	if errStack != "" && !foundErrStackKey {
		r.AddAttrs(slog.String("errstack", errStack))
	}

	return h.JSONHandler.Handle(ctx, r)
}
